[Back to Table of contents](../index.md)

## Overlay qualifiers

Overlay qualifiers refine/qualify when an overlay is applied to a YAML document within a YAML file path.  
YAML Overlay Tool has two kinds of qualifiers that can be used together, separately, or not at all:  
1. [documentQuery](#documentquery-overlay-qualifier)  
1. [documentIndex](#documentindex-overlay-qualifier)  



## Overlay qualifiers Table of Contents
<!-- @import "[TOC]" {cmd="toc" depthFrom=3 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [documentQuery overlay qualifier](#documentquery-overlay-qualifier)
  - [documentQuery top-level keys](#documentquery-top-level-keys)
    - [documentQuery conditions keys](#documentquery-conditions-keys)
  - [documentQuery examples](#documentquery-examples)
    - [Single condition](#single-condition)
    - [With multiple conditions](#with-multiple-conditions)
    - [With no Value, same conditions, but expressed in JSONPath](#with-no-value-same-conditions-but-expressed-in-jsonpath)
    - [Multiple documentQuery groups](#multiple-documentquery-groups)
- [documentIndex overlay qualifier](#documentindex-overlay-qualifier)

<!-- /code_chunk_output -->


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
| query | The key to search for within a YAML document expressed as a JSONPath query or dot-notation. (can accept multiple queries)| string |
| value | (optional) The value that the JSONPath query must return from one of the results of the `query` before an overlay action is applied to a document. | string |

>**NOTE:** If `value` is not provided then the condition will return true if any results are found from the query

#### documentQuery examples

The following example demonstrates use of `commonOverlays` with a `documentQuery` to qualify when the overlay will be applied.  This allows you a great deal of flexibility on when to apply overlays.  All query/value pairs within each `conditions` item has to contain a valid matched result within the YAML document prior to the overlay's application.  

Think of each group of `conditions` as "match this" or "match this" (implicit "or").  Think of each condition within a group of `conditions` as "match this" and "match this" (implicit "and").  

##### Single condition

```yaml
---
yamlFiles:
  - path: /file/to/modify.yaml
    overlays:
      - name: Change the namespace for all k8s Deployments
        query: metadata.namespace
        value: my-namespace
        action: replace
        documentQuery:
          - conditions:
              - query: kind
                value: Deployment
```


##### With multiple conditions

```yaml
---
# documentQuery says this must be a Deployment with a specific label to get applied
yamlFiles:
  - path: /file/to/modify.yaml
    overlays:
      - name: Change the namespace for all k8s Deployments with name label of cool-app
        query: metadata.namespace
        value: my-namespace
        action: replace
        documentQuery:
          - conditions:
              - query: kind
                value: Deployment
              - query: metadata.labels.["app.kubernetes.io/name"]
                value: cool-app
```


##### With no Value, same conditions, but expressed in JSONPath

```yaml
---
yamlFiles:
  - path: /file/to/modify.yaml
    overlays:
      - name: Change the namespace for all k8s Deployments with name label of cool-app
        query: metadata.namespace
        value: my-namespace
        action: replace
        documentQuery:
          - conditions:
              - query: $[?($.kind == "Deployment")]
              - query: metadata.labels.[?(@.name == "cool-app")]`
```


##### Multiple documentQuery groups
The following example demonstrates use of multiple `documentQuery` groups.  Any single one of these query/value conditions groups have to match within the YAML document prior to the overlay's application. Think of each group of conditions as "match this" or "match this" (implicit "or").  

```yaml
---
# must be a k8s deployment OR a k8 service to be applied
yamlFiles:
  - path: /file/to/modify.yaml
    overlays:
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


### documentIndex overlay qualifier

The `documentIndex` qualifier is used on the `overlays` key on a file path, but cannot be used under the `documents` key.  The purpose of a `documentIndex` is to qualify an overlay by specifying which specific YAML documents within a file should receive the overlay.  The `documentIndex` is a list, and should be expressed as:

```yaml
---
# only apply this to document 0, 1, or 3 in file /file/to/modify.yaml
yamlFiles:
  - path: /file/to/modify.yaml
    overlays:
      - name: Change the namespace for all k8s Deployments or Services
        query: metadata.namespace
        value: my-namespace
        action: replace
        documentIndex: [0,1,3]
```

or

```yaml
---
# only apply this to document 0, 1, or 3 in file /file/to/modify.yaml
yamlFiles:
  - path: /file/to/modify.yaml
    overlays:
      - name: Change the namespace for all k8s Deployments or Services
        query: metadata.namespace
        value: my-namespace
        action: replace
        documentIndex:
          - 0
          - 1
          - 3
```


[Back to Table of contents](../index.md)  
[Next Up: Format markers](formatMarkers.md)