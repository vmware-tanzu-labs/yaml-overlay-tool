# onMissing actions introduction

Earlier in the lesson we saw a Kubernetes `Deployment` that did not have any `metadata.labels`.  Let's see how we can ensure we add labels in a more advanced use case.

Yot, by default, does not act if a `query` returns no results.  However, there are times when we may want to add data.  We address this with an `onMissing` `action`.

Copy the following code block by clicking the copy icon.

```yaml
  - name: Ensure frontend-deployment has labels
    query: metadata.labels
    action: merge
    value:
      app.kubernetes.io/app: guesbook
      app.kubernetes.io/tier: frontend
    onMissing:
      action: inject
    documentQuery:
      - conditions:
          - query: kind
            value: Deployment
          - query: metadata.name
            value: frontend
```{{ copy }}

Paste this code block on a new line under your last `commonOverlays` item.

Your commonOverlays section should now a total of three overlays, and look like this:

```yaml
commonOverlays:
  - name: prefix labels
    query:
      - metadata.labels
      - spec.selector.matchLabels
      - spec.template.metadata.labels
    action: merge
    value:
      app.kubernetes.io/%k: "%v"
  - name: Prefix labels for Service selectors
    query: spec.selector
    action: merge
    value:
      app.kubernetes.io/%k: "%v"
    documentQuery:
      - conditions:
          - query: kind
            value: Service
  - name: Ensure frontend-deployment has labels
    query: metadata.labels
    action: merge
    value:
      app.kubernetes.io/app: guestbook
      app.kubernetes.io/tier: frontend
    onMissing:
      action: inject
    documentQuery:
      - conditions:
          - query: kind
            value: Deployment
          - query: metadata.name
            value: frontend
```

Go ahead and run the following command to see the changes to the `frontend-deployment.yaml`:

`yot -i yot.yaml -s`{{ execute }}

Let's take a look at how we could have made this change without the use of a common overlay in the next step.