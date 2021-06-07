[Back to Table of Contents](../documentation.md)

# Instructions file YAML schema

What is an instructions file?  
The instructions file is the instructions for manipulating a set of YAML files and documents.  It is a YAML document that contains an array of paths to YAML documents to perform operations.  Each instructions file starts with the optional key `commonOverlays`, which is an array, along with a required key `yamlFiles` which is also an array.  


## Top-level commonOverlays Keys

The `commonOverlays` key is optional. It provides overlays to apply to every item in the `yamlFiles` array, also referred to as `yamlFiles[*].path` or `yamlFiles[*].documents.path` as defined in the instructions. Each array item in the `commonOverlays` key is treated as a dictionary/map with the following top-level keys:

| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change you are performing, and used for either on-screen output or self-documentation. | None | string |
| query | yes | JSONPath query or JSONPath fully-qualified (dot-notation) path to value you would like to manipulate. If the `query` is not a fully-qualified path (such as a.b.c.d) and returns no matches, you need to specify the `onMissing` key (i.e. metadata.labels VS metadata.fake.*). | None | string or list/array |
| value | yes | The desired value to take action with if `query` is found. This is not required when using the `delete` action, as we are removing a value altogether. | None | int, bool, str, dict/map, list/array |
| action | yes | The action to take when the JSONPath expression is found in the YAML document. Can be one of `combine`, `delete`, `merge`, or `replace`.  See [Actions](actions.md) for details. | None | string |
| onMissing.action | no | Instructions for what to do if the JSONPath `query` is not found. Can be `ignore` or `inject`. Only applies to the actions `merge` and `replace`| `ignore` | string |
| onMissing.injectPath | no | If your JSONPath expression is not a fully-qualified path (dot-notation) then an `injectPath` is required to qualify your action. Only applies if `onMissing: {'action': 'inject'}` is set. This is the path or list/array of paths to inject the value if your JSONPath expression was not found in the YAML document. | None | string or list/array |
| documentQuery | no | A qualifier to refine which documents the common overlay is applied to.  If not set, the overlay applies to all files in `yamlFiles`.  See [Qualifiers](#qualifiers) for more details. | None | dictionary/map |

## Top-Level yamlFiles Keys

Each list item in the `yamlFiles` key is treated as a dictionary/map with the following top-level keys:

| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change/file you are performing, and used only for on-screen output or self-documentation. | None | string |
| path | yes | A fully qualified path to the YAML file to modify, or a path relative to the location of the instructions file. Can be a path to a YAML file or a directory containing YAML files. | None | string |
| overlays | no | List/array of overlay operations to apply. If your YAML file contains multiple documents separated by `---`, then this would apply to every YAML document first, unless a qualifier or combination of qualifiers `documentQuery` and `documentIndex` are provided.  If you need to apply overlays only to a specific YAML document in a multi-document YAML file, then see the `documents` key. See [overlays keys](#overlays-keys) for available dictionary/map keys. | None | list/array of dictionaries |
| documents | no | List/array of overlay operations to apply to a multi-document YAML file.  When each document from a multi-document YAML file is loaded, an overlay can be applied by addressing the document by its index.  See [documents keys](#documents-keys) for available dictionary/map keys. | None | list/array of dictionaries/maps |
| outputPath | no | alter the name of the output path. if a filename is given this will also alter the outputted filename. absolute paths are not currently supported, all paths are relative to outputDirectory defined by -o. if a directory is given then this directory will be prepended to the filename after the outputDirectory. | None | string |

### `overlays` keys

The `overlays` key is the main place to set your overlay operation instructions, but is an optional setting.  If working with multi-document YAML files, the items set under the `overlays` key will apply to all YAML documents in the file, unless a qualifier or combination of qualifiers `documentQuery` and `documentIndex` are provided. The `overlays` are processed prior to overlays in the [documents key](#documents-keys) instructions.  Each `overlays` list/array item can have the following keys set:


| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change you are performing, and used only for on-screen output or self-documentation. | None | string |
| query | yes | JSONPath query or JSONPath fully-qualified (dot-notation) path to value you would like to manipulate. If the `query` is not a fully-qualified path (such as a.b.c.d) and returns no matches, you need to specify the `onMissing` key (i.e. metadata.labels VS metadata.fake.*). | None | string or list/array|
| value | yes | The desired value to take action with if `query` is found. | None | str, dict, list/array |
| action | yes | The action to take when the JSONPath expression is found in the YAML document. Can be one of `delete`, `merge`, or `replace`. | None | string |
| onMissing.action | no | What to do if the JSONPath expression is not found. Can be one of `ignore` or `inject`. Only applies to the actions `merge` and `replace`| `ignore` | string |
| onMissing.injectPath | no | If your JSONPath expression was not a fully-qualified path (dot-notation) then an `injectPath` is required to qualify your action. Only applies if `onMissing: {'action': 'inject'}` is set. This should be a path or list/array of paths to inject the value if your JSONPath expression was not found in the YAML document | None | string or list/array |
| documentQuery | no | A qualifier to refine which documents the overlay is applied to.  If not set, the overlay applies to all documents in the YAML file.  Can be used in conjunction with `documentIndex`. See [Qualifiers](qualifiers.md) for more details.| None | dictionary/map |
| documentIndex | no | A qualifier to refine which documents the overlay is applies to. A list/array of YAML document indices in a multi-document YAML file.  When this is set, the overlay will only be applied to this list/array of YAML documents within the file.  Can be used in conjunction with the `documentQuery`.  See [Qualifiers](qualifiers.md) for more details. | None | list/array |


### documents keys

The `documents` list/array applies to multi-document YAML files only.  It is optional, just like the top-level `overlays` key.  If you require changes to a specific YAML document in the multi-document YAML file, this is where you define them.  Actions in the `documents` key are processed after the actions in the top-level `overlays` key.  Consider the `commonOverlays` key as a place to perform "common" changes on all YAML documents within all files listed in `yamlFiles`.  Consider the `overlays` key as the place to perform "common" changes on all YAML documents in a single file.  Actions defined in the `documents` key are for making specific changes to a specific document within the YAML file.

The keys in the `documents` list are the same as those in the [Top-Level Instructions Keys](#top-level-instructions-keys). The only difference is for you to see the `path` key as a numeric value.  This numeric value represents the positional index of the YAML document within the multi-document YAML file.  To determine the numeric value, refer to your file, and count each document (separated by `---`) starting at `0`.  

>**NOTE**: Qualifiers such as the `documentQuery` and `documentIndex` are not available here because we are performing actions on a specific document within a YAML file, and there is nothing to qualify.

In the following example of a YAML file with 3 documents, in indices which would be used for the `path` key are marked with comments to illustrate how to refer to them from the `documents` list/array `path` key.  Counting always begins at `0`.

```yaml
---
# 0
this: is_a_value
and: another_value
1: more
---
# 1
foo: bulous
app: website
---
# 2
custom: ride
bar: foo
drink: juice
```


[Back to Table of Contents](../documentation.md)  
[Next Up: Overlay Actions](actions.md)