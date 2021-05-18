[Back to Table of Contents](../documentation.md)

# Overlay Actions

Within `yot` there are four types of actions you can take to apply a change to a YAML document.  These actions help you achieve a desired outcome when it comes to manipulating existing YAML files.

There are also actions you can take when a JSONPath `query` returns no results, by making use of `onMissing`.

## Overlay Actions

### 1. Delete

The `delete` action provides a user of `yot` a mechanism to remove unwanted pieces of a YAML document.  


### 2. Format

The `format` action provides a user of `yot` a mechanism to take existing data looked up by a JSONPath query, and represented by `%s`, to do something with it.  Typically `format` is used for either putting text before the existing value, after it, or a combination thereof.  

>**NOTE:** *If you would like to use `%s` at the start of the `value` field, you will need to wrap the value in double-quotation marks like `"%sSomeNewText"`.*


#### Format Usage Example

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

In the example above, the `name` label will now be called `app.kubernetes.io/name`.


### 3. Merge

The `merge` action provides a user of `yot` a mechanism to merge new data with existing data.  Depending on the type of data being merged, merge behaves differently.  `merge` is best used with lists/arrays and dictionaries/maps.

Please see [Details on How Data Types are Handled with Merge Actions](mergeTypeFunctionality.md)


### 4. Replace

The `replace` action provides a user of `yot` a mechanism to completely replace existing data with new data.



## On Missing Actions

`onMissing` actions are used to instruct `yot` on what to do if there are no results from your JSONPath `query`.


### 1. Ignore

`ignore` is always the default action when no results are obtained from your `query`.  It is not required to add the `onMissing` key if you do not care that something is missing, unless you want to be explicit to anyone reading your instructions file.  `yot` will not act on data that does not exist.

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

`inject` is required if your `query` returned no results, and you wish to insert the data anyways.

#### Inject Path
If your initial `query` used some of the more powerful JSONPath operations, rather than a dot-notation style path (e.g: `a.b.c.d`), and obtained no results, an `injectPath` is also required.  An `injectPath` can be either a `string` or a `list/array` for cases where you may want to inject the same data to multiple-locations within the file.

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

The example above illustrates a simple use-case for missing labels that you would like to inject.

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