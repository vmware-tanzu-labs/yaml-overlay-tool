[Back to Table of contents](../index.md)  

# Environment variables

Added in Yot v0.5.0.

Environment variables can be set to eliminate the need to pass Yot CLI flags/parameters.  They can also be used to override configuration file settings as needed.

Your current Yot settings can be checked by running `yot env`.

The following table will display the available settings, their default values, and their available options.

| Environment Variable | Default Value | Options | Description | Version Added |
| --- | --- | --- | --- | --- |
| YOT_CONFIG_FILE | "" | any path to a Yot configuration file you like | Yot configuration file location. Corresponds to the `--config` CLI parameter. | v0.5.0 |
| YOT_INDENT_LEVEL | "2" | 2-9 | How much to indent the new YAML.  Corresponds to the `-I` or `--indent-level` CLI parameter. | v0.5.0 |
| YOT_LOG_LEVEL | "error" | critical, error, warning, notice, info, debug | What log level to run with.  Corresponds to the `-v` or `--log-level` CLI parameter. | v0.5.0 |
| YOT_OUTPUT_DIRECTORY | "./output" | any path you like | Path where you would like the new YAML files to be output. Corresponds to the `-o` or `--output-directory` CLI parameter. | v0.5.0 |
| YOT_OUTPUT_STYLE | "[normal]" | normal, tagged, doubleQuoted, singleQuoted, literal, folded, flow | Style of the new YAML file output. Corresponds to the `-S` or `--output-style` CLI parameter. | v0.5.0 |
| YOT_REMOVE_COMMENTS | "false" | false, true | Removes existing comments prior to performing overlays.  Corresponds to the `--remove-comments` CLI parameter. | v0.5.0 |
| YOT_STDOUT | "false" | false, true | Whether or not to output to `stdout`/standard out.  Corresponds to the `-s` or `--stdout` CLI parameter. | v0.5.0 |
| YOT_DEFAULT_ON_MISSING_ACTION | "ignore" | ignore, inject | Sets the default `onMissing` action, which is defaulted to `ignore`. | v0.5.0 |


[Back to Table of contents](../index.md)  
[Next Up: Command line interface usage and overview](commandUsage.md)