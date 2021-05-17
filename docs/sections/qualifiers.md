[Back to Table of Contents](../documentation.md)

## Overlay Qualifiers

Overlay Qualifiers are a means to further refine/qualify when an overlay is applied to a YAML document within a YAML file path.  Currently `yot` has two kinds of qualifiers, `documentQuery` and `documentIndex`.  These can be used together or separately, or not at all.


### documentQuery Overlay Qualifier

The `documentQuery` qualifier can be used on either `commonOverlays` or the `overlays` key on a `yamlFiles.path`, but cannot be used under the `documents` key.  The purpose of a `documentQuery` is to qualify an overlay operation by checking for a value or multiple values contained in a YAML document within a file.

Think of the `documentQuery` as groups of conditions that must be met before applying this overlay to the YAML document. Only one group of `conditions` must all match prior to qualifying the application of an overlay.

The `documentQuery` key is a list/array which contains a list of the following top-level keys:


#### documentQuery Top-level keys

| Key | Description | Type |
| --- | --- | --- |
| conditions | A grouping of conditions that must exist in a YAML document to qualify application of an overlay. Each `documentQuery` list/array can contain one or many `conditions` groups.  Each list/array of `conditions` contains a list/array of key/value pairs that must all return valid matches with expected values prior to qualifying application of an overlay. Each group of `conditions` is treated as an implicit "or", while the key/value conditions in each grouping is treated as an implicit "and". | list/array |


##### documentQuery conditions keys

| Key | Description | Type |
| --- | --- | --- |
| key | The key to search for within a YAML document expressed as a JSONPath query or dot-notation. | string |
| value | The value that the JSONPath query must return from one of the results of the `key`'s query before an overlay action will be applied to a document. | string |


#### documentQuery Examples

The following example demonstrates use of `commonOverlays` with a `documentQuery` to qualify when the overlay will be applied.  All key/value pairs within each `conditions` item would have to contain a valid matched result within the YAML document prior to the overlay's application.  

Think of each grouping of `conditions` as "match this" or "match this" (implicit "or").  Think of each condition within a group of `conditions` as "match this" and "match this" (implicit "and").  This allows you a great deal of flexibility on when to apply overlays. 

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

The following example demonstrates use of multiple `documentQuery` groupings.  Any single one of these key/value conditions groups would need to match within the YAML document prior to the overlay's application. Think of each group of conditions as "match this" or "match this".

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
    - key: kind
      value: Service
```

#### documentIndex Overlay Qualifier

The `documentIndex` qualifier can be used on the `overlays` key on a file path, but cannot be used under the `documents` key.  The purpose of a `documentIndex` is to qualify an overlay by specifying which specific YAML documents within a file should receive the overlay.  The `documentIndex` is a list, and should be expressed as:

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
  query: metadata.labels # required JSONPath (dot-notation)
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
      injectPath: "metadata.labels" # if your key (metadata.labels) in this instance was a JSONPath expression, we can't exactly inject to an expression.  We need a real path to plug it into. If you had a JSONPath expression and no onMissing.injectPath we would assume ignore and print a warning
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
      value: [] # the desired value of the JSONPath expression, in this case [], does not matter on a delete action
      action: delete
```

[Back to Table of Contents](../documentation.md)  
[Next Up: Order of Operations/Processing](orderOfOperations.md)