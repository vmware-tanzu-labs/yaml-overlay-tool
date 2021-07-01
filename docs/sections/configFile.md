[Back to Table of contents](../documentation.md)  


## Configuration file

Added in Yot v0.5.0.

A configuration file can be placed in 3 locations on your filesystem:

1. `/etc/yot/yot.config`  
2. `$HOME/.yot/yot.config`  
3. `./yot.config` where this is the current working directory  

Item three will always override all other configuration files, however environment variables and `yot` command flags can still override the configuration file's settings.

If Yot is behaving in an unexpected behavior, you can check your current configuration settings by running `yot env`.


### Configuration file settings

The following table will display the available settings, their default values, and their available options.

| Setting Key | Default Value | Options |
| --- | --- | --- |
| indentLevel | 2 | 2-9 |
| logLevel | error | "error" | critical, error, warning, notice, info, debug |
| outputDirectory | ./output | any path you like |
| outputStyle | [normal] | normal, tagged, doubleQuoted, singleQuoted, literal, folded, flow |
| removeComments | false | false, true |
| stdout | false | false, true |
| defaultOnMissingAction | ignore | ignore, inject |


[Back to Table of contents](../documentation.md)  
[Next Up: Environment variables](envVars.md)