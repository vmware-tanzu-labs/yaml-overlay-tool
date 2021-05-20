[Back to Table of Contents](../documentation.md)  

# Comment Preservation and Injection

Due to the method YAML Overlay Tool (`Yot`) uses to operate on YAML nodes, it can preserve header (above a piece of data), footer (below a piece of data), and line (on the same line as a piece of data) comments. This is useful if you want to retain comment data for informational purposes, or if you require comments to be preserved for another tool to consume.`Yot` is also unique because it can inject comments into YAML files.

## Comment Preservation

By default (v0.1.0), Yot preserves all existing comments that are operated on within a YAML file.  A future release of `Yot` will contain a configuration file with options on which comments (header, footer, or line) to retain will be added.  Until then, all comments are preserved.


## Comment Injection

You can inject comments into YAML files in the same way that `Yot` preserves existing comments in YAML files.

WIP


[Back to Table of Contents](../documentation.md)  
[Next Up: Order of Operations/Processing](orderOfOperations.md)
