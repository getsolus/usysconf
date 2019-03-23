/*
 * This file is part of usysconf.
 *
 * Copyright © 2017-2018 Solus Project
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 */

#define _GNU_SOURCE

#include "config.h"
#include "context.h"
#include "files.h"
#include "util.h"

static const char *unit_paths[] = {
        SYSTEMD_UTIL_DIR "/systemd" /* /usr/lib/systemd/systemd */
};

/**
 * Ask systemd to reexec when the binary has been updated
 *
 */
static UscHandlerStatus usc_handler_systemd_reexec_exec(UscContext *ctx, const char *path)
{
        const char *command[] = {
                "/usr/bin/systemctl", "daemon-reexec", NULL, /* Terminator */
        };

        if (access(path, X_OK) != 0) {
                return USC_HANDLER_SKIP;
        }

        usc_context_emit_task_start(ctx, "Re-executing systemd");

        if (usc_context_has_flag(ctx, USC_FLAGS_CHROOTED) ||
            usc_context_has_flag(ctx, USC_FLAGS_LIVE_MEDIUM) ||
            usc_context_has_flag(ctx, USC_FLAGS_CONTAINER)) {
                usc_context_emit_task_finish(ctx, USC_HANDLER_SKIP);
                return USC_HANDLER_SKIP | USC_HANDLER_BREAK;
        }

        int ret = usc_exec_command((char **)command);
        if (ret != 0) {
                usc_context_emit_task_finish(ctx, USC_HANDLER_FAIL);
                return USC_HANDLER_FAIL | USC_HANDLER_BREAK;
        }
        usc_context_emit_task_finish(ctx, USC_HANDLER_SUCCESS);
        /* Only want to run once for all of our globs */
        return USC_HANDLER_SUCCESS | USC_HANDLER_BREAK;
}

const UscHandler usc_handler_systemd_reexec = {
        .name = "systemd-reexec",
        .description = "Re-execute systemd",
        .required_bin = "/usr/bin/systemctl",
        .exec = usc_handler_systemd_reexec_exec,
        .paths = unit_paths,
        .n_paths = ARRAY_SIZE(unit_paths),
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
