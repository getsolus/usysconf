usysconf-config(5)

# NAME

usysconf = configuration files

# DESCRIPTION

A usysconf configuration file is written in TOML and should be saved in either
the system or user specified directories.

The configuration includes the following sections:

*description* 
	This provides context of what this configuration file should do.

*[[bins]]*
	This contains the information needed to run the binaries.

	*task*
		The text that will print during run.

	*bin*
		The binary that will be run.

	*args*
		A comma separated list of arguments that will run with the bin.

	*[replace]*
		Special section to allow the run the same bin multiple times
		with different paths.
		
		*paths*
			A comma separated list of paths, supports glob.

		*exclude*
			A comma separated list of paths to not run from paths,
			supports glob.

*[skip]*
	This contains the information for when the configuration should be
	skipped.

	*chroot*
		Pass true if the configuration will be skipped in a chrooted
		environment
	
	*live_media*
		Pass true if the configuration will be skipped in a live medium
		environment

	*paths*
		If the path exists in the comma separated list then the 
		configuration will be skipped, supports glob.

*[check]*
	This contains the information to check if a path exists, and if it does
	not then the configuration will be skipped.

	*paths*
		If the path does not exists in the comma separated list then 
		the configuration will be skipped, supports glob.

*env*
	Set environment variable for during the execution of the configuration.

*remove*
	Remove specified paths and files.
		
	*paths*
		A comma separated list of paths, supports glob.

	*exclude*
		A comma separated list of paths to remove, supports glob.

# SEE ALSO

*usysconf*(1) *usysconf-run*(1)

# AUTHORS

Maintained by the Solus team <https://getsol.us>.  For more information about
usysconf development, see <https://github.com/getsolus/usysconf>.
