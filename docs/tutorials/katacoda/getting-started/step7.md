# Improving our first overlay

As we saw in the previous steps we were able to update a single label with the `app.kubernetes.io/` prefix, but what if we want to update the same value in other JSONPaths within the YAML files?

Let's improve our first overlay.

---

The `query` key supports either a string (what we saw previously) or a list of JSONPath queries.

Let's switch our `query` to a list and add some additional queries to our overlay.

Click the copy icon to copy the new query code block:

```yaml
    query:
      - metadata.labels.app~
      - spec.selector.matchLabels.app~
      - spec.template.metadata.labels.app~
      - spec.selector.app~
```{{ copy }}

Go ahead and replace the line that is currently `query: metadata.labels.app~` by highlighting it and pasting the new code in your editor.

This new `query` in list format will replace all instances where `app` is a label key with `app.kubernetes.io/app` so that everything related to that label key will now function properly when you attempt to deploy it to a Kubernetes cluster.


Your `commonOverlays` should now look like this:

```yaml
commonOverlays:
  - name: prefix a single label
    query:
      - metadata.labels.app~
      - spec.selector.matchLabels.app~
      - spec.template.metadata.labels.app~
      - spec.selector.app~
    action: merge
    value: app.kubernetes.io/%v
```

Your editor will auto-save the yot.yaml, so let's go ahead and try the changes out:

`yot -i yot.yaml -s`{{ execute }}  

Let's continue building off of this example in the next steps.