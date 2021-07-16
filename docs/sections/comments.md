[Back to Table of contents](../index.md)  

# Comment preservation and injection

Due to the method YAML Overlay Tool (Yot) uses to operate on YAML nodes, the tool can preserve header (above a piece of data), footer (below a piece of data), and line (on the same line as a piece of data) comments.  This is extremely useful if you want to retain comment data for informational purposes, or require comments to be preserved for some other tool to consume.  YAML Overlay Tool (Yot) is also unique because it can inject comments into YAML files as well.


## Comment preservation

By default, Yot preserves all existing comments within any YAML file that Yot has operated on.  In a future release, a configuration file with options on which comments (header, footer, or line) to retain will be added.

In release (v0.5.0), a command-line option `--remove-comments` will allow a user to remove all comments from the original YAML files being operated on.  This can also be set via an [environment variable](envVars.md) or a Yot [configuration file](configFile.md) setting.

You can inject comments into YAML files by simply adding a comment above (head comments), below (foot comments), or on the same line (line comments) as the data within the `value` key of any overlay within the [instructions file](instructionsFileSpec.md).  

>**NOTE:** Due to [some minor bugs within Go's yaml.v3 library](https://github.com/go-yaml/yaml/issues/610), head comments in a map/dictionary do not always apply where they should, and will be addressed in a future version of Yot.  **However,** ***line comments*** **can be** reliably **injected** today.  Head and foot comments are currently considered experimental.


[Back to Table of contents](../index.md)  
[Next Up: Using templating within instructions files](instructionsFileTemplating.md)