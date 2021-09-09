# Document specific overlays

Copy the following block of code by clicking the copy icon.

```yaml
    documents:
      - path: 0
        overlays:
          - name: add another new label
            query: metadata.labels
            action: merge
            value:
              app.kubernetes.io/owner: me
```{{ copy }}

Paste the copied code block on a new line at the end of the `yamlFiles` section of yot.yaml.


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
    documents:
      - path: 0
        overlays:
          - name: add another new label
            query: metadata.labels
            action: merge
            value:
              app.kubernetes.io/owner: me
```

Notice we did not add an `onMissing` action.  This is because `metadata.labels` was created in a previous step, and we are now just merging in a new key/value pair to the existing map.

Since the `frontend-deployment.yaml` only has a single YAML document, our document path is `0`.  If there were 2 documents, we'd refer to this path as `1` to make changes to the 2nd document, etc.  Other than seeing the `documents` key and the `path` key used with an integer, you should be familiar with the remainder of the Yot overlay specification.

Go ahead and apply the changes by running:

`yot -i yot.yaml -o /tmp/new-manifests`{{ execute }}

Let's take a look at the specific file we added a new label to:

`cat /tmp/new-manifests/frontend-deployment.yaml`{{ execute }}