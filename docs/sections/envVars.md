[Back to Table of contents](../index.md)  

# Environment variables

Added in Yot v0.5.0.

Environment variables can be set to eliminate the need to pass Yot CLI flags/parameters.  They can also be used to override configuration file settings as needed.

Your current Yot settings can be checked by running `yot env`.

The following table will display the available settings, their default values, and their available options.

| Environment Variable | Default Value | Options |
| --- | --- | --- |
| YOT_CONFIG_FILE | "" | any path to a file you like |
| YOT_INDENT_LEVEL | "2" | 2-9 |
| YOT_LOG_LEVEL | "error" | critical, error, warning, notice, info, debug |
| YOT_OUTPUT_DIRECTORY | "./output" | any path you like |
| YOT_OUTPUT_STYLE | "[normal]" | normal, tagged, doubleQuoted, singleQuoted, literal, folded, flow |
| YOT_REMOVE_COMMENTS | "false" | false, true |
| YOT_STDOUT | "false" | false, true |
| YOT_DEFAULT_ON_MISSING_ACTION | "ignore" | ignore, inject |


[Back to Table of contents](../index.md)  
[Next Up: Command line interface usage and overview](commandUsage.md)