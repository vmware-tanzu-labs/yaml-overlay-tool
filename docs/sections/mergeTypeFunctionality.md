## Details on How Types are Handled with Merge Actions

The `action` of `merge` can affect how the `value` data gets applied to a YAML document, depending on the type of data it is.  This is fairly intuitive by design, but there are a few things to be aware of so you can harness the full feature set of `yot`.  If you do not wish to use these features, simply use the `replace` action.


### Dictionary/Map Merge Functionality

A dictionary in YAML is a set of key value pairs, represented as such:

```yaml
# short-form
dictionaryExample1: {key1: value1, key2: value2}
# long-form
dictionaryExample2:
  key1: value1
  key2: value2 
```

When merging dictionary data, `yot` performs a deep merge on the original dictionary data with the new dictionary data.  This means any new keys are added into the existing values, and any identical keys with new values are simply updated. If this approach does not work for your situation, consider using the `replace` action.


### Array/Sequence Merge Functionality

An array/sequence in YAML appears as a key with a list of items prefaced by a `- ` pattern or a comma-separated list contained in square brackets as such:

```yaml
# short-form
arrayExample1: [item1, item2, item3]
# long-form
arrayExample2:
  - item1
  - item2
  - item3
```

When merging list data, `yot` takes the original list data, and extends it with the new list data.  This is fairly intuitive, but worth calling out for clarity.


### Integer/Float Merge Functionality

An integer in YAML appears as an unquoted set of digits, without a decimal.  A float (floating point numbers) appears as an unquoted set of digits with a decimal as such:

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

When merging an integer/float with an existing integer/float, `yot` performs a math behavior by adding the two values together.


### Boolean Merge Functionality
A boolean in YAML is represented as `true`/`false`, `yes`/`no`, `y`/`n`, or `on`/`off` as such:

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

When merging boolean type data, `yot` performs a boolean math behavior.  **In this scenario, `false`/`no`/`n`/`off` always wins.**


### String Merge Functionality

In YAML a string is typically represented as an unquoted set of alphanumeric characters.  If you wish to represent a number as a string you would enclose the number in double-quotation marks (`"`) as such:

```yaml
stringExample1: dog
stringExample2: "dog"
stringExample3: dog1
stringExample4: "dog1"
stringExample5: "1234"
stringExample6: "1234.56"
```

When merging string data, `yot` takes the original string data and concatenates it with the new string data.  This is not initially intuitive, but can provide some interesting use-cases.  


#### String Merge Examples

A few use-cases that come to mind is adding on to a Kubernetes `apiVersion` (i.e. v1 + alpha2 => outputs a change of v1alpha2).  In a templated instructions file (templating is not implemented in v0.1.0) with multiple values files for differing Kubernetes clusters, a user could have the value of `site` set differently in each values file, such as:

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

