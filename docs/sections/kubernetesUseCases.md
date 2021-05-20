[Back to Table of Contents](../documentation.md)

# Kubernetes common use-cases and patterns

The following set of examples will help you quickly achieve common tasks in the context of Kubernetes YAML manifests.  

All these examples are available in [examples/kubernetes](../../examples/kubernetes) and are intended to be launched from the root of your local copy of the YAML Overlay Tool repository.

Example:

```bash
yot -i examples/kubernetes/< example you wish to run>.yaml -o < desired output path >
```

<br/>


### Example Kubernetes YAML manifests

Within the [examples/kubernetes/manifests](../../examples/kubernetes/manifests) directory of the YAML Overlay Tool repository, you'll find the two example Kubernetes YAML Manifests that are manipulated in the following set of example use-cases:

```yaml
# my-app.yaml
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    my.custom.annotation/fake: idk
  labels:
    name: my-web-page
  name: my-web-page
  namespace: my-web-page
spec:
  containers:
    - image: nginx:latest
      name: my-web-page
      ports:
        - containerPort: 443
      resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always

```

```yaml
# my-service.yaml
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
  labels:
    name: my-web-page
  name: my-service
  namespace: my-web-page
spec:
  ports:
    - name: 8443-443
      port: 8443
      protocol: TCP
      targetPort: 443
  selector:
    app: my-service
  type: LoadBalancer

```

<br/>


## Add additional labels and selectors to all YAML files in a directory

In this example, the `merge` action lets you add two new labels to a set of YAML files contained in a directory.

```yaml
---
# addLabels.yaml
commonOverlays:
  - name: Add additional labels
    query:
      - metadata.labels
      - spec.template.metadata.labels
      - spec.selector.matchLabels
      - spec.selector
    value:
      app.kubernetes.io/owner: Jeff Smith
      app.kubernetes.io/purpose: static-webpage
    action: merge

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now, you can apply the changes by generating a new set of YAML files.

`yot -i ./examples/kubernetes/addLabels.yaml -o /tmp/new`


## Prepend a private registry URL to all container images

In this example the `format` action lets you take images that are currently pointing to the Docker Hub registry (only base imagename:tag), and prepending with a private registry URL.

```yaml
# privateContainerRegistry.yaml
commonOverlays:
  - name: Set our private container registry in manifests
    query: ..image
    value: my-private-reg/%s
    action: format

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now, you can apply the changes by generating a new set of YAML files.

`yot -i ./examples/kubernetes/privateContainerRegistry.yaml -o /tmp/new`


## Modify the label key name

In this example, the `name` label key of `app.kubernetes.io/name` is manipulated by using the `format` action and retaining the existing value.  The `~` character in JSONPath+ always returns the value of the key, rather than the key/value pair.

```yaml
# formatLabelKey.yaml
commonOverlays:
  - name: Update name label's key to app.kubernetes.io/name
    query: metadata.labels.name~
    value: app.kubernetes.io/%s
    action: format

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now, you can apply the changes by generating a new set of YAML files.

`yot -i ./examples/kubernetes/formatLabelKey.yaml -o /tmp/new`

## Replace the label key name

In this example, the `name` label is replaced with `my-new-label` by using the `replace` action and retaining the existing value. The `~` character in JSONPath+ always returns the value of the key, rather than the key/value pair.

```yaml
# replaceLabelKey.yaml
commonOverlays:
  - name: Replace name label's key to my-new-label
    query: metadata.labels.name~
    value: my-new-label
    action: replace

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now, you can apply the changes by generating a new set of YAML files.

`yot -i ./examples/kubernetes/replaceLabelKey.yaml -o /tmp/new`


## Remove all annotations

There are times when annotations that are set for specific environments may not apply to your environment. To remove annotations, use the `delete` action.

```yaml
# removeAnnotations.yaml
commonOverlays:
  - name: Remove all annotations
    query: metadata.annotations
    action: delete

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now, you can apply the changes by generating a new set of YAML files.

`yot -i ./examples/kubernetes/removeAnnotations.yaml -o /tmp/new`


## Remove annotations from specific Kubernetes object types

To build on the previous example, there are times when you may want to remove annotations from specific Kubernetes types, or a combination of conditions.  To remove these annotations, use the `delete` action and a `documentQuery`.

```yaml
# removeAnnotationsWithConditions.yaml
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
    path: ./examples/kubernetes/manifests
```

Now, you can apply the changes by generating a new set of YAML files.

`yot -i ./examples/kubernetes/removeAnnotationsWithConditions.yaml -o /tmp/new`


[Back to Table of Contents](../documentation.md)  
[Next Up: Interactive Tutorials and Learning Paths](tutorials.md)
