[Back to Table of contents](../index.md)


## yot

yot (YAML Overlay Tool) is a YAML overlay tool.

### Synopsis

yot (YAML Overlay Tool) is a YAML overlay tool which uses a YAML schema to 
	define overlay operations on a set of YAML documents. yot only produces valid YAML 
	documents on output, and can preserve and inject comments.

```
yot [flags]
```

### Examples

```
yot -i instructions.yaml -o /tmp/output
```

### Options

```
Available Commands:
  completion  Generate shell auto-completion scripts
  help        Help about any command

Flags:
  -a, --action actions.Action        Action to take when the JSONPath 'query' is found with a YAML document. 
                                     Can be one of combine, delete, merge, or replace. Typically used for one-off overlays from the CLI.
  -h, --help                         help for yot
  -I, --indent-level int             Number of spaces to be used for indenting YAML output (min: 2, max: 9) (default 2)
  -i, --instructions string          Path to the instructions file (required)
  -v, --log-level logLevel[=debug]   Log-level to display to stdout, one of: 
                                        CRITICAL: {"critical", "crit", "c"},
                                        ERROR:    {"error", "err", "e"},
                                        WARNING:  {"warning", "warn", "w"},
                                        NOTICE:   {"notice", "note", "n"},
                                        INFO:     {"info", "i"},
                                        DEBUG:    {"debug", "d", "verbose", "v"} * used if no argument is provided
                                      (default error)
  -o, --output-directory string      Path to a directory for writing the YAML files which were operated on by Yot (default "./output")
  -S, --output-style style           Style to be used for rendering final documents.
                                     Multiple values can be provided to achieve the desired result, valid values are:
                                        NORMAL:       {"normal", "n"},
                                        TAGGED:       {"tagged", "tag", "t"},
                                        DOUBLEQUOTED: {"doubleQuoted", "doubleQuote", "double", "dq"},
                                        SINGLEQUOTED: {"singleQuoted", "singleQoute", "single", "sq"},
                                        LITERAL:      {"literal", "l"},
                                        FOLDED:       {"folded", "fold", "fo"},
                                        FLOW:         {"flow", "fl"}
                                      (default [normal])
  -p, --path string                  YAML files outside of the instructions file to process. Using a 'path' of "-" will read from stdin. 
                                     Typically used for one-off overlays from the CLI.
  -q, --query overlays.Queries       JSONPath query or JSONPath fully-qualified (dot-notation) path you would like to manipulate.
                                     these will get added to the common overlays in conjunction with 'action' and 'value'
      --remove-comments              Remove all comments from the source YAML files prior to overlayment
  -s, --stdout                       Output YAML files which were operated on by Yot to stdout
  -x, --value string                 Desired 'value' to take 'action' with if 'query' is found within the YAML document. 
                                     Typically used for one-off overlays from the CLI.
  -f, --values-file stringArray      Path to a values file for use with templating an instructions file.
                                     Takes multiple values files in case you would like to better organize the values. 
                                     Each subsequent file passed with -f will be merged over the values 
                                     from the previous. Values are applied to your instructions file when using templating.
      --version                      version for yot
```

[Back to Table of contents](../index.md)  
[Next Up: Example CLI usage](exampleUsage.md)
