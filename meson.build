project(
    'usysconf',
    ['c'],
    version: '0.5.14',
    license: [
        'GPL-2.0-only',
    ],
    meson_version: '>= 0.57.0',
    default_options: [
        'default_library=static',
        'warning_level=3',
    ],
)

# Get configuration bits together
path_prefix = get_option('prefix')
path_sysconfdir = join_paths(path_prefix, get_option('sysconfdir'))
path_mandir = join_paths(path_prefix, get_option('mandir'))
path_datadir = join_paths(path_prefix, get_option('datadir'))
path_sbindir = join_paths(path_prefix, get_option('sbindir'))
path_vardir = join_paths(path_prefix, get_option('localstatedir'), 'lib', meson.project_name())


# Sort out config.h now
cdata = configuration_data()

# General options..
cdata.set_quoted('PACKAGE_NAME', meson.project_name())
cdata.set_quoted('PACKAGE_VERSION', meson.project_version())
cdata.set_quoted('PACKAGE_URL', 'https://getsol.us')

# Allow building a static binary. For now let's default to dynamic so we can
# actually valgrind the thing.
with_static_binary = get_option('with-static-binary')

with_cbm = get_option('with-clr-boot-manager')
if with_cbm == true
    cdata.set('HAVE_CBM', '1')
endif
with_ldm = get_option('with-linux-driver-management')
if with_ldm == true
    cdata.set('HAVE_LDM', '1')
endif
with_qol = get_option('with-qol-assist')
if with_qol == true
    cdata.set('HAVE_QOL_ASSIST', '1')
endif
with_apparmor = get_option('with-apparmor')
if with_apparmor == true
    cdata.set('HAVE_APPARMOR', '1')
endif
with_mono_certs = get_option('with-mono-certs')
if with_mono_certs == true
    cdata.set('HAVE_MONO_CERTS', '1')
endif

# Work out systemd support
with_systemd = get_option('with-systemd')
if with_systemd == true
    cdata.set('HAVE_SYSTEMD', '1')

    message('Asking pkg-config for systemd\'s directories')
    dep_systemd = dependency('systemd')
    systemd_unit_dir = dep_systemd.get_variable(pkgconfig : 'systemdsystemunitdir')
    systemd_tmpfiles_dir = dep_systemd.get_variable(pkgconfig : 'tmpfilesdir')
    systemd_sysusers_dir = dep_systemd.get_variable(pkgconfig : 'sysusersdir')
    systemd_util_dir = dep_systemd.get_variable(pkgconfig : 'systemdutildir')
    cdata.set_quoted('SYSTEMD_UNIT_DIR', systemd_unit_dir)
    cdata.set_quoted('SYSTEMD_TMPFILES_DIR', systemd_tmpfiles_dir)
    cdata.set_quoted('SYSTEMD_SYSUSERS_DIR', systemd_sysusers_dir)
    cdata.set_quoted('SYSTEMD_UTIL_DIR', systemd_util_dir)

    with_systemd_reexec = get_option('with-systemd-reexec')
    if with_systemd_reexec == true
        cdata.set_quoted('HAVE_SYSTEMD_REEXEC', '1')
    endif

    with_vbox_restart = get_option('with-vbox-restart')
    if with_vbox_restart == true
        cdata.set_quoted('HAVE_VBOX_RESTART', '1')
    endif
endif

# Tracking files
cdata.set_quoted('USYSCONF_TRACK_DIR', path_vardir)
with_status_file = join_paths(path_vardir, 'status')
cdata.set_quoted('USYSCONF_STATUS_FILE', with_status_file)

with_log_dir = get_option('with-log-dir')
with_rewind_log_file = join_paths(with_log_dir, '@0@.rewind.log'.format(meson.project_name()))
with_log_file = join_paths(with_log_dir, '@0@.log'.format(meson.project_name()))
cdata.set_quoted('USYSCONF_LOG_DIR', with_log_dir)
cdata.set_quoted('USYSCONF_REWIND_LOG_FILE', with_rewind_log_file)
cdata.set_quoted('USYSCONF_LOG_FILE', with_log_file)

# Write config.h now
config_h = configure_file(
     configuration: cdata,
     output: 'config.h',
)
config_h_dir = include_directories('.')

# Get libuf built to provide convenience functions
libuf = subproject('libuf',
    default_options: [
        'with-tests=false',
        'with-static=true',
    ]
)

# Install manpage
subdir('man')

# Now go build the source
subdir('src')

report = [
    '    Build configuration:',
    '    ====================',
    '',
    '    prefix:                                 @0@'.format(path_prefix),
    '    datadir:                                @0@'.format(path_datadir),
    '    sysconfdir:                             @0@'.format(path_sysconfdir),
    '    mandir:                                 @0@'.format(path_mandir),
    '    sbindir:                                @0@'.format(path_sbindir),
    '    static binary:                          @0@'.format(with_static_binary),
    '    tracking directory:                     @0@'.format(path_vardir),
    '    status file:                            @0@'.format(with_status_file),
    '    rewind log file:                        @0@'.format(with_rewind_log_file),
    '    main log file:                          @0@'.format(with_log_file),
    '',
    '    Extra modules:',
    '    ==============',
    '',
    '    clr-boot-manager:                       @0@'.format(with_cbm),
    '    linux-driver-management:                @0@'.format(with_ldm),
    '    qol-assist:                             @0@'.format(with_qol),
    '    systemd support:                        @0@'.format(with_systemd),
    '    apparmor (aa-lsm-hook)                  @0@'.format(with_apparmor),
    '    mono-certs (cert-sync)                  @0@'.format(with_mono_certs),
]

if with_systemd == true
report += [
    '',
    '    systemd configuration:',
    '    ======================',
    '',
    '    systemd unit directory:                 @0@'.format(systemd_unit_dir),
    '    systemd tmpfiles directory:             @0@'.format(systemd_tmpfiles_dir),
    '    systemd sysusers directory:             @0@'.format(systemd_sysusers_dir),
    '    systemd binary directory:               @0@'.format(systemd_util_dir),
    '    allow systemd reexec:                   @0@'.format(with_systemd_reexec),
    '    support vboxdrv.service:                @0@'.format(with_vbox_restart),
]
endif

# Output some stuff to validate the build config
message('\n\n\n' + '\n'.join(report) + '\n\n')
