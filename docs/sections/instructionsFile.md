[Back to Table of Contents](../documentation.md)

# Instructions File YAML Schema

What is an instructions file?  A YAML document that contains an array of paths to YAML documents to perform operations on.  Each instructions file starts with the optional key `commonOverlays`, which is an array, along with a required key `yamlFiles` which is also an array.  The instructions file is the instructions for manipulating a set of YAML files and documents.

## Top-level commonOverlays Keys

The `commonOverlays` key is completely optional, and is a means to providing overlays that should be applied to every item in the `yamlFiles` array or in other words, `yamlFiles[*].path` or `yamlFiles[*].documents.path`entries defined in the instructions. Each array item in the `commonOverlays` key is treated as a dictionary/map with the following top-level keys:

| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change you are performing, and used for either on-screen output or self-documentation. | None | string |
| query | yes | JSONPath query or JSONPath fully-qualified (dot-notation) path to value you would like to manipulate. If the `query` is not a fully-qualified path (such as a.b.c.d) and returns no matches, you need to specify the `onMissing` key (i.e. metadata.labels VS metadata.fake.*). | None | string or list/array |
| value | yes | The desired value to take action with if `query` is found. This is not required when using the `delete` action, as we are removing a value altogether. | None | int, bool, str, dict/map, list/array |
| action | yes | The action to take when the JSONPath expression is found in the YAML document. Can be one of `delete`, `format`, `merge`, or `replace`.  See [Actions](actions.md) for details. | None | string |
| onMissing.action | no | What to do if the JSONPath expression is not found. Can be one of `ignore` or `inject`. Only applies to the actions `merge` and `replace`| `ignore` | string |
| onMissing.injectPath | no | If your JSONPath expression was not a fully-qualified path (dot-notation) then an `injectPath` is required to qualify your action. Only applies if `onMissing: {'action': 'inject'}` is set. This should be a path or list/array of paths to inject the value if your JSONPath expression was not found in the YAML document | None | string or list/array |
| documentQuery | no | A qualifier to refine which documents the common overlay is applied to.  If not set, the overlay applies to all files in `yamlFiles`.  See [Qualifiers](#qualifiers) for more details. | None | dictionary/map |

## Top-Level yamlFiles Keys

Each list item in the `yamlFiles` key is treated as a dictionary/map with the following top-level keys:

| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change/file you are performing, and used only for on-screen output or self-documentation. | None | string |
| path | yes | A fully qualified path to the YAML file to modify, or a path relative to where `yot` was launched from (i.e. relative path from `pwd`). Can be a path to a YAML file or a path containing YAML files. | None | string |
| overlays | no | List/array of overlay operations to apply. If your YAML file contains multiple documents separated by `---`, then this would apply to every YAML document first, unless a qualifier or combination of qualifiers `documentQuery` and `documentIndex` are provided.  If you need to apply overlays only to a specific YAML document in a multi-document YAML file, then see the `documents` key. See [overlays keys](#overlays-keys) for available dictionary/map keys. | None | list/array of dictionaries |
| documents | no | List/array of overlay operations to apply to a multi-document YAML file.  When each document from a multi-document YAML file is loaded, an overlay can be applied by addressing the document by its index.  See [documents keys](#documents-keys) for available dictionary/map keys. | None | list/array of dictionaries/maps |

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
| documentQuery | no | A qualifier to refine which documents the overlay is applied to.  If not set, the overlay applies to all documents in the YAML file.  Can be used in conjunction with `documentIndex`. See [Qualifiers](#qualifiers) for more details.| None | dictionary/map |
| documentIndex | no | A qualifier to refine which documents the overlay is applied to, which is a list/array of YAML document indexes in a multi-document YAML file.  When this is set, the overlay will only be applied to this list/array of documents within the file.  Can be used in conjunction with the `documentQuery`.  See [Qualifiers](#qualifiers) for more details. | None | list/array |


### documents keys

The `documents` list/array applies only to multi-document YAML files and is completely optional the same as the top-level `overlays` key is.  If you require changes to a specific YAML document in the multi-document YAML file, then this is where you define them.  Actions in the `documents` key are processed after actions in the top-level `overlays` key.  Think of the `commonOverlays` key as a place to perform changes on all files listed in `yamlFiles`, while the `overlays` key is a place to perform your "common" changes on all YAML documents in a single file, and actions defined here in the `documents` key are for making specific changes to a specific document within the YAML file.

The keys in the `documents` list are the same as found in the [Top-Level Instructions Keys](#top-level-instructions-keys), except you will refer to the `path` key as a numeric value.  This numeric value represents the positional index of the YAML document within the multi-document YAML file.  You can determine this numeric value by referring to your file, and counting each document starting at `0`.  Qualifiers such as the `documentQuery` and `documentIndex` are not available here, because we are performing actions on a specific document within a YAML file, and therefore do not need to qualify anything.

Consider the following example, a YAML file with 3 documents:

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

In the example above, the indexes which would be used for the `path` key have been marked with comments to illustrate how to refer to them from the `documents` list/array `path` key.  Counting always begins at `0`.


[Back to Table of Contents](../documentation.md)  
[Next Up: Actions](qualifiers.md)