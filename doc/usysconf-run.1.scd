usysconf-run(1)

# NAME

usysconf - run command

# SYNOPSIS

_usysconf_ run [options...]

# OPTIONS
*-c*, *--chroot*
	Specify that the configuration is being run in a chrooted environment

*-d*, *--debug*
	Debug output

*-f*, *--force*
	Force run the configuration regardless if it should be skipped

*-l*, *--live*
	Specify that the configuration is being run from a live medium

*-n*, *--names* [name...]
	Specify the name of the files to be run in a comma separated list

*--norun*
	Test the configuration files without executing the specified binaries
	and arguments

# SEE ALSO

*usysconf*(1) *usysconf-config*(5)

# AUTHORS

Maintained by the Solus team <https://getsol.us>.  For more information about
usysconf development, see <https://github.com/getsolus/usysconf>.
