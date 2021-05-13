## YAML Overlay Tool - Common Kubernetes Use-Cases

### Adding Additional Labels and Selectors To All YAML Files In a Directory

In this example we use the `merge` action to add in 2 new labels to a set of YAML files contained in a directory.

```yaml
---
# instructions.yaml
commonOverlays:
  - name: Add additional labels
    query: 
      - metadata.labels
      - spec.template.metadata.labels
      - spec.selector.matchLabels
    value:
      newlabel1: newValue1
      newlabel2: newValue2
    action: merge

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: /tmp/k8s
```

Now apply the changes by generating a new set of YAML files:
`yot -i instructions.yaml -o /tmp/new`


### Prepend a Private Registry URL to All Container Images

In this example we use the `format` action to take the images which are currently pointing to the Docker Hub registry (only base imagename:tag), and prepending with a private registry URL.

```yaml
# instructions.yaml
commonOverlays:
  - name: Set our private container registry in manifests
    query: ..image
    value: my-private-reg/%s
    action: format

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: /tmp/k8s
```

Now apply the changes by generating a new set of YAML files:
`yot -i instructions.yaml -o /tmp/new`


### Modify the Name of a Label's Key

In this example we will manipulate the `name` label key with `app.kubernetes.io/name` by using the `format` action and retaining the existing value.  The `~` character in JSONPath+ always returns the value of the key, rather than the value of the key/value pair.

```yaml
# instructions.yaml
commonOverlays:
  - name: Update name label's key to app.kubernetes.io/name
    query: metadata.labels.name~
    value: app.kubernetes.io/%s
    action: format

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: /tmp/k8s
```

Now apply the changes by generating a new set of YAML files:
`yot -i instructions.yaml -o /tmp/new`

### Replace the Name of a Label's Key

In this example we will replace the `name` label with `my-new-label` by using the `replace` action and retaining the existing value. The `~` character in JSONPath+ always retures the value of the key, rather than the value of the key/value pair.

```yaml
# instructions.yaml
commonOverlays:
  - name: Replace name label's key to my-new-label
    query: metadata.labels.name~
    value: my-new-label
    action: replace

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: /tmp/k8s
```

Now apply the changes by generating a new set of YAML files:
`yot -i instructions.yaml -o /tmp/new`


### Remove All Annotations

Often times annotations are set for certain environments that may not apply to your environment.  To remove all annotations we will use the `delete` action.

```yaml
# instructions.yaml
commonOverlays:
  - name: Remove all annotations
    query: metadata.annotations
    action: delete

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: /tmp/k8s
```

Now apply the changes by generating a new set of YAML files:
`yot -i instructions.yaml -o /tmp/new`


### Remove Annotations from Specific Kubernetes Object Types

To build on the previous example, there are often times when you may want to remove annotations from specific Kubernetes types, or a combination of conditions.  To remove these annotations, we will use the `delete` action and a `documentQuery`.

```yaml
commonOverlays:
  - name: Remove all annotations with conditions
    query: metadata.annotations
    action: delete
    documentQuery:
      - conditions:
          - key: kind
            value: Service
          - key: metadata.namespace
            value: my-web-page
      - conditions:
          - key: metadata.name
            value: my-service

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: /tmp/k8s
```