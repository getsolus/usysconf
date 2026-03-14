/*
 * This file is part of usysconf.
 *
 * Copyright © 2026 Solus Project
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

#include <libgen.h>
#include <string.h>

static const char *system_presets_paths[] = {
    "/usr/lib/systemd/system/*.path",
    "/usr/lib/systemd/system/*.service",
    "/usr/lib/systemd/system/*.socket",
    "/usr/lib/systemd/system/*.timer",
};

static const char *user_presets_paths[] = {
    "/usr/lib/systemd/user/*.path",
    "/usr/lib/systemd/user/*.service",
    "/usr/lib/systemd/user/*.socket",
    "/usr/lib/systemd/user/*.timer",
};

static UscHandlerStatus usc_handler_system_presets_exec(UscContext *context, const char *path) {
    if (!usc_file_exists(path)) {
        return USC_HANDLER_SKIP;
    }

    const char *service_name = basename((char *) path);

    usc_context_emit_task_start(context, "Updating systemd system service preset: %s", service_name);

    const char *command[] = {
        "/usr/bin/systemctl",
        "preset",
        service_name,
        "--root=/", /* Ensure no tom-foolery with dbus */
        "--force",
        NULL        /* Terminator */
    };

    int ret = usc_exec_command((char **) command);

    if (ret != 0) {
        usc_context_emit_task_finish(context, USC_HANDLER_FAIL);
        return USC_HANDLER_FAIL;
    }

    usc_context_emit_task_finish(context, USC_HANDLER_SUCCESS);
    /* Only want to run once for all of our globs */
    return USC_HANDLER_SUCCESS;
}

static UscHandlerStatus usc_handler_user_presets_exec(UscContext *context, const char *path) {
    if (!usc_file_exists(path)) {
        return USC_HANDLER_SKIP;
    }

    const char *service_name = basename((char *) path);

    usc_context_emit_task_start(context, "Updating systemd user service preset: %s", service_name);

    const char *command[] = {
        "/usr/bin/systemctl",
        "preset",
        service_name,
        "--root=/", /* Ensure no tom-foolery with dbus */
        "--force",
        "--global", /* Enable for all logins of all users */
        NULL        /* Terminator */
    };

    int ret = usc_exec_command((char **) command);

    if (ret != 0) {
        usc_context_emit_task_finish(context, USC_HANDLER_FAIL);
        return USC_HANDLER_FAIL;
    }

    usc_context_emit_task_finish(context, USC_HANDLER_SUCCESS);
    return USC_HANDLER_SUCCESS;
}

const UscHandler usc_handler_system_presets = {
    .name = "system-presets",
    .description = "Update systemd system service presets",
    .required_bin = "/usr/bin/systemctl",
    .exec = usc_handler_system_presets_exec,
    .paths = system_presets_paths,
    .n_paths = ARRAY_SIZE(system_presets_paths),
};

const UscHandler usc_handler_user_presets = {
    .name = "user-presets",
    .description = "Update systemd user service presets",
    .required_bin = "/usr/bin/systemctl",
    .exec = usc_handler_user_presets_exec,
    .paths = user_presets_paths,
    .n_paths = ARRAY_SIZE(user_presets_paths),
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
