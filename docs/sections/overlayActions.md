[Back to Table of contents](../index.md)

# Actions

There are four types of actions that you can use to apply changes to a YAML document within Yot.

1. **[Combine](#1-combine)**
1. **[Delete](#2-delete)**
1. **[Merge](#3-merge)**
1. **[Replace](#4-replace)**


There are two types of actions that you can use to apply changes to a YAML document when a `query` returns no results (`onMissing`) within Yot.

1. **[Ignore](#1-ignore)**
1. **[Inject](#2-inject)**


## Overlay actions

If the `action` key is omitted from an overlay, an action of `merge` is assumed.


### 1. Combine

The `combine` action lets a Yot user combine booleans with booleans, integers with integers, and strings with strings.  The product of this action is a new value.

1. Combining booleans produces the result of boolean addition.
1. Combining integers produces the result of integer addition.
1. Combining strings produces the result of string concatenation.


### 2. Delete

The `delete` action lets a Yot user remove unwanted pieces of a YAML document.   


### 3. Merge

As mentioned above, `merge` is the default action, should the `action` keyword be omitted.  
The `merge` action lets a Yot user merge new data with existing data, and is primarily used for performing a deep merge on a map/dictionary and arrays.  

However, `merge` can also be used to format string/scalar data with some special [Format Markers](formatMarkers.md).  When merging scalar data, it is treated as a `replace`.  However, you can use the format markers to manipulate the existing data with `%v`, or insert a new value or a new line comment.


### 4. Replace

The `replace` action lets a Yot user replace existing data with new data.


## OnMissing actions

`onMissing` actions instruct Yot on what to do if there are no results from your JSONPath `query`.


### 1. Ignore

The `ignore` action is the default `onMissing` action if there are no results found from your `query`.  Use of the `onMissing` key is optional. Use of `ignore` can be added for the sake of being explicit to anyone reading your instructions file.  

```yaml
yamlFiles:
  - path: /some/yaml/file.yaml
    overlays:
      - name: Replace labels if they exist
        query: metadata.labels
        value:
          label1: newLabel
          label2: newLabel
        action: replace
        # the following 2 lines are not required, and this would be considered long-form
        onMissing:
          action: ignore
```


### 2. Inject

Use `inject` if your `query` returned no results, but you still want to insert data.

#### injectPath

If your initial `query` used some of JSONPath's advanced features (`../`, `*`, etc) rather than a dot-notation style path (e.g: `a.b.c.d`), and no results were obtained, an `injectPath` is also required to allow for properly building the YAML paths.  An `injectPath` can either be a `string` or a `list/array` that you can use to inject the same data to multiple-locations within the file.

The following example illustrates a simple use-case for missing labels that you would like to inject if `metadata.labels` was missing in the YAML file.

```yaml
yamlFiles:
  - path: /some/yaml/file.yaml
    overlays:
      - name: Replace labels if they exist, otherwise inject them
        query: metadata.labels
        value:
          label1: newLabel
          label2: newLabel
        action: replace
        onMissing:
          action: inject
```


```yaml
yamlFiles:
  - path: /some/yaml/file.yaml
    overlays:
      - name: Find some data, and inject if it does not exist to multiple locations
        query: ..image
        value: nginx:latest
        action: replace
        onMissing:
          action: inject
          injectPath:
            - spec.template.spec.containers[0].image
```

[Back to Table of contents](../index.md)  
[Next Up: Overlay qualifiers](overlayQualifiers.md)
