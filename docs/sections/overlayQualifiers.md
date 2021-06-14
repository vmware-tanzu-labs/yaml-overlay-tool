[Back to Table of contents](../documentation.md)

## Overlay qualifiers

Overlay qualifiers refine/qualify when an overlay is applied to a YAML document within a YAML file path.  
YAML Overlkay Tool has two kinds of qualifiers that can be used together, separately, or not at all:  
* `documentQuery`  
* `documentIndex`  


### documentQuery overlay qualifier

The `documentQuery` qualifier can be used on either `commonOverlays` or the `overlays` key on a `yamlFiles.path`, but cannot be used under the `documents` key.  The purpose of a `documentQuery` is to qualify an overlay operation by checking for a value or multiple values contained in a YAML document within a file.

The `documentQuery` includes groups of conditions that must be met before applying this overlay to the YAML document. Only one group of `conditions` must all match prior to qualifying the application of an overlay.

The `documentQuery` key is a list/array which contains a list of the following top-level keys:


#### documentQuery top-level keys

| Key | Description | Type |
| --- | --- | --- |
| conditions | A group of conditions that must exist in a YAML document to qualify application of an overlay. Each `documentQuery` list/array can contain one or more `conditions` groups.  Each list/array of `conditions` contain a list/array of query/value pairs that must all return valid matches with expected values prior to qualifying the application of an overlay. Each group of `conditions` is treated as an implicit "or", while the query/value conditions in each grouping is treated as an implicit "and". | list/array |


##### documentQuery conditions keys

| Key | Description | Type |
| --- | --- | --- |
| query | The key to search for within a YAML document expressed as a JSONPath query or dot-notation. | string |
| value | The value that the JSONPath query must return from one of the results of the `query` before an overlay action is applied to a document. | string |


#### documentQuery examples

The following example demonstrates use of `commonOverlays` with a `documentQuery` to qualify when the overlay will be applied.  This allows you a great deal of flexibility on when to apply overlays.  All query/value pairs within each `conditions` item has to contain a valid matched result within the YAML document prior to the overlay's application.  

Think of each group of `conditions` as "match this" or "match this" (implicit "or").  Think of each condition within a group of `conditions` as "match this" and "match this" (implicit "and").  

```yaml
commonOverlays:
- name: Change the namespace for all k8s Deployments
  query: metadata.namespace
  value: my-namespace
  action: replace
  documentQuery:
  - conditions:
    - query: kind
      value: Deployment

# With multiple conditions, must be a Deployment with a specific label to get applied
commonOverlays:
- name: Change the namespace for all k8s Deployments with name label of cool-app
  query: metadata.namespace
  value: my-namespace
  action: replace
  documentQuery:
  - conditions:
    - query: kind
      value: Deployment
    - query: metadata.labels.`app.kubernetes.io/name`
      value: cool-app
```

The following example demonstrates use of multiple `documentQuery` groups.  Any single one of these query/value conditions groups have to match within the YAML document prior to the overlay's application. Think of each group of conditions as "match this" or "match this" (implicit "or").  

```yaml
commonOverlays:
- name: Change the namespace for all k8s Deployments or Services
  query: metadata.namespace
  value: my-namespace
  action: replace
  documentQuery:
  - conditions:
    - query: kind
      value: Deployment
  - conditions:
    - query: kind
      value: Service
```

#### documentIndex overlay qualifier

The `documentIndex` qualifier is used on the `overlays` key on a file path, but cannot be used under the `documents` key.  The purpose of a `documentIndex` is to qualify an overlay by specifying which specific YAML documents within a file should receive the overlay.  The `documentIndex` is a list, and should be expressed as:

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


[Back to Table of contents](../documentation.md)  
[Next Up: Format markers](formatMarkers.md)