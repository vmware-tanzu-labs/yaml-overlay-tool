# File specific overlays

In the previous step we added missing labels within a common overlay using overlay qualifiers.  There's an easier way to achieve this by making use of file specific overlays.


Copy the following code block by clicking on the copy icon.

```yaml
  - name: Frontend deployment
    path: source-manifests/frontend-deployment.yaml
    overlays:
      - name: Ensure frontend-deployment has labels
        query: metadata.labels
        action: merge
        value:
          app.kubernetes.io/app: guesbook
          app.kubernetes.io/tier: frontend
        onMissing:
          action: inject
```{{ copy }}

On a new line under the `yamlFiles` section, paste the copied code block.

Your `yamlFiles` section should now look like this:

```yaml
yamlFiles:
  - name: source manifests
    path: source-manifests
  - name: Frontend deployment
    path: source-manifests/frontend-deployment.yaml
    overlays:
      - name: Ensure frontend-deployment has labels
        query: metadata.labels
        action: merge
        value:
          app.kubernetes.io/app: guesbook
          app.kubernetes.io/tier: frontend
        onMissing:
          action: inject
```

Be sure to remove the 3rd common overlay from the `commonOverlays` section which we created in the previous step so that there are only two overlays `commonOverlays`, and now looks like this:

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

Go ahead and run the following command to apply our changes:

`yot -i yot.yaml -s`{{ execute }}

You should see that `frontend-deployment.yaml` still has labels, but our instructions are expressed a little differently now.

Let's discuss this in the next step.