[Back to Table of contents](../index.md)  


## Order of operations/Processing Order

### 1a. Dynamic instructions (Templating)

Go templating was implemented in v0.4.0, and Sprig function support for Go templating was added in v0.5.0.

When one or more values files have been passed to `yot`  

A single or multiple values files can be passed with the `-f` or `--values-file` option.

If you pass multiple default values files, the first file passed serves as the base values.  Each additional values files passed with the `-f` option are merged over the top of the base values sequentially in the order they were passed to `yot`.    


### 1b. Static instructions (No templating) - v0.1.0 behavior

If no values files were passed, then the instructions file skips the templating mechanism and is loaded directly, as long as the document is valid YAML.

Consider the following example:
```bash
yot -i instructions.yaml -o ./output
```

`yot` will read the instructions.yaml file and any outputs will be placed directly in ./output/ .

### 2. Processing instructions

After the instructions have either been rendered or read, they are processed.

Processing begins at the instruction's `commonOverlays` if they exist in the instruction set.  `yot` processes the `commonOverlays` first if they exist.  

If no `commonOverlays` have been defined, processing starts at the first path of the `yamlFiles` array and processes one YAML file `path` at a time sequentially.  Within each YAML file `path`, `overlays` are processed first if set, followed by the items within the `documents` key, which applies overlays to specific documents within a multi-yaml document YAML file. A single YAML document in a YAML file could still be referred to by `path: 0` from the `documents` key if desired.

### 3. Output

Output is sent to the output directory specified at runtime, or if it was not provided, the default path of `./output/`.


[Back to Table of contents](../index.md)  
[Next Up: Output directory structure](outputDirStructure.md)