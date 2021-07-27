[Back to Table of contents](../index.md)  


## Example CLI usage

<!-- @import "[TOC]" {cmd="toc" depthFrom=3 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Use with an Instructions file](#use-with-an-instructions-file)
  - [Specify an output path](#specify-an-output-path)
  - [Direct output to stdout](#direct-output-to-stdout)
  - [Pipe output to kubectl](#pipe-output-to-kubectl)
  - [Manipulating the output style](#manipulating-the-output-style)
  - [Manipulating the output indentation level](#manipulating-the-output-indentation-level)
  - [Remove source YAML file comments prior to overlayment](#remove-source-yaml-file-comments-prior-to-overlayment)
  - [Provide variable values to a templated Instructions file](#provide-variable-values-to-a-templated-instructions-file)
  - [Obtain logging or debug information](#obtain-logging-or-debug-information)
- [Use without an Instructions file](#use-without-an-instructions-file)
  - [One-off overlay example](#one-off-overlay-example)
  - [One-off overlay example from stdin](#one-off-overlay-example-from-stdin)
  - [One-off overlay example in conjunction with an Instructions file](#one-off-overlay-example-in-conjunction-with-an-instructions-file)
- [Check YAML Overlay Tool's Configuration Setting](#check-yaml-overlay-tools-configuration-setting)

<!-- /code_chunk_output -->


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

Yot allows [templating within the Instructions file](instructionsFileTemplating.md).  To use templating, values must be supplied at run-time for the particular templating engine to be able to render the instructions into a usable Instructions set.

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

As of v0.5.0, an overlay can be applied purely from the CLI, and can be processed with or in combination with an instructions file.  [Overlay qualifiers](overlayQualifiers.md) are not supported via the CLI in v0.5.0.

To process an overlay without an instructions file, a few parameters are required:  
    * **Query:** `q`, `query`
    * **Action:** `a`, `action` ( default is merge)
    * **Value:** `x`, `value` (not required for action delete)
    * **Path:** `p`, `path`

#### One-off overlay example

```bash
yot -q metadata.labels -x "{app.kubernetes.io/owner: Jeff Smith}" -p /path/to/source/yaml/file.yaml -o /tmp/new
```

#### One-off overlay example from stdin

```bash
cat /path/to/yaml/files/*.yaml | yot -q metadata.labels -x "{app.kubernetes.io/owner: Jeff Smith}" -p - -o /tmp/new
```

#### One-off overlay example in conjunction with an Instructions file

An additional overlay can be added to be processed in addition to an Instructions file.  CLI based overlays are processed as if they were `commonOverlays`, and when specified in addition to an Instructions file they are always processed as the ***last*** `commonOverlay`.  

```bash
yot -i /path/to/my/instructions/file.yaml -q metadata.labels -x "{app.kubernetes.io/owner: Jeff Smith}" -a merge -p /path/to/source/yaml/file.yaml -o /tmp/new
```

### Check YAML Overlay Tool's Configuration Setting

An option to override default settings was added in v0.5.0.  This can be done with either [environment variables](envVars.md) or a [configuration file](configFile.md).  To check Yot's current settings the following command can be run:

```bash
yot env
```


[Back to Table of contents](../index.md)  
[Next Up: Instructions file introduction](instructionsFileIntro.md)
