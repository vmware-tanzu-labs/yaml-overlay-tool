[Back to Actions](actions.md#3-merge)  
[Back to Table of Contents](../documentation.md)

## How data types are handled with merge actions

The type of data you are merging affects how the `merge` `action` applies `value` data is applied to a YAML document. This concept is intuitive by design, but there are a few things to be aware of to harness the full feature set of `Yot`.  You can also use the `replace` action if you do not want to merge data.


### Dictionary/Map Merge Functionality

A dictionary in YAML is a set of key value pairs. 

For example,

```yaml
# short-form
dictionaryExample1: {key1: value1, key2: value2}
# long-form
dictionaryExample2:
  key1: value1
  key2: value2 
```

During a dictionary data merge, `Yot` performs a deep merge on the original dictionary data with the new dictionary data. Any new keys added into the existing values, or any identical keys with new values are automatically updated. 

If you don't want to merge data, but would rather replace it, use the `replace` action instead.


### Array/sequence merge functionality

An array/sequence in YAML appears as a key with a list of items prefaced by a `- ` pattern or a comma-separated list contained in square brackets.

For example,

```yaml
# short-form
arrayExample1: [item1, item2, item3]
# long-form
arrayExample2:
  - item1
  - item2
  - item3
```

When merging list data, `Yot` takes the original list data and extends it with the new list data.  


### Integer/float merge functionality

In YAML, an integer appears as an unquoted set of digits without a decimal and a float (floating point numbers) appears as an unquoted set of digits with a decimal.

For example,

```yaml
# source: https://www.tutorialspoint.com/yaml/yaml_scalars_and_tags.htm

# integers
canonical: 10
decimal: +1,024
seagecimal: 1:23:45
octal: 014
hexadecimal: 0xC

# floats
canonical: 1.2345
exponential: 1.2345e+02
sexagecimal: 20:30.15
fixed: 1,234.56
negativeInfinity: -.inf
```

When merging an integer/float with an existing integer/float, `Yot` adds the two values together.


### Boolean merge functionality
A boolean in YAML is represented as `true`/`false`, `yes`/`no`, `y`/`n`, or `on`/`off`.

For example,

```yaml
booleanExample1: true
booleanExample2: false
booleanExample3: yes
booleanExample4: no
booleanExample5: y
booleanExample6: n
booleanExample7: on
booleanExample8: off
```

When merging boolean type data, `yot` performs in boolean.  

For example, 

**In this scenario, `false`/`no`/`n`/`off` always wins.**


### String merge functionality

In YAML, a string is typically represented as an unquoted set of alphanumeric characters.  To represent a number as a string, enclose the number in double-quotation marks (`"`). 

For example,

```yaml
stringExample1: dog
stringExample2: "dog"
stringExample3: dog1
stringExample4: "dog1"
stringExample5: "1234"
stringExample6: "1234.56"
```

When merging string data, `yot` takes the original string data and concatenates it with the new string data. This is not initially intuitive, but can provide some interesting use-cases.  


#### String merge examples

In a templated instructions file with multiple values files for differing Kubernetes clusters, you could have the value of `site` set differently in each values file. 

**Note:** Templating is not implemented in v0.1.0.

The following example use-cases show adding on to a Kubernetes `apiVersion`. For example, v1 + alpha2 => outputs a change of v1alpha2. 

```yaml
# dev.yaml
site: "DEV"
```

```yaml
# qa.yaml
site: "QA"
```

```yaml
# prod.yaml
site: "PROD"
```
With an overlay of:

```yaml
yamlFiles:
  - path: "examples/manifests/test.yaml"
    overlays:
      - query: metadata.name
        value: -{{ site }}
        action: merge
```

```yaml
# test.yaml
...
metadata:
  name: my-cool-app
...
```

This will render three versions of the file test.yaml for DEV, QA, and PROD, where the metadata.name field will have been extended as such `my-cool-app-DEV`, `my-cool-app-QA`, and `my-cool-app-PROD`.

[Back to Actions](actions.md#3-merge)  
[Back to Table of Contents](../documentation.md)
