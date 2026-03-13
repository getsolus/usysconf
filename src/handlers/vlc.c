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

static const char *vlc_modules_paths[] = {
        "/usr/lib/vlc/plugins/",
        "/usr/lib64/vlc/plugins/",
};

/**
 * Create a VLC plugin cache
 */
static UscHandlerStatus usc_handler_vlc_exec(UscContext *ctx, const char *path)
{
        autofree(char) *fp = NULL;
        char *command[] = {
                "/usr/lib64/vlc/vlc-cache-gen",
                NULL, /* /usr/lib64/vlc/plugins/ */
                NULL, /* Terminator */
        };

        if (!usc_file_is_dir(path)) {
                return USC_HANDLER_SKIP;
        }

        command[1] = (char *)path,

        usc_context_emit_task_start(ctx, "Creating VLC plugins cache");
        int ret = usc_exec_command(command);
        if (ret != 0) {
                usc_context_emit_task_finish(ctx, USC_HANDLER_FAIL);
                return USC_HANDLER_FAIL | USC_HANDLER_BREAK;
        }
        usc_context_emit_task_finish(ctx, USC_HANDLER_SUCCESS);
        /* Only want to run once for all of our globs */
        return USC_HANDLER_SUCCESS | USC_HANDLER_BREAK;
}

const UscHandler usc_handler_vlc = {
        .name = "vlc",
        .description = "Create VLC Plugins cache",
        .required_bin = "/usr/lib64/vlc/vlc-cache-gen",
        .exec = usc_handler_vlc_exec,
        .paths = vlc_modules_paths,
        .n_paths = ARRAY_SIZE(vlc_modules_paths),
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
