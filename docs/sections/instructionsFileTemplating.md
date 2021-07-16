[Back to Table of contents](../index.md) 

## Using templating within instructions files

### Why use templating?

Templating creates reusability within your overlays.  

### Supported templating languages

* As of v0.4.0 YAML Overlay Tool supports [Go Templating](https://golang.org/pkg/text/template/).  
* As of v0.5.0 YAML Overlay Tool supports Go Templating with [Sprig functions](https://masterminds.github.io/sprig/).  

Future versions of Yot ***may*** add additional configurable templating language support.


### How templating variables are handled

Templating variables are passed into Yot at run-time by the `-f` or `--values-file` parameter.  This parameter can be passed more than once as needed.  If passing more than one values-file, the initial file will act as the base values and each subsequent file will merge with these values.  ***Ordering does matter*** if there are any overlapping values.

The values file is expected to be a YAML file containing values you would like to insert into your templated instructions file.

Prior to processing your instructions file, if the `-f` or `--values-file` parameter is passed, the templated instructions will be rendered.  If the template fails to produce valid YAML, the instructions will fail to be read and an error will occur.

See the [Example CLI usage](exampleUsage.md#provide-variable-values-to-a-templated-instructions-file) page for an example.


[Back to Table of contents](../index.md)  
[Next Up: Order of operations/processing order](orderOfOperations.md)