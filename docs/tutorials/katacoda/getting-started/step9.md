# Overlay qualifiers

In the previous step we were able to prefix the labels for almost all of the Kubernetes manifests.  However, we would have run into an issue with the `spec.selector` query.  This is because on a Kubernetes `Deployment` the labels are under `spec.selector.matchLabels`, and on a Kubernetes `Service` the labels are under `spec.selector`.  If we had both of those JSONPath queries listed in our `query`, the `Deployment`'s `spec.selector.matchLabels` key would be transformed to `spec.selector.app.kubernetes.io/matchLabels` and the data would have been lost.

---

That brings us to the topic of **Overlay qualifiers**.

An Overlay qualifier allows us to add additional conditions as to when an overlay should be applied by making use of a `query` and an expected `value`.

Copy the code block below by clicking on the copy icon.

```yaml
  - name: Prefix labels for Deployment selectors
    query: spec.selector.matchLabels
    action: merge
    value:
      app.kubernetes.io/%k: "%v"
    documentQuery:
      - conditions:
          - query: kind
            value: Deployment
```{{ copy }}

Paste this on a new line below your previous commonOverlay in the yot.yaml.

You should now have a `commonOverlays` section that looks like this:

```yaml
commonOverlays:
  - name: prefix a single label
    query:
      - metadata.labels
      - spec.selector.matchLabels
      - spec.template.metadata.labels
      - spec.selector
    action: merge
    value:
      app.kubernetes.io/%k: "%v"
  - name: Prefix labels for Deployment selectors
    query: spec.selector.matchLabels
    action: merge
    value:
      app.kubernetes.io/%k: "%v"
    documentQuery:
      - conditions:
          - query: kind
            value: Deployment
```

Let's break this down in the next step and also try it out.