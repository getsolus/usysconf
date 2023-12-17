/*
 * This file is part of usysconf.
 *
 * Copyright © 2017-2019 Solus Project
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 */

#define _GNU_SOURCE

#include "context.h"
#include "files.h"
#include "util.h"

static const char *gio_modules_paths[] = {
        "/usr/lib64/gio/modules/",
};

/**
 * Create a module cache with metadata from gio modules. Without this gio has to open each module which can cause bugs.
 */
static UscHandlerStatus usc_handler_glib2_gio_exec(__usc_unused__ UscContext *ctx, const char *path)
{
        autofree(char) *fp = NULL;
        char *command[] = {
                "/usr/bin/gio-querymodules",
                NULL, /* /usr/lib64/gio/modules */
                NULL, /* Terminator */
        };

        if (!usc_file_is_dir(path)) {
                return USC_HANDLER_SKIP;
        }

        command[1] = (char *)path,

        usc_context_emit_task_start(ctx, "Creating GIO modules cache");
        int ret = usc_exec_command(command);
        if (ret != 0) {
                usc_context_emit_task_finish(ctx, USC_HANDLER_FAIL);
                return USC_HANDLER_FAIL | USC_HANDLER_BREAK;
        }
        usc_context_emit_task_finish(ctx, USC_HANDLER_SUCCESS);
        /* Only want to run once for all of our globs */
        return USC_HANDLER_SUCCESS | USC_HANDLER_BREAK;
}

const UscHandler usc_handler_glib2_gio = {
        .name = "glib2-gio",
        .description = "Create glib2 GIO modules cache",
        .required_bin = "/usr/bin/gio-querymodules",
        .exec = usc_handler_glib2_gio_exec,
        .paths = gio_modules_paths,
        .n_paths = ARRAY_SIZE(gio_modules_paths),
};

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
