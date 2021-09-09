# Creating our first overlay (continued)

Your editor will auto-save your work, so you should now have a yot.yaml that looks like this:

```yaml
---
commonOverlays:
  - name: prefix a single label
    query: metadata.labels.app~
    action: merge
    value: app.kubernetes.io/%v

yamlFiles:
  - name: source manifests
    path: source-manifests
```

---

In your terminal, let's test out our first overlay:  
`yot -i yot.yaml -s`{{ execute }}

Let's disect the command:

1. `yot` is the YAML Overlay Tool binary
1. `-i` represents a Yot instructions file with file `yot.yaml`
1. `-s` sends the output to standard out (stdout)

When you review the output you will see all the prior instances of `metadata.labels.app` are now prefixed with `app.kubernetes.io/` so that the key is now `app.kubernetes.io/app`.

Now try outputting the manipulated YAML to an output directory:  
`yot -i yot.yaml -o /tmp/new-manifests`{{ execute }}

The `-o` represents an output directory with a value of `/tmp/new-manifests`.  

Take a look at the contents of the new copy of your YAML files:  
`ls -1 /tmp/new-manifests`{{ execute }}

Take a look at one of the files:  
`cat /tmp/new-manifests/frontend-service.yaml`{{ execute }}

You should see that the `app` label is now `app.kubernetes.io/app: guestbook`

Let's continue to build off of this example now that you've created your first overlay.

