[Back to Table of Contents](../documentation.md)  

# Comment Preservation and Injection

Due to the nature of how YAML Overlay Tool operates on YAML nodes, the tool has the ability to preserve header (above a piece of data), footer (below a piece of data), and line (on the same line as a piece of data) comments.  This is extremely useful if you wanted to retain comment data for informational purposes, or require comments to be preserved for some other tool to consume.  YAML Overlay Tool is also unique, because the tool can also inject comments into YAML files.

## Comment Preservation

By default (v0.1.0), YAML Overlay Tool preserves all existing comments within a YAML file that has been operated on.  In a future version of `yot`, a configuration file with options on which comments (header, footer, or line) to retain will be added.  Until those configurable options are added, all comments will be preserved.


## Comment Injection

Comments can be injected into YAML files similarly to how `yot` preserves existing comments in YAML files.

WIP


[Back to Table of Contents](../documentation.md)  
[Next Up: Order of Operations/Processing](orderOfOperations.md)