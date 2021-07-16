[Back to Overlay actions](overlayActions.md#3-merge)  
[Back to Overlay qualifiers](overlayQualifiers.md)  
[Back to Table of contents](../index.md) 


# Format markers


<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Format marker introduction](#format-marker-introduction)
- [Format marker parameters](#format-marker-parameters)
- [Manipulation of content with format markers](#manipulation-of-content-with-format-markers)
  - [Content manipulation examples](#content-manipulation-examples)
    - [Simple sed string replacement](#simple-sed-string-replacement)
    - [sed read command example](#sed-read-command-example)
    - [Forcing types with YAML tags](#forcing-types-with-yaml-tags)
      - [Problem](#problem)
      - [Solution](#solution)
      - [YAML tags reference](#yaml-tags-reference)

<!-- /code_chunk_output -->


## Format marker introduction

Format markers provide a Yot user the ability to access the original values returned from the JSONPath query.  This can be used to add some additional text to the original value, and they can be used more than once within the overlay's `value` key.  

## Format marker parameters

| Marker | Description | Action where available | Status |
| --- | --- | --- | --- |
| %k | Original key's value returned from the JSONPath `query`. If querying for a key by appending the `~` character, use `%v` | combine (comments only), merge, replace (comments only), and currently only available on scalar types | experimental in (v0.1.0) |
| %v | Original value returned from the JSONPath `query` | combine, merge, replace | stable |
| %h | Original value of the head comment (comment above original value) returned from the JSONPath `query` | combine (comments only), merge, replace (comments only) | experimental in (v0.1.0) |
| %l | Original value of the line comment (comment on the same line as the value) | combine (comments only), merge, replace (comments only) | stable |
| %f | Original value of the foot comment (comment below original value) | combine (comments only), merge, replace (comments only) | experimental in (v0.1.0) |


## Manipulation of content with format markers

Added in v0.6.0.  The values returned from a query can be further customized, allowing a Yot user to manipulate the original value of the content.  This can be achieved by appending a `sed` command within curly braces `{}` directly after a format marker.  The supported `sed` features can be found within the [Go-Sed](https://github.com/rwtodd/Go.Sed) package documentation, but the main difference is the implementation of regular expressions.

It is encouraged to familiarize yourself with the capabilities of [sed](https://www.gnu.org/software/sed/manual/sed.html#sed-scripts).

>**NOTE:** It is important to call out that these can only be used on [YAML scalars](https://www.tutorialspoint.com/yaml/yaml_scalars_and_tags.htm).


### Content manipulation examples

#### Simple sed string replacement

Using `%v{s/hello/world/g}` within the `value` key of an overlay would search the original value returned by the `query` and replace the text 'hello' with 'world'.  

Here's a more detailed example of this behavior.

```yaml
# examples/kubernetes/manifests/my-app.yaml
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    my.custom.annotation/fake: idk
  labels:
    name: my-web-page
    app: my-web-page
    owner: my-web-page
  name: my-web-page
  namespace: my-web-page
spec:
  containers:
    - image: nginx:latest
      name: my-web-page
      ports:
        - containerPort: 443
      resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
```

```yaml
# instructions.yaml file
---
yamlFiles:
  - path: "examples/kubernetes/manifests/my-app.yaml"
    overlays:
      - name: "replace 'page' with 'site' in original label contents"
        query: "metadata.labels"
        value:
          "%k": "%v{s/page/site/g}"
        action: "merge"
```

Which will produce the following output when run with `yot -i instructions.yaml -s`:

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    my.custom.annotation/fake: idk
  labels:
    name: my-web-site
    app: my-web-site
    owner: my-web-site
  name: my-web-page
  namespace: my-web-page
spec:
  containers:
    - image: nginx:latest
      name: my-web-page
      ports:
        - containerPort: 443
      resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
```


#### sed read command example

`%v{r ~/myfile.yaml}` would insert the contents of 'myfile.yaml' on a new line after the existing value.  This could be useful for populating Kubernetes Secret and ConfigMap contents.


#### Forcing types with YAML tags

##### Problem

There are times when you may want to simply inject a comment into an existing YAML document while retaining a key's original value.  However, problems can arise if the original type was an integer.  This scenario would look like this:

```yaml
...
- name: inject a comment
  query: some.jsonpath
  value: "%v"  # new comment
  action: merge
```

Since format markers begin with the `%` character, and the YAML specification does not allow a `%` character to begin a value, we must wrap the format marker in quotation marks.  Once you have wrapped the format marker in quotation marks, you have affectively transformed the original integer (pretend we were manipulating an integer) type to that of a string.


##### Solution

To solve the above problem, we must use explicit YAML tags to force the value into a desired type.  This can be achieved by doing the following:

```yaml
...
- name: inject a comment
  query: some.jsonpath
  value: !!int "%v"  # new comment
  action: merge
```


##### YAML tags reference

Acceptable YAML tags are:

| Tag | Related Type |
| --- | --- |
| `!!int` | Integer |
| `!!float` | Floating point decimal |
| `!!null` | `null` value |
| `!!map` | Map/Dictionary |
| `!!seq` | Sequence/Array |
| `!!str` | String |


[Back to Overlay actions](overlayActions.md#3-merge)  
[Back to Overlay qualifiers](overlayQualifiers.md)  
[Back to Table of contents](../index.md)  
[Next Up: Comment preservation and injection](comments.md)