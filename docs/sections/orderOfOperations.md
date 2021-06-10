[Back to Table of contents](../documentation.md)  


## Order of operations/Processing Order

### 1a. Dynamic instructions (Templating) - Not currently implemented in v0.1.0

When one or more values files have been passed to `yot` a `defaults.yaml` or `defaults.yml` file must exist within that path, when passing a directory with the `-v` option.  

Alternatively, a single or multiple default value files can be passed with the `-d` option.
If you only have one site and a desire to template, then only provide a default values file with `-d` or a file named `defaults.yaml` or `defaults.yml` within a directory passed with the `-v` option.

If you pass multiple default values files with the `-d` option, the first file passed serves as the base values.  Each additional default values files passed with the `-d` option are merged over the top of the base values sequentially in the order they were passed to `yot`.  Additionally, if a `defaults.yaml` or `defaults.yml` file was also present in a path passed with the `-v` option, those values will be merged over the base values last.

Any additional values files (files passed with `-v`) are treated as additional sites (or site values), and the values contained in those files are merged over the top of the `defaults.yaml` file's values.  The base filename of the additional values files are used to output the modified yaml files into.  

#### Templating Example - Not currently implemented in v0.1.0

Consider the following example:

```bash
$ tree
├── instructions.yaml
└── values
    ├── bar
    ├── defaults.yaml
    ├── foo.yaml
    └── test.yml

./yot -v values -i instructions.yaml -o ./output
```

`yot` only takes values files with the file extension of `.yml` or `.yaml`, therefore the file `bar` above would not be read.  The values contained in values/defaults.yaml would be read and then the instructions file would be rendered by applying the values of foo.yaml over defaults.yaml, and then render the instructions.yaml template into memory for processing.  Then the test.yml will merge over defaults.yaml and render the instructions.yaml again, producing a total of 2 unique instruction sets.  When processing is complete, the modified yaml documents will be placed in `./output/yaml_files/foo` and `./output/yaml_files/test` respectively.  
> **NOTE:**  Please be careful about naming your values files, as a values file called test.yaml and test.yml would dump both of their rendered contents into ./output/test

### 1b. Static instructions (No templating) - v0.1.0 behavior

If no values files were passed, then the instructions file skips the templating mechanism and is loaded directly, as long as the document is valid yaml.

Consider the following example:
```bash
yot -i instructions.yaml -o ./output
```

`yot` will read the instructions.yaml file and any outputs will be placed directly in ./output/ .

### 2. Processing instructions

After the instructions have either been rendered or read, they are processed.

Processing begins at the instruction's `commonOverlays` if they exist in the instruction set.  `yot` combines the `commonOverlays` with the `yamlFiles.path.overlays` first if they exist.  If `yamlFiles.path.overlays` do not exist, `yot` combines `commonOverlays` with `yamlFiles.path.documents.path.overlays`.  In both of these scenarios, the `commonOverlays` will always be applied prior to the more granular overlays.

If no `commonOverlays` have been defined, processing starts at the first path of the `yamlFiles` array and processes one YAML file `path` at a time sequentially.  Within each YAML file `path`, `overlays` are processed first if set, followed by the items within the `documents` key, which applies overlays to specific documents within a multi-yaml document yaml file. A single yaml document in a YAML file could still be referred to by `path: 0` from the `documents` key if desired.

### 3. Output

If no templating was performed, then output is sent to the output directory specified at runtime, or if it was not provided, the default path of `./output/`.

If values files in addition to a defaults.yaml were passed at runtime, then the output for each values file will be `< output directory >/< value file basename >` (not implemented in v0.1.0 of `yot`).

If only a defaults.yaml value file was passed, then the output will be placed in the < output directory specified at runtime >/, or if it was not provided, the default path of `./output/`.


[Back to Table of contents](../documentation.md)  
[Next Up: Output directory structure](outputDirStructure.md)