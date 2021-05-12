## Usage

```
yot --help

usage: yot [-h] [-d DEFAULT_VALUES_FILE] [-v VALUES_PATH] -i INSTRUCTION_FILE
           [-o OUTPUT_DIRECTORY] [-s] [-r] [-l LOG_FILE] [-V]

yot (YAML overlay tool) is a yaml overlay tool which allows for the templating
of overlay instruction data with jinja2, and the application of rendered
overlays "over the top" of a yaml file. yot only produces valid yaml documents
on output.

optional arguments:
  -h, --help            show this help message and exit
  -d DEFAULT_VALUES_FILE, --default-values-file DEFAULT_VALUES_FILE
                        Path to your default values file. If not set, you must
                        pass a values file of "defaults.yaml" or
                        "defaults.yml" within a path from the "-v" option.
                        Takes multiple default values files in case you would
                        like to separate out some of the values. After the
                        first default values file, each subsequent file passed
                        with "-d" will be merged with the values from the
                        first. If a "defaults.yaml" or "defaults.yml" file is
                        discovered in one of your "-v" paths, it will be
                        merged with these values last.
  -v VALUES_PATH, --values-path VALUES_PATH
                        Values file path. May be a path to a file or directory
                        containing value files ending in either .yml or .yaml.
                        This option can be provided multiple times as
                        required. A file named "defaults.yaml" or
                        "defaults.yml" is required within the path(s) if not
                        using the "-d" option, and you may have only 1 default
                        value file in that scenario. Additional values files
                        are merged over the defaults.yaml file values. Each
                        values file is treated as a unique "site" and will
                        render your instructions differently based on its
                        values
  -i INSTRUCTION_FILE, --instruction-file INSTRUCTION_FILE
                        Instruction file path. Defaults to ./instructions.yaml
  -o OUTPUT_DIRECTORY, --output-directory OUTPUT_DIRECTORY
                        Path to directory to write the overlayed yaml files
                        to. If value files were supplied in addition to a
                        defaults.yaml/.yml then the rendered templates will
                        land in <output dir>/<addl value file name>.
  -s, --stdout          Render output to stdout. Templated instructions files
                        will still be output to the "--output-directory."
  -r, --dump-rendered-instructions
                        If using a templated instructions file, you can dump
                        the rendered instructions to stdout to allow for
                        reviewing how they were rendered prior to a full run
                        of yot. Equivalent to a dry-run. Exits with return
                        code 0 prior to processing instructions
  -l LOG_FILE, --log-file LOG_FILE
                        debug log file output path. Where to output the log
                        to. Defaults to ./yot.log
  -V, --log-verbosity   Log level verbosity (-V for critical, -VV for error,
                        -VVV for warning, -VVVV for info, -VVVVV for debug).
                        Defaults to -VVV (warning and above)
```

### Example Usage

```bash
# with no templating
./yot \
    -i ./examples/static/instructions.yaml \
    -o ./out

# with templating of overlays from a directory path
./yot \
    -i ./examples/templated/instructions.yaml \
    -v ./examples/values \
    -o ./out

# more to come soon
```

# From old Quickstart
`yot` also includes a templating feature by making use of the popular [Jinja2](https://jinja.palletsprojects.com/en/master/templates/) templating language to allow for templated overlay values which are rendered into memory and processed at run-time.  Use of the templating engine is completely optional, and instruction files can be static yaml documents.  When using the templating engine, anything within the instructions file can be templated, but keep in mind the document must template into a valid yaml document.  The template feature is useful if you are managing multi-environment yaml configurations.  Values files are also treated as templates, and can contain jinja2 content.  This is most useful in a scenario with a lot of values, where you would like to organize them into separate files and use the `{% include 'addl_values.yaml' %}` [Jinja2](https://jinja.palletsprojects.com/en/master/templates/#include) tag.  This is also extremely useful for large instruction files.