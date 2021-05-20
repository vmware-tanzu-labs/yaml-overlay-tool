[Back to Table of Contents](../documentation.md)

# Actions

There are four types of actions that you can use to apply changes to a YAML document within Yot.

* **Delete**
* **Format**
* **Merge**
* **Replace**


There are two types of actions that you can use to apply changes to a YAML document when a `query` returns no results (`onMissing`) within Yot.

* **Ignore**
* **Inject**


## Overlay actions

### 1. Delete

The `delete` action lets a Yot user remove unwanted pieces of a YAML document.   


### 2. Format

The `format` action lets a Yot user do something with data looked up by a JSONPath query, where the returned value is represented by `%s`.  You can use `format` to put text before an existing value, after an existing value, or both before and after the existing value.  

>**NOTE:** To use `%s` at the start of the `value` field, wrap the value in double-quotation marks. For example,`"%sSomeNewText"`.


#### Format usage example

In the following example, the `name` label is now called `app.kubernetes.io/name`.

```yaml
commonOverlays:
  - name: Update all 'name' labels to 'app.kubernetes.io/name'
    query: metadata.labels.name~
    value: app.kubernetes.io/%s
    action: format

yamlFiles:
  - name: Pile of YAML files
    path: /tmp/yamls
```



### 3. Merge

The `merge` action lets a Yot user merge new data with existing data. You'll find this action works best with lists/arrays and dictionaries/maps, and that the `merge` behavior differs according to the type of data being merged.

See [Details on How Data Types are Handled with Merge Actions](mergeTypeFunctionality.md).


### 4. Replace

The `replace` action lets a Yot user replace existing data with new data.



## On Missing Actions

`onMissing` actions instruct Yot on what to do if there are no results from your JSONPath `query`.


### 1. Ignore

The `ignore` action is the default `onMissing` action if there are no results found from your `query`.  Use of the `onMissing` key is optional. Use of `ignore` can be added for the sake of being explicit to anyone reading your instructions file.  

```yaml
yamlFiles:
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

#### Inject Path

If your initial `query` used some of JSONPath's advanced features (`../`, `*`, etc) rather than a dot-notation style path (e.g: `a.b.c.d`), and no results were obtained, an `injectPath` is also required to allow for properly building the YAML paths.  An `injectPath` can either be a `string` or a `list/array` that you can use to inject the same data to multiple-locations within the file.

The following example illustrates a simple use-case for missing labels that you would like to inject if `metadata.labels` was missing in the YAML file.

```yaml
yamlFiles:
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
  - name: Find some data, and inject if it does not exist to multiple locations
    query: ..image
    value: nginx:latest
    action: replace
    onMissing:
      action: inject
      injectPath:
        - spec.template.spec.containers[0].image
```

[Back to Table of Contents](../documentation.md)  
[Next Up: Overlay Qualifiers](qualifiers.md)
