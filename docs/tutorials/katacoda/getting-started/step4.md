# Creating our first overlay

Make sure you have `yot.yaml`{{ open }} in your editor.  

Click the copy icon below to copy this block of code.

```yaml
  - name: prefix a single label
    query: metadata.labels.app~
    action: merge
    value: app.kubernetes.io/%v
```{{ copy }}

On a new line below `commonOverlays`, paste the block of code so it appears in your editor as such:

```yaml
commonOverlays:
  - name: prefix a single label
    query: metadata.labels.app~
    action: merge
    value: app.kubernetes.io/%v
```

---

Let's disect what we have just added to our yot.yaml.

1. `name` is an optional key that helps us document why we're doing this.
1. `query` is a JSONPath query to lookup data from within a YAML file or collection of YAML files. Notice the `~` character at the end of the query.  This will return the key of `app` rather than it's value of `guestbook`, a JSONPath notation.  The `query` key can also be a list of queries, which we'll look at later in the lab.
1. `action` is "how" to manipulate the data returned from our `query`.  `merge` is the default action and could be omitted in this case.  There are additional actions of `combine`, `delete`, and `replace`.
1. `value` is the intended value to update the data returned from the `query` with.  Yot allows for editing the original value with format markers.  In this instance the `%v` represents the original value of the data, which we've prefixed with `app.kubernetes.io/`.


Let's list out the files we'd like to manipulate in the next step.