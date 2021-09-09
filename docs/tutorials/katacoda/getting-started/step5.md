# Creating our first overlay (continued)

In the previous step, we created a `commonOverlay`.  Now we need to list out the files it applies to.  

Click the copy icon below to copy the code block.

```yaml
  - name: source manifests
    path: source-manifests
```{{ copy }}

Paste the copied code on a new line below the `yamlFiles:` line so that is looks like this:

```yaml
yamlFiles:
  - name: source manifests
    path: source-manifests
```

---

Let's disect what this means.

1. `name` is an optional key to describe the files we're working on.
1. `path` is a required key that describes the files to manipulate relative to the Yot instructions file (yot.yaml).  This value can be either a single file or an entire directory.  In this case we've declared that we'd like to operate on an entire directory.

There are additional keys available under each `yamlFiles` item that we'll introduce in later steps.
