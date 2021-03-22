# YAML Overlay Tool (yot)

A yaml overlay tool with templating tendencies.

## Table of Contents
<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Table of Contents](#table-of-contents)
- [Why?](#why)
- [Usage](#usage)
	- [Example Usage](#example-usage)
- [Setup](#setup)
- [Quick Start](#quick-start)
	- [The Instructions File](#the-instructions-file)
		- [Top-level common_overlays Keys](#top-level-commonoverlays-keys)
		- [Top-Level yaml_files Keys](#top-level-yamlfiles-keys)
		- [overlays keys](#overlays-keys)
		- [documents keys](#documents-keys)
		- [Qualifiers](#qualifiers)
			- [document_query Qualifier](#documentquery-qualifier)
			- [document_query Example](#documentquery-example)
			- [document_index Qualifier](#documentindex-qualifier)
		- [Instructions File Full-Specification Example](#instructions-file-full-specification-example)
- [Order of Operations](#order-of-operations)
	- [1a. Dynamic Instructions (Templating)](#1a-dynamic-instructions-templating)
		- [Templating Example](#templating-example)
	- [1b. Static Instructions (No Templating)](#1b-static-instructions-no-templating)
	- [2. Processing Instructions](#2-processing-instructions)
	- [3. Output](#3-output)
- [Details on How Types are Handled with Merge Actions](#details-on-how-types-are-handled-with-merge-actions)
	- [Dictionary Merge Action](#dictionary-merge-action)
	- [List Merge Action](#list-merge-action)
	- [String Merge Action](#string-merge-action)
		- [String Merge Examples](#string-merge-examples)
- [Author](#author)
- [License](#license)
- [Contributing](#contributing)
- [Code of Conduct](#code-of-conduct)
- [Communication](#communication)
	- [E-Mail](#e-mail)

<!-- /TOC -->


## Why?

`yot` was designed to be flexible, simple, and familiar.  Whether you want the ability to use a templating language to transform yaml data, or just change a couple values in a yaml document, `yot` makes it possible.  

Our philosophy is to treat 3rd party yaml manifests as source code.  We don't want to manage templated yaml files, we want to manage patches (overlays) for yaml files, and furthermore we want to keep any potential templating outside of the source yaml files.  

Templated manifests become hard to manage overtime and can be difficult to read.  `yot` allows us to take yaml documents from multiple sources and transform them through overlays to fit our environment's requirements.  This allows us to take any generic yaml file and manipulate it to suit our purpose, without contaminating the original file.

This practice gets you out of the cycle of updating and managing complex yaml templates.  At the same time `yot`'s instructions file specification serves as documentation as code, where you have now essentially documented all the required changes to source yaml files in one place.

The use of JSONpath queries and Jinja2 templating give the tool familiar interfaces, making adoption easier, and a more pleasant end-user experience.  The specification, or instructions file, is put together in a declarative way, where we only operate on what has been defined.  We take actions based on JSONpath query results.  We provide flexibility by allowing your instructions to be templated if needed.  The following sections will help get you moving along with `yot`!

## Usage
`./yot -h`

### Example Usage
```bash
# with no templating
./yot \
    -i ./examples/static/instructions.yaml \
    -o ./out

# with templating of overlays from a directory path
./yot \
    -i ./examples/templated/instructions.yaml \
    -v ./examples/values \
    -o ./out

# more to come soon
```

## Setup
Install the required Python3 libraries:
```bash
pip3 install -r ./requirements.txt
```

If you would like to try an interactive tutorial on installing `yot`, go [here](https://katacoda.com/ahuffman/scenarios/tool-setup)

## Quick Start
`yot` is not a templating tool in the sense of a traditional text-based templating tool.  `yot` is primarily an overlayment tool, meaning we take fragments of yaml configuration and apply or inject them over the top of an existing yaml configuration.  

`yot` also includes a templating feature by making use of the popular [Jinja2](https://jinja.palletsprojects.com/en/master/templates/) templating language to allow for templated overlay values which are rendered into memory and processed at run-time.  Use of the templating engine is completely optional, and instruction files can be static yaml documents.  When using the templating engine, anything within the instructions file can be templated, but keep in mind the document must template into a valid yaml document.  The template feature is useful if you are managing multi-environment yaml configurations.  Values files are also treated as templates, and can contain jinja2 content.  This is most useful in a scenario with a lot of values, where you would like to organize them into separate files and use the `{% include 'addl_values.yaml' %}` [Jinja2](https://jinja.palletsprojects.com/en/master/templates/#include) tag.  This is also extremely useful for large instruction files.

Each overlay operation can be performed with a JSONpath query.  If a JSONpath query produces no results in the yaml document, a desired value can be either ignored (default behavior) or injected (`on_missing: {'action': 'inject'}`) and provided a specific path (`on_missing: {'inject_path': []}`) to inject the value if the query was not a fully-qualified JSONpath (example: `metadata.labels.*` <= query VS. `metadata.labels` <= fully-qualified path).

### The Instructions File
What is an instructions file?  A yaml document that contains a list of yaml documents to perform operations on.  Each instructions file starts with the key `yaml_files`, which is a list, along with an optional key `common_overlays` which is also a list.

#### Top-level common_overlays Keys
The `common_overlays` key is completely optional, and is a means to providing overlays that should be applied to every `yaml_files` list item or `yaml_files[*].documents[*].path` defined in the instructions. Each list item in the `common_overlays` key is treated as a dictionary with the following top-level keys:

| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change you are performing, and used only for on-screen output or self-documentation. | None | string |
| query | yes | JSONpath query or JSONpath fully-qualified (dot-notation) path to value you would like to manipulate. If the `query` is not a fully-qualified path (such as a.b.c.d) and returns no matches, you need to specify the `on_missing` key (i.e. metadata.labels VS metadata.fake.*). | None | string or list|
| value | yes | The desired value to take action with if `query` is found. | None | str, dict, list |
| action | yes | The action to take when the JSONPath expression is found in the yaml document. Can be one of `delete`, `merge`, or `replace`. | None | string |
| on_missing.action | no | What to do if the JSONpath expression is not found. Can be one of `ignore` or `inject`. Only applies to the actions `merge` and `replace`| `ignore` | string |
| on_missing.inject_path | no | If your JSONpath expression was not a fully-qualified path (dot-notation) then an `inject_path` is required to qualify your action. Only applies if `on_missing: {'action': 'inject'}` is set. This should be a path or list of paths to inject the value if your JSONpath expression was not found in the yaml document | None | string or list |
| document_query | no | A qualifier to refine which documents the common overlay is applied to.  If not set, the overlay applies to all files in `yaml_files`.  See [Qualifiers](#qualifiers) for more details. | None | dictionary |

#### Top-Level yaml_files Keys
Each list item in the `yaml_files` key is treated as a dictionary with the following top-level keys:

| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change/file you are performing, and used only for on-screen output or self-documentation. | None | string |
| path | yes | A fully qualified path to the yaml file to modify, or a path relative to where `yot` was launched from (i.e. relative path from `pwd`). Can be a path to a yaml file or a path containing yaml files. | None | string |
| overlays | no | List of overlay operations to apply. If your yaml file contains multiple documents separated by `---`, then this would apply to every yaml document first, unless a qualifier or combination of qualifiers `document_query` and `document_index` are provided.  If you need to apply overlays to a specific yaml document in a multi-document yaml file, then see the `documents` key. See [overlays keys](#overlays-keys) for available dictionary keys. | None | list of dictionaries |
| documents | no | List of overlay operations to apply to a multi-document yaml file.  When each document from a multi-document yaml file is loaded, an overlay can be applied by addressing the document by its index.  See [documents keys](#documents-keys) for available dictionary keys. | None | list of dictionaries |

#### overlays keys
The `overlays` key is the main place to set your overlay operation instructions, but is an optional setting.  If working with multi-document yaml files, the items set under the `overlays` key will apply to all yaml documents in the file, unless a qualifier or combination of qualifiers `document_query` and `document_index` are provided. The `overlays` are processed prior to overlays in the [documents key](#documents-keys) instructions.  Each `overlays` list item can have the following keys set:


| key | required | description | default | type |
| --- | --- | --- | --- | --- |
| name | no | An optional description of the change you are performing, and used only for on-screen output or self-documentation. | None | string |
| query | yes | JSONpath query or JSONpath fully-qualified (dot-notation) path to value you would like to manipulate. If the `query` is not a fully-qualified path (such as a.b.c.d) and returns no matches, you need to specify the `on_missing` key (i.e. metadata.labels VS metadata.fake.*). | None | string or list|
| value | yes | The desired value to take action with if `query` is found. | None | str, dict, list |
| action | yes | The action to take when the JSONPath expression is found in the yaml document. Can be one of `delete`, `merge`, or `replace`. | None | string |
| on_missing.action | no | What to do if the JSONpath expression is not found. Can be one of `ignore` or `inject`. Only applies to the actions `merge` and `replace`| `ignore` | string |
| on_missing.inject_path | no | If your JSONpath expression was not a fully-qualified path (dot-notation) then an `inject_path` is required to qualify your action. Only applies if `on_missing: {'action': 'inject'}` is set. This should be a path or list of paths to inject the value if your JSONpath expression was not found in the yaml document | None | string or list |
| document_query | no | A qualifier to refine which documents the overlay is applied to.  If not set, the overlay applies to all documents in the yaml file.  Can be used in conjunction with `document_index`. See [Qualifiers](#qualifiers) for more details.| None | dictionary |
| document_index | no | A qualifier to refine which documents the overlay is applied to, which is a list of yaml document indexes in a multi-document yaml file.  When this is set, the overlay will only be applied to this list of documents within the file.  Can be used in conjunction with the `document_query`.  See [Qualifiers](#qualifiers) for more details. | None | list |


#### documents keys
The `documents` list applies only to multi-document yaml files and is completely optional the same as the top-level `overlays` key is.  If you require changes to a specific yaml document in the multi-document yaml file, then this is where you define them.  Actions in the `documents` key are processed after actions in the top-level `overlays` key.  Think of the `common_overlays` key as a place to perform changes on all files listed in `yaml_files`, while the `overlays` key is a place to perform your "common" changes within a single file, and actions defined here in the `documents` key are for making specific changes to a specific document within the yaml file.

The keys in the `documents` list are the same as found in the [Top-Level Instructions Keys](#top-level-instructions-keys), except you will refer to the `path` key as a numeric value.  This numeric value represents the positional index of the yaml document within the multi-document yaml file.  You can determine this numeric value by referring to your file, and counting each document starting at `0`.  Qualifiers such as the `document_query` and `document_index` are not available here, because we are performing actions on a specific document within a file.

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
Qualifiers are a means to further refine when an overlay is applied to a yaml document within a file path.  Currently `yot` has two kinds of qualifiers, `document_query` and `document_index`.  These can be used together or separately, or not at all.

##### document_query Qualifier

The `document_query` qualifier can be used on either `common_overlays` or the `overlays` key on a file path, but cannot be used under the `documents` key.  The purpose of a `document_query` is to qualify an overlay with another value contained in a yaml document within a file.

| Key | Description | Type |
| --- | --- | --- |
| key | The key to search for within a yaml document written as a fully-qualified JSONpath expression (dot-notation) | string |
| value | The value that the query must return from the `key` query, before an overlay action will be applied to a document. | string |

##### document_query Example

```yaml
common_overlays:
- name: Change the namespace for all k8s Deployments
  query: metadata.namespace
  value: my-namespace
  action: replace
  document_query:
    key: kind
    value: Deployment
```

##### document_index Qualifier

The `document_index` qualifier can be used on the `overlays` key on a file path, but cannot be used under the `documents` key.  The purpose of a `document_index` is to qualify an overlay by specifying which specific yaml documents within a file should receive the overlay.  The `document_index` is a list, and should be expressed as:

```yaml
document_index: [0,1,3]
```

or

```yaml
document_index:
  - 0
  - 1
  - 3
```


#### Instructions File Full-Specification Example

```yaml
---
common_overlays: # optional way to apply overlays to all 'yaml_files'
- name: Apply common label only to k8s services # optional key
  query: metadata.labels # required JSONpath (dot-notation)
  value: # desired value to perform an action on matches of the query with
    some: label
  action: merge # merge | replace | delete
  on_missing: # optional - what to do if 'query' not found in yaml
    action: inject # inject | ignore, default of ignore if on_missing not set
  document_query: # qualifier
    key: kind # search for the 'kind' key in the yaml doc
    value: Service # we expect the result of the 'kind' key to be this value before applying the overlay
yaml_files: # what to overlay onto
- name: "some arbitrary descriptor" # Name is Optional
  path: "path/relative/to/directory/of/execution.yaml" # or
  # path: "/fully/qualified/path.yaml"
  overlays: # if multi-doc yaml file, applies to all docs, gets applied first
  - name: Inject label to documents 0 2 or 4 if a Deployment
    query: metadata.labels.foo
    value: {{ foo }} # example with jinja2 templating
    action: "replace" # merge, replace, delete
    on_missing:
      action: "inject" # inject | ignore
      inject_path: "metadata.labels" # if your key (metadata.labels) in this instance was a jsonpath expression, we can't exactly inject to an expression.  We need a real path to plug it into. If you had a jsonpath expression and no on_missing.inject_path we would assume ignore and print a warning
    document_query: # qualifier, only modify if a k8s Deployment
      key: kind
      value: Deployment
    document_index: # qualifier, only modify docs 0, 2, and 4 in multi-yaml doc
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

## Order of Operations

### 1a. Dynamic Instructions (Templating)
When one or more values files have been passed to `yot` a `defaults.yaml` or `defaults.yml` file must exist within that path, when passing a directory with the `-v` option.  

Alternatively, a single or multiple default value files can be passed with the `-d` option.
If you only have one site and a desire to template, then only provide a default values file with `-d` or a file named `defaults.yaml` or `defaults.yml` within a directory passed with the `-v` option.

If you pass multiple default values files with the `-d` option, the first file passed serves as the base values.  Each additional default values files passed with the `-d` option are merged over the top of the base values sequentially in the order they were passed to `yot`.  Additionally, if a `defaults.yaml` or `defaults.yml` file was also present in a path passed with the `-v` option, those values will be merged over the base values last.

Any additional values files (files passed with `-v`) are treated as additional sites (or site values), and the values contained in those files are merged over the top of the `defaults.yaml` file's values.  The base filename of the additional values files are used to output the modified yaml files into.  

#### Templating Example
Consider the following example:

```bash
$ tree
├── instructions.yaml
└── values
    ├── bar
    ├── defaults.yaml
    ├── foo.yaml
    └── test.yml

./yot -v values -i instructions.yaml -o ./output
```

`yot` only takes values files with the file extension of `.yml` or `.yaml`, therefore the file `bar` above would not be read.  The values contained in values/defaults.yaml would be read and then the instructions file would be rendered by applying the values of foo.yaml over defaults.yaml, and then render the instructions.yaml template into memory for processing.  Then the test.yml will merge over defaults.yaml and render the instructions.yaml again, producing a total of 2 unique instruction sets.  When processing is complete, the modified yaml documents will be placed in `./output/foo` and `./output/test` respectively.  
> **NOTE:**  Please be careful about naming your values files, as a values file called test.yaml and test.yml would dump both of their rendered contents into ./output/test

### 1b. Static Instructions (No Templating)
If no values files were passed, then the instructions file skips the templating mechanism and is loaded directly, as long as the document is valid yaml.

Consider the following example:
```bash
yot -i instructions.yaml -o ./output
```

`yot` will read the instructions.yaml file and any outputs will be placed directly in ./output/ .

### 2. Processing Instructions
After the instructions have either been rendered or read, they are processed.

Processing begins at the instruction's `common_overlays` if they exist in the instruction set.  `yot` combines the `common_overlays` with the `yaml_files.path.overlays` first if they exist.  If `yaml_files.path.overlays` do not exist, `yot` combines `common_overlays` with `yaml_files.path.documents.path.overlays`.  In both of these scenarios, the `common_overlays` will always be applied prior to the more granular overlays.

If no `common_overlays` have been defined, processing starts at the first item of the `yaml_files` list and processes one yaml file `path` at a time sequentially.  Within each yaml file `path`, `overlays` are processed first if set, followed by the items within the `documents` key, which applies overlays to specific documents within a multi-yaml document yaml file. A single yaml document in a yaml file could still be referred to by `path: 0` from the `documents` key if desired.

### 3. Output
If no templating was performed, then output is sent to the output directory specified at runtime, or if it was not provided, the default path of `./output`.

If values files in addition to a defaults.yaml were passed at runtime, then the output for each values file will be `<output directory>/<value file basename>`.

If only a defaults.yaml value file was passed, then the output will be placed in the output directory specified at runtime, or if it was not provided, the default path of `./output`.

## Details on How Types are Handled with Merge Actions
The action of `merge` can affect how the `value` data gets applied to a yaml document, depending on the type of data it is.  This is fairly intuitive by design, but there are a few things to know so you can harness the full feature set of `yot`.


### Dictionary Merge Action
When merging dictionary data, `yot` performs a deep merge on the original dictionary data with the new dictionary data.  This means any new keys are added into the existing values, and any identical keys with new values are simply updated. If this approach does not work for your situation, consider using the `replace` action.

### List Merge Action
When merging list data, `yot` takes the original list data, and extends it with the new list data.  This is fairly intuitive, but worth calling out for clarity.

### String Merge Action
When merging string data, `yot` takes the original string data and concatenates it with the new string data.  This is not initially intuitive, but can provide some interesting use-cases.  

#### String Merge Examples
A few use-cases that come to mind is adding on to a kubernetes `apiVersion` (i.e. v1 + alpha2 => outputs a change of v1alpha2).  In a templated instructions file with multiple values files for differing kubernetes clusters, a user could have the value of `site` set differently in each values file, such as:

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
yaml_files:
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

## Author
[Andrew J. Huffman](https://github.com/ahuffman)


## License
[MIT](LICENSE)  
[NOTICE](NOTICE)


## Contributing
Please see our [Contribution Guide](CONTRIBUTING.md)


## Code of Conduct
Please see our project's [Code of Conduct](CODE-OF-CONDUCT.md)

## Communication

### E-Mail
Please join our mailing list on Google Groups: [yaml-overlay-tool-users](https://groups.google.com/g/yaml-overlay-tool-users)
