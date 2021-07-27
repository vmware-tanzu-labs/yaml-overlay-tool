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

| Setting Key | Default Value | Options | Description | Version Added |
| --- | --- | --- | --- | --- |
| indentLevel | 2 | 2-9 | How much to indent the new YAML.  Corresponds to the `-I` or `--indent-level` CLI parameter. | v0.5.0 |
| logLevel | "error" | critical, error, warning, notice, info, debug | What log level to run with.  Corresponds to the `-v` or `--log-level` CLI parameter. | v0.5.0 |
| outputDirectory | ./output | any path you like | Path where you would like the new YAML files to be output. Corresponds to the `-o` or `--output-directory` CLI parameter. | v0.5.0 |
| outputStyle | [normal] | normal, tagged, doubleQuoted, singleQuoted, literal, folded, flow | Style of the new YAML file output. Corresponds to the `-S` or `--output-style` CLI parameter. | v0.5.0 |
| removeComments | false | false, true | Removes existing comments prior to performing overlays.  Corresponds to the `--remove-comments` CLI parameter. | v0.5.0 |
| stdout | false | false, true | Whether or not to output to `stdout`/standard out.  Corresponds to the `-s` or `--stdout` CLI parameter. | v0.5.0 |
| defaultOnMissingAction | ignore | ignore, inject | Sets the default `onMissing` action, which is defaulted to `ignore`. | v0.5.0 |


[Back to Table of contents](../index.md)  
[Next Up: Environment variables](envVars.md)