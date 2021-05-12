### The Instructions File

What is an instructions file?  A yaml document that contains a list of yaml documents to perform operations on.  Each instructions file starts with the key `yamlFiles`, which is a list, along with an optional key `commonOverlays` which is also a list.

#### Top-level commonOverlays Keys

The `commonOverlays` key is completely optional, and is a means to providing overlays that should be applied to every `yamlFiles` list item or `yamlFiles[*].documents[*].path` defined in the instructions. Each list item in the `commonOverlays` key is treated as a dictionary with the following top-level keys:

| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change you are performing, and used only for on-screen output or self-documentation. | None | string |
| query | yes | JSONpath query or JSONpath fully-qualified (dot-notation) path to value you would like to manipulate. If the `query` is not a fully-qualified path (such as a.b.c.d) and returns no matches, you need to specify the `onMissing` key (i.e. metadata.labels VS metadata.fake.*). | None | string or list|
| value | yes | The desired value to take action with if `query` is found. | None | str, dict, list |
| action | yes | The action to take when the JSONPath expression is found in the yaml document. Can be one of `delete`, `merge`, or `replace`. | None | string |
| onMissing.action | no | What to do if the JSONpath expression is not found. Can be one of `ignore` or `inject`. Only applies to the actions `merge` and `replace`| `ignore` | string |
| onMissing.injectPath | no | If your JSONpath expression was not a fully-qualified path (dot-notation) then an `injectPath` is required to qualify your action. Only applies if `onMissing: {'action': 'inject'}` is set. This should be a path or list of paths to inject the value if your JSONpath expression was not found in the yaml document | None | string or list |
| documentQuery | no | A qualifier to refine which documents the common overlay is applied to.  If not set, the overlay applies to all files in `yamlFiles`.  See [Qualifiers](#qualifiers) for more details. | None | dictionary |

#### Top-Level yamlFiles Keys

Each list item in the `yamlFiles` key is treated as a dictionary with the following top-level keys:

| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change/file you are performing, and used only for on-screen output or self-documentation. | None | string |
| path | yes | A fully qualified path to the yaml file to modify, or a path relative to where `yot` was launched from (i.e. relative path from `pwd`). Can be a path to a yaml file or a path containing yaml files. | None | string |
| overlays | no | List of overlay operations to apply. If your yaml file contains multiple documents separated by `---`, then this would apply to every yaml document first, unless a qualifier or combination of qualifiers `documentQuery` and `documentIndex` are provided.  If you need to apply overlays to a specific yaml document in a multi-document yaml file, then see the `documents` key. See [overlays keys](#overlays-keys) for available dictionary keys. | None | list of dictionaries |
| documents | no | List of overlay operations to apply to a multi-document yaml file.  When each document from a multi-document yaml file is loaded, an overlay can be applied by addressing the document by its index.  See [documents keys](#documents-keys) for available dictionary keys. | None | list of dictionaries |

#### overlays keys

The `overlays` key is the main place to set your overlay operation instructions, but is an optional setting.  If working with multi-document yaml files, the items set under the `overlays` key will apply to all yaml documents in the file, unless a qualifier or combination of qualifiers `documentQuery` and `documentIndex` are provided. The `overlays` are processed prior to overlays in the [documents key](#documents-keys) instructions.  Each `overlays` list item can have the following keys set:


| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change you are performing, and used only for on-screen output or self-documentation. | None | string |
| query | yes | JSONpath query or JSONpath fully-qualified (dot-notation) path to value you would like to manipulate. If the `query` is not a fully-qualified path (such as a.b.c.d) and returns no matches, you need to specify the `onMissing` key (i.e. metadata.labels VS metadata.fake.*). | None | string or list|
| value | yes | The desired value to take action with if `query` is found. | None | str, dict, list |
| action | yes | The action to take when the JSONPath expression is found in the yaml document. Can be one of `delete`, `merge`, or `replace`. | None | string |
| onMissing.action | no | What to do if the JSONpath expression is not found. Can be one of `ignore` or `inject`. Only applies to the actions `merge` and `replace`| `ignore` | string |
| onMissing.injectPath | no | If your JSONpath expression was not a fully-qualified path (dot-notation) then an `injectPath` is required to qualify your action. Only applies if `onMissing: {'action': 'inject'}` is set. This should be a path or list of paths to inject the value if your JSONpath expression was not found in the yaml document | None | string or list |
| documentQuery | no | A qualifier to refine which documents the overlay is applied to.  If not set, the overlay applies to all documents in the yaml file.  Can be used in conjunction with `documentIndex`. See [Qualifiers](#qualifiers) for more details.| None | dictionary |
| documentIndex | no | A qualifier to refine which documents the overlay is applied to, which is a list of yaml document indexes in a multi-document yaml file.  When this is set, the overlay will only be applied to this list of documents within the file.  Can be used in conjunction with the `documentQuery`.  See [Qualifiers](#qualifiers) for more details. | None | list |


#### documents keys

The `documents` list applies only to multi-document yaml files and is completely optional the same as the top-level `overlays` key is.  If you require changes to a specific yaml document in the multi-document yaml file, then this is where you define them.  Actions in the `documents` key are processed after actions in the top-level `overlays` key.  Think of the `commonOverlays` key as a place to perform changes on all files listed in `yamlFiles`, while the `overlays` key is a place to perform your "common" changes within a single file, and actions defined here in the `documents` key are for making specific changes to a specific document within the yaml file.

The keys in the `documents` list are the same as found in the [Top-Level Instructions Keys](#top-level-instructions-keys), except you will refer to the `path` key as a numeric value.  This numeric value represents the positional index of the yaml document within the multi-document yaml file.  You can determine this numeric value by referring to your file, and counting each document starting at `0`.  Qualifiers such as the `documentQuery` and `documentIndex` are not available here, because we are performing actions on a specific document within a file.

Consider the following example, a yaml file with 3 documents:

```yaml
---
# 0
this: is_a_value
and: another_value
1: more
---
# 1
foo: bulous
app: c
---
# 2
custom: value
bar: foo
drink: juice
```

In the example above, the index which would be used for the `path` key have been marked with comments to illustrate how to refer to them from the `documents` list `path` key.  Counting always begins at `0`.


#### Qualifiers

Qualifiers are a means to further refine when an overlay is applied to a yaml document within a file path.  Currently `yot` has two kinds of qualifiers, `documentQuery` and `documentIndex`.  These can be used together or separately, or not at all.


##### documentQuery Qualifier

The `documentQuery` qualifier can be used on either `commonOverlays` or the `overlays` key on a `yamlFiles.path`, but cannot be used under the `documents` key.  The purpose of a `documentQuery` is to qualify an overlay operation by checking for a value or multiple values contained in a yaml document within a file.

Think of the `documentQuery` as groups of conditions that must be met before applying this overlay to the yaml document. Only one group of `conditions` must all match prior to qualifying the application of an overlay.

The `documentQuery` key is a list/array which contains a list of the following top-level keys:


###### documentQuery Top-level keys

| Key | Description | Type |
| --- | --- | --- |
| conditions | A grouping of conditions that must exist in a yaml document to qualify application of an overlay. Each `documentQuery` list can contain one or many condition groups.  Each list of `conditions` contains a list of key/value pairs that must all return valid matches with expected values prior to qualifying application of an overlay. | list |


###### documentQuery conditions keys

| Key | Description | Type |
| --- | --- | --- |
| key | The key to search for within a yaml document expressed as a JSONPath query or dot-notation. | string |
| value | The value that the JSONPath query must return from one of the results of the `key`'s query before an overlay action will be applied to a document. | string |


##### documentQuery Examples

The following example demonstrates use of `commonOverlays` with a `documentQuery` to qualify when the overlay will be applied.  All key/value pairs within each `conditions` item would have to contain a valid matched result within the yaml document prior to the overlay's application.  

Think of each grouping of `conditions` as "match this" or "match this".  Think of each condition within a group of `conditions` as "match this" and "match this".

```yaml
commonOverlays:
- name: Change the namespace for all k8s Deployments
  query: metadata.namespace
  value: my-namespace
  action: replace
  documentQuery:
  - conditions:
    - key: kind
      value: Deployment

# With multiple conditions, must be a Deployment with a specific label to get applied
commonOverlays:
- name: Change the namespace for all k8s Deployments with name label of cool-app
  query: metadata.namespace
  value: my-namespace
  action: replace
  documentQuery:
  - conditions:
    - key: kind
      value: Deployment
    - key: metadata.labels.`app.kubernetes.io/name`
      value: cool-app
```

The following example demonstrates use of multiple `documentQuery` groupings.  Any single one of these key/value conditions groups would need to match within the yaml document prior to the overlay's application. Think of each group of conditions as "match this" or "match this".

```yaml
commonOverlays:
- name: Change the namespace for all k8s Deployments or Services
  query: metadata.namespace
  value: my-namespace
  action: replace
  documentQuery:
  - conditions:
    - key: kind
      value: Deployment
  - conditions:
    - key: kind:
      value: Service
```

##### documentIndex Qualifier

The `documentIndex` qualifier can be used on the `overlays` key on a file path, but cannot be used under the `documents` key.  The purpose of a `documentIndex` is to qualify an overlay by specifying which specific yaml documents within a file should receive the overlay.  The `documentIndex` is a list, and should be expressed as:

```yaml
documentIndex: [0,1,3]
```

or

```yaml
documentIndex:
  - 0
  - 1
  - 3
```


#### Instructions File Full-Specification Example

```yaml
---
commonOverlays: # optional way to apply overlays to all 'yamlFiles'
- name: Apply common label only to k8s services # optional key
  query: metadata.labels # required JSONpath (dot-notation)
  value: # desired value to perform an action on matches of the query with
    some: label
  action: merge # merge | replace | delete
  onMissing: # optional - what to do if 'query' not found in yaml
    action: inject # inject | ignore, default of ignore if onMissing not set
  documentQuery: # qualifier
  # array/list of condition groupings. Each array of conditions is treated separately
  # each grouping of conditions must all match. If 1 group of conditions returns
  ## True, then the overlay will get applied
  - conditions:
    - key: kind # search for the 'kind' key in the yaml doc
      value: Service # we expect the result of the 'kind' key to be this value before applying the overlay
yamlFiles: # what to overlay onto
- name: "some arbitrary descriptor" # Name is Optional
  path: "path/relative/to/directory/of/execution.yaml" # or
  # path: "/fully/qualified/path.yaml"
  overlays: # if multi-doc yaml file, applies to all docs, gets applied first
  - name: Inject label to documents 0 2 or 4 if a Deployment
    query: metadata.labels.foo
    value: {{ foo }} # example with jinja2 templating
    action: "replace" # merge, replace, delete
    onMissing:
      action: "inject" # inject | ignore
      injectPath: "metadata.labels" # if your key (metadata.labels) in this instance was a jsonpath expression, we can't exactly inject to an expression.  We need a real path to plug it into. If you had a jsonpath expression and no onMissing.injectPath we would assume ignore and print a warning
    documentQuery: # qualifier, only modify if a k8s Deployment
      key: kind
      value: Deployment
    documentIndex: # qualifier, only modify docs 0, 2, and 4 in multi-yaml doc
    - 0
    - 2
    - 4
  documents: # optional and only used for multi-doc yaml files
  # need to refer to their path by their index
  - name: the manifest that does something
    path: 0
    overlays:
    - query: a.b.c.d
      value: [] # the desired value of the JSONpath expression, in this case [], does not matter on a delete action
      action: delete
```