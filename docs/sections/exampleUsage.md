[Back to Table of contents](../documentation.md)  

## Example CLI usage

### Use with an Instructions file

YAML Overlay Tool's main intended use is with an [Instructions file](instructionsFileIntro.md).  
To use Yot with an Instructions file, simply provide the `-i` parameter followed by a path to the instructions.  Yot's default output path is `./output`.

```bash
yot -i < instructions file >
```


#### Specify an output path

While the default output path is `./output`, a Yot user may specify an output path on the filesystem.

```bash
yot -i < instructions file > -o < output path >
```


#### Direct output to stdout

Output is not exclusive to the filesystem, and may be sent directly to standard out.

```bash
yot -i < instructions file > -s
```


#### Pipe output to kubectl

Using standard out to direct overlayed Kubernetes manifests into a cluster is simple.

```bash
yot -i < instructions file > -s | kubectl apply -f -
```


#### Manipulating the output style

By default, Yot outputs your YAML documents in the same style they started, and any changes going in to the documents will maintain their original style.  However, Yot has the ability to manipulate the final output style of your YAML documents with the `-S` parameter followed by the desired style option.  
There are seven styles to choose from, all of which are considered valid YAML:
  * **Normal:** `n`, `normal`, `NORMAL`
  * **Tagged:** `t`, `tagged`, `tag`, `TAGGED`
  * **Double-quoted:** `dq`, `double`, `doubleQuote`, `doubleQuoted`, `DOUBLEQUOTED`
  * **Single-quoted:** `sq`, `single`, `singleQuote`, `singleQuoted`, `SINGLEQUOTED`
  * **Literal:** `l`, `literal`, `LITERAL`
  * **Folded:** `fo`, `fold`, `folded`, `FOLDED`
  * **Flow:** `fl`, `flow`, `FLOW`

```bash
yot -i < instructions file > -S < style option >
```


#### Manipulating the output indentation level

By default, Yot outputs your YAML documents with two spaces for all indentation.  To modify this behavior, use the `-I` parameter followed by an integer which represents the number of spaces to indent.  The minimum number of spaces for indentation is 2, while the maximum is 9.

```bash
# indent 4 spaces in output
yot -i < instructions file > -I 4
```

#### Remove source YAML file comments prior to overlayment

By default, Yot preserves all comments within a YAML document.  If a Yot user would like to remove the original comments within a YAML document prior to performing overlays, pass the `--remove-comments` parameter.  Comments may still be injected from the Instructions file, but the original comments will be removed first.

```bash
yot -i < instructions file > --remove-comments
```


#### Provide variable values to a templated Instructions file

Yot allows [templating within the Instructions file](instructionsFileTemplating.md).  To use templating, values must be supplied at run time for the particular templating engine to be able to render the instructions into a usable Instructions set.

Pass variable values into Yot by using the `-f` paramater followed by the path to a YAML file containing values.  This parameter may be passed multiple times, where any overlapping values will be overridden from the newest instantiation of the `-f` parameter.

```bash
yot -i < instructions file > -f /my/template/values1.yaml -f /my/template/values2.yaml
```


#### Obtain logging or debug information

Yot has a logging facility to assist with debugging, and to determine what might be causing something to not work in the way the user expected.  
To see the fully verbose logs when running Yot, pass only the `-v` parameter which is by default the debug logs.  

```bash
# default debug logs
yot -i < instructions file > -v
```

Additionally, the `-v` parameter can be followed by an equal sign `=` and a desired log-level: `-v=error`.

There are six available log-levels that can be used at run-time:

  * **Critical:** `c`, `crit`, `critical`, `CRITICAL`
  * **Error:** `e`, `err`, `error`, `ERROR`
  * **Warning:** `w`, `warn`, `warning`, `WARNING`
  * **Notice:** `n`, `note`, `notice`, `NOTICE`
  * **Info:** `i`, `info`, `INFO`
  * **Debug:** `d`, `debug`, `DEBUG`, `v`, `verbose`

```bash
# warning logs
yot -i < instructions file > -v=WARNING
```


### Use without an Instructions file

Yot may also be used without an Instructions file, and behaves as a [`commonOverlay`](instructionsFileSpec.md#top-level-commonoverlays-keys) does (added in v0.5.0).  


#### One-off overlays

An overlay can be applied purely from the CLI, and can be performed with or without an instructions file.  

There are a WIP


[Back to Table of contents](../documentation.md)  
[Next Up: Instructions file introduction](instructionsFileIntro.md)