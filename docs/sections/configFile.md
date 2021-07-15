[Back to Table of contents](../index.md)  


## Configuration file

Added in Yot v0.5.0.

A configuration file can be placed in 3 default locations on your filesystem, but can also be overridden by the environment variable or the runtime flag.  

This is the order of precedence for the configuration file:

1. `/etc/yot/yot.config`  
1. `$HOME/.yot/yot.config`  
1. `./yot.config` i.e. the current working directory  
1. `YOT_CONFIG_FILE` environment variable
1. `--config` CLI option


Item five will always override all other configuration file locations.  Each configuration file settings are not be merged together.

If Yot is behaving in an unexpected behavior, you can check your current configuration settings by running `yot env`.


### Configuration file settings

The following table will display the available settings, their default values, and their available options.

| Setting Key | Default Value | Options |
| --- | --- | --- |
| indentLevel | 2 | 2-9 |
| logLevel | "error" | critical, error, warning, notice, info, debug |
| outputDirectory | ./output | any path you like |
| outputStyle | [normal] | normal, tagged, doubleQuoted, singleQuoted, literal, folded, flow |
| removeComments | false | false, true |
| stdout | false | false, true |
| defaultOnMissingAction | ignore | ignore, inject |


[Back to Table of contents](../index.md)  
[Next Up: Environment variables](envVars.md)