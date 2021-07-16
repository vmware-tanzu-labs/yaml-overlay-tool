[Back to Table of contents](../index.md)  


## Order of operations/processing order

### Processing overlays with a Yot instructions file
#### 1a. Dynamic instructions (Templating)

Go templating was implemented in v0.4.0, and Sprig function support for Go templating was added in v0.5.0.

The `-f` or `--values-file` option may be used one or many times.

If you pass multiple values files, the first file passed in with `-f` or `--values-file` serves as the base set of values.  Each additional values file passed is merged over the top of the base values sequentially in the order they were passed to `yot` producing a final set of values to be applied to the templated instructions.  If a value is repeated in an additional values file, it will be overridden with the latest instance of that value winning.  


#### 1b. Static instructions (No templating) - v0.1.0 behavior

If no values files were passed, then the instructions file skips the templating mechanism and is loaded directly, as long as the document is valid YAML.

Consider the following example:
```bash
yot -i instructions.yaml -o ./output
```

`yot` will read the instructions.yaml file and any outputs will be placed directly in ./output/ .

#### 2. Processing instructions

After the instructions have either been rendered or read, they are processed.

Yot provides three levels of specificity, being `commonOverlays`, `yamlFiles[].overlays`, and `yamlFiles[]documents.overlays`.

Processing begins at the instruction's `commonOverlays` if they exist in the instruction set.  `yot` processes the `commonOverlays` first if they exist by prepending them to each instance of yamlFiles[].overlays

If no `commonOverlays` have been defined, processing starts at the first path of the `yamlFiles` array and processes one YAML file `path` at a time sequentially.  Within each YAML file `path`, `overlays` are processed first if set, followed by the items within the `documents` key, which applies overlays to specific documents within a multi-yaml document YAML file. A single YAML document in a YAML file could still be referred to by `path: 0` from the `documents` key if desired.

#### 3. Output

Output is sent to the output directory specified at runtime, or if it was not provided, the default path of `./output/`, or to the location set within the `yot.config` or the `YOT_OUTPUT_DIRECTORY` environment variable.


### Processing overlays with the CLI (no instructions file)

When performing one-off overlays via the CLI, Yot treats the `-q` or `--query` parameter as part of a `commonOverlay`.  Based on the `-a` or `--action`, and `-x` or `--value` parameters a `commonOverlay` will be constructed within a in-memory instruction file, where `-p` or `--path` is loaded into the `yamlFiles` array.  

Templating is not allowed in a CLI-only or one-off overlay operation.  Processing occurs as described in [Processing overlays with a Yot instructions file](#processing-overlays-with-a-yot-instructions-file).


[Back to Table of contents](../index.md)  
[Next Up: Output directory structure](outputDirStructure.md)