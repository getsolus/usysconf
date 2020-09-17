/*
 * This file is part of usysconf.
 *
 * Copyright © 2020 Solus Project
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
        SYSTEMD_UNIT_DIR "/sockets.target.wants/" /* /usr/lib/systemd/system/sockets.target.wants/ */
};

/**
 * Ask systemd to restart sockets.target when new vendor-enabled .socket units have been added/updated.
 *
 * This assumes that a daemon-reload has been performed first.
 */
static UscHandlerStatus usc_handler_systemd_sockets_exec(UscContext *ctx, const char *path)
{
        const char *command[] = {
                "/usr/bin/systemctl", "restart", "sockets.target", NULL, /* Terminator */
        };

        if (access(path, X_OK) != 0) {
                return USC_HANDLER_SKIP;
        }

        usc_context_emit_task_start(ctx, "Re-starting vendor-enabled .socket units");

        if (usc_context_has_flag(ctx, USC_FLAGS_CHROOTED) ||
            usc_context_has_flag(ctx, USC_FLAGS_LIVE_MEDIUM)) {
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

const UscHandler usc_handler_systemd_sockets = {
        .name = "systemd-sockets",
        .description = "Re-start systemd sockets.target",
        .required_bin = "/usr/bin/systemctl",
        .exec = usc_handler_systemd_sockets_exec,
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
