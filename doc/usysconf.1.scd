usysconf(1)

# NAME

usysconf - Universal system configuration

# SYNOPSIS

_usysconf_ [options...] [command]

# OPTIONS

*-h*, *--help*
	Prints the usage information and exits.

*-v*, *--version*
	Prints the installed version of usysconf and exits.

# DESCRIPTION

usysconf is a system configuration agent designed for integration within the
software update process. The main job for usysconf is to run OS configuration
tasks within a well defined centralised location.

usysconf will run a variety of triggers to ensure system consistency and health.

For information on the how to run configuration files, see *usysconf-run*(1).

# CONFIGURATION

usysconf searches for config files in the following location, in this order:

. /usr/share/defaults/usysconf.d (or whatever path is set for USRDIR)
. /etc/usysconf.d (or whatever path is set for SYSDIR)

For information on the config file format, see *usysconf-config*(5).

# SEE ALSO

*usysconf-run*(1) *usysconf-config*(5)

# AUTHORS

Maintained by the Solus team <https://getsol.us>.  For more information about
usysconf development, see <https://github.com/getsolus/usysconf>.
