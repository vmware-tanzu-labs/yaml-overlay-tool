# Overlay qualifier introduction

An Overlay qualifier allows us to add additional conditions as to when an overlay should be applied by making use of a `query` and an expected `value`.

Copy the code block below by clicking on the copy icon.

```yaml
  - name: Prefix labels for Service selectors
    query: spec.selector.matchLabels
    action: merge
    value:
      app.kubernetes.io/%k: "%v"
    documentQuery:
      - conditions:
          - query: kind
            value: Service
```{{ copy }}

Paste this on a new line below your previous commonOverlay in the yot.yaml.

You should now have a `commonOverlays` section that looks like this:

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
```

Let's break this down in the next step and also try it out.