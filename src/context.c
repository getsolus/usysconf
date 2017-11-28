/*
 * This file is part of usysconf.
 *
 * Copyright © 2017 Solus Project
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 */

#define _GNU_SOURCE

#include <errno.h>
#include <fcntl.h>
#include <glob.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/sendfile.h>
#include <sys/stat.h>

#include "config.h"
#include "context.h"
#include "files.h"
#include "handlers.h"
#include "state.h"

/* libuf */
#include "map.h"

/**
 * TODO: Expose this more generically through libuf API
 */
#ifndef UF_INT_TO_PTR
#define UF_INT_TO_PTR(x) ((void *)((uintptr_t)x))
#endif

/* Table of supported handlers */
static const UscHandler *usc_handlers[] = {
        &usc_handler_ldconfig, /**<Get library cache in order first */

#ifdef HAVE_CBM
        &usc_handler_cbm, /**<Now CBM can successfully run */
#endif

#ifdef HAVE_QOL_ASSIST
        &usc_handler_qol_assist, /**<Schedule migration on boot */
#endif

        &usc_handler_depmod, /**< Run depmod after cbm does its thing */

        /** Middleware */
        &usc_handler_hwdb, /**<Want hwdb updated before calling LDM (PCI) */

#ifdef HAVE_LDM
        &usc_handler_ldm, /**<Update drivers/GL-links/etc */
#endif

        /* Very likely that LDM caused a cache invalidation for lib dirs */
        &usc_handler_ldconfig,

#ifdef HAVE_SYSTEMD
        &usc_handler_sysusers,
        &usc_handler_tmpfiles,
        &usc_handler_systemd_reload,
        &usc_handler_systemd_reexec,
#endif

        /** Enter userspace. */
        &usc_handler_glib2,
        &usc_handler_fonts,
        &usc_handler_mime,
        &usc_handler_icon_cache,
        &usc_handler_desktop_files,
        &usc_handler_gconf,
        &usc_handler_dconf,

        /* GTK immodules */
        &usc_handler_gtk2_immodules,
        &usc_handler_gtk3_immodules,

        /* Special cases */
        &usc_handler_mandb,
        &usc_handler_ssl_certs,
        &usc_handler_sshd,
};

/**
 * Opaque implementation allows us to avoid potential issues with methods
 * using our internal state in an invalid fashion, and will also allow us
 * to enforce const usage.
 */
struct UscContext {
        unsigned int flags; /**<A bitwise set of flags specified for the context */

        UfHashmap *skip_map; /**<Allow implementations to track skips */
};

UscContext *usc_context_new()
{
        UscContext *ret = NULL;

        ret = calloc(1, sizeof(UscContext));
        if (!ret) {
                return NULL;
        }

        if (usc_is_chrooted()) {
                ret->flags |= USC_FLAGS_CHROOTED;
        }

        /* Skip map only contains a 1 value */
        ret->skip_map =
            uf_hashmap_new_full(uf_hashmap_string_hash, uf_hashmap_string_equal, free, NULL);
        if (!ret->skip_map) {
                usc_context_free(ret);
                return NULL;
        }

        return ret;
}

void usc_context_free(UscContext *self)
{
        if (!self) {
                return;
        }
        uf_hashmap_free(self->skip_map);
        free(self);
}

bool usc_context_has_flag(UscContext *self, unsigned int flag)
{
        if (!self) {
                return false;
        }
        if ((self->flags & flag) == flag) {
                return true;
        }
        return false;
}

/**
 * An item failed, so wind back our log and spit it back out to the tty
 */
static void usc_spit_fail_log(void)
{
        struct stat st = { 0 };
        int fd = -1;
        ssize_t written;
        ssize_t total;

        fd = open(USYSCONF_LOG_FILE, O_RDONLY);
        if (fd < 0) {
                fprintf(stderr,
                        "open(%s): failed to open log file: %s\n",
                        USYSCONF_LOG_FILE,
                        strerror(errno));
                return;
        }
        fstat(fd, &st);

        /* Pump entire file to stderr */
        total = st.st_size;
        for (;;) {
                written = sendfile(STDERR_FILENO, fd, NULL, (size_t)total);
                if (written == total) {
                        break;
                } else if (written < 0) {
                        fprintf(stderr, "sendfile(): %s\n", strerror(errno));
                        return;
                }
                total -= written;
        }
}

static void usc_handle_one(const UscHandler *handler, UscContext *context, UscStateTracker *tracker)
{
        UscHandlerStatus status = USC_HANDLER_MIN;
        bool record_remain = false;

        for (size_t i = 0; i < handler->n_paths; i++) {
                glob_t glo = { 0 };
                const char *path = NULL;

                path = handler->paths[i];

                if (glob(path, GLOB_NOSORT, NULL, &glo) != 0) {
                        continue;
                }

                for (size_t i = 0; i < glo.gl_pathc; i++) {
                        char *resolved = glo.gl_pathv[i];
                        bool record_path = false;

                        /* Don't try to do anything for the remainder of this glob. */
                        if (record_remain) {
                                goto push_entry;
                        }

                        /* Do we need to handle this dude ? */
                        if (!usc_state_tracker_needs_update(tracker, resolved)) {
                                continue;
                        }

                        status = handler->exec(context, resolved);

                        if ((status & USC_HANDLER_SUCCESS) == USC_HANDLER_SUCCESS) {
                                record_path = true;
                        }
                        if ((status & USC_HANDLER_FAIL) == USC_HANDLER_FAIL) {
                                fputs("Failed\n", stderr);
                                usc_spit_fail_log();
                                continue;
                        }
                        if ((status & USC_HANDLER_SKIP) == USC_HANDLER_SKIP) {
                                continue;
                        }

                        if ((status & USC_HANDLER_BREAK) == USC_HANDLER_BREAK) {
                                /* Record all paths as updated */
                                record_remain = true;
                        }

                        if (!record_path) {
                                continue;
                        }

                push_entry:
                        /* We won't record the new entry here and we won't write it back out either
                         */
                        if ((status & USC_HANDLER_DROP) == USC_HANDLER_DROP) {
                                continue;
                        }
                        if (!usc_state_tracker_push_path(tracker, resolved)) {
                                fprintf(stderr, "Failed to record path %s\n", resolved);
                        }
                }
                globfree(&glo);
        }
}

bool usc_context_run_triggers(UscContext *context, const char *name)
{
        autofree(UscStateTracker) *tracker = NULL;
        bool ran_trigger = false;

        tracker = usc_state_tracker_new();
        if (!tracker) {
                fputs("Cannot continue without valid UscStateTracker\n", stderr);
                return false;
        }

        /* Crack on regardless. */
        if (!usc_state_tracker_load(tracker)) {
                fputs("Invalid state has been removed\n", stderr);
        }

        /* Before we go, make sure log directory exists */
        if (!usc_file_exists(USYSCONF_LOG_DIR) && mkdir(USYSCONF_LOG_DIR, 00755) != 0) {
                fprintf(stderr,
                        "Cannot construct log directory %s: %s\n",
                        USYSCONF_LOG_DIR,
                        strerror(errno));
                return false;
        }

        /* Just test the main loop iteration jank for now */
        for (size_t i = 0; i < ARRAY_SIZE(usc_handlers); i++) {
                if (name && strcmp(usc_handlers[i]->name, name) != 0) {
                        continue;
                }
                usc_handle_one(usc_handlers[i], context, tracker);
                ran_trigger = true;
        }

        if (!ran_trigger) {
                fprintf(stderr, "Unknown trigger '%s'\n", name);
                return false;
        }

        /* Dump it back to disk */
        if (!usc_state_tracker_write(tracker)) {
                fputs("Failed to write state file!\n", stderr);
                return false;
        }

        return true;
}

bool usc_context_push_skip(UscContext *self, char *skip_item)
{
        if (!self) {
                return false;
        }
        return uf_hashmap_put(self->skip_map, strdup(skip_item), UF_INT_TO_PTR(1));
}

bool usc_context_should_skip(UscContext *self, char *skip_item)
{
        if (!self) {
                return false;
        }
        /* Only 0==NULL on glibc, 1 will be set and thus not NULL */
        return uf_hashmap_get(self->skip_map, skip_item) != NULL;
}

/*
 * Editor modelines  -  https://www.wireshark.org/tools/modelines.html
 *
 * Local variables:
 * c-basic-offset: 8
 * tab-width: 8
 * indent-tabs-mode: nil
 * End:
 *
 * vi: set shiftwidth=8 tabstop=8 expandtab:
 * :indentSize=8:tabSize=8:noTabs=true:
 */
