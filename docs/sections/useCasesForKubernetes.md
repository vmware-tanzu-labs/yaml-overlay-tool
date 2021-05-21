[Back to Table of Contents](../documentation.md)

# Kubernetes common use-cases and patterns

The following set of examples will help you quickly achieve common tasks in the context of Kubernetes YAML manifests.  

All examples are available for your convenience in [examples/kubernetes](../../examples/kubernetes), and are intended to be launched from the root of your local copy of the YAML Overlay Tool repository:

```bash
yot -i examples/kubernetes/< example you wish to run >.yaml -o < desired output path >
```

<br/>


### Example Kubernetes YAML Manifests

Within the [examples/kubernetes/manifests](../../examples/kubernetes/manifests) directory of the YAML Overlay Tool repository, you will find the two example Kubernetes YAML Manifests which we will be manipulating in the following set of example use-cases:

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


## Adding additional labels and selectors to all YAML files in a directory

In the following example, the `merge` action adds 2 new labels to a set of YAML files contained in a directory.

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

Now, apply the changes by generating a new set of YAML files.  
>`yot -i ./examples/kubernetes/addLabels.yaml -o /tmp/new`


## Prepend a private registry URL to all container images

In the following example, the `merge` action takes the images which are currently pointing to the Docker Hub registry (only base imagename:tag), and prepends them with a private registry URL.  

This example also demonstrates how the `%v` format marker is used to inject the original value into a new line comment.

```yaml
---
# privateContainerRegistry.yaml
commonOverlays:
  - name: Set our private container registry in manifests
    query: ..image
    value: my-private-reg/%v  # old value was: %v
    action: merge

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now, apply the changes by generating a new set of YAML files.  
>`yot -i ./examples/kubernetes/privateContainerRegistry.yaml -o /tmp/new`


## Modify the name of a label's key

In the following example, the `name` label key and `app.kubernetes.io/name` are manipulated by using the `merge` action and retaining the existing value.  The `~` character in JSONPath+ always returns the value of the key, rather than the value of the key/value pair.

```yaml
---
# formatLabelKey.yaml
commonOverlays:
  - name: Update name label's key to app.kubernetes.io/name
    query: metadata.labels.name~
    value: app.kubernetes.io/%v
    action: merge

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now apply the changes by generating a new set of YAML files:  
>`yot -i ./examples/kubernetes/formatLabelKey.yaml -o /tmp/new`

## Replace the name of a label's key

In the following example, the `name` label is replaced with `my-new-label` by using the `replace` action and retaining the existing value. The `~` character in JSONPath+ always retures the value of the key, rather than the value of the key/value pair.

```yaml
---
# replaceLabelKey.yaml
commonOverlays:
  - name: Replace name label's key to my-new-label
    query: metadata.labels.name~
    value: my-new-label  # old label key was: %v
    action: replace

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now, apply the changes by generating a new set of YAML files.  
>`yot -i ./examples/kubernetes/replaceLabelKey.yaml -o /tmp/new`


## Remove all annotations

Annotations are set for specific environments that may not apply to your environment. 

Use the `delete` action to remove all annotations in the following example:

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

Now, apply the changes by generating a new set of YAML files.  
>`yot -i ./examples/kubernetes/removeAnnotations.yaml -o /tmp/new`


## Remove annotations from specific kubernetes object types

To build on the previous example, there are times when you may want to remove annotations from specific Kubernetes types, or a combination of conditions.  To remove annotations, use the `delete` action and a `documentQuery`.

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

Now, apply the changes by generating a new set of YAML files.  
>`yot -i ./examples/kubernetes/removeAnnotationsWithConditions.yaml -o /tmp/new`


## Inject comments

Adds additional comments about why you did a specific task. You can also use it to document specific settings in the event you later want to restore them to a previous state.

Also lets you add comments as annotations for consumption by another application.

Yot **can** inject comments!

```yaml
---
# injectComments.yaml
commonOverlays:
  - name: Inject line comments via merge
    query:
      - metadata.annotations['my.custom.annotation/fake']
      - metadata.annoations['service.beta.kubernetes.io/aws-load-balancer-type']
    value: "%v" # inject a line comment
    action: merge
  - name: Inject a line comment via replace
    query: spec.containers[0].image
    value: new-image:latest # old value was: %v
    action: replace
  - name: Inject head, foot, and line comments via merge
    query: metadata.labels
    value:
      # inject a head comment
      app.kubernetes.io/owner: Jeff Smith  # inject a line comment
      app.kubernetes.io/purpose: static-webpage  # inject another line comment
      # inject a foot comment
    action: merge
yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now, apply the changes by generating a new set of YAML files.  
>`yot -i ./examples/kubernetes/injectComments.yaml -o /tmp/new`


[Back to Table of Contents](../documentation.md)  
[Next Up: Interactive Tutorials and Learning Paths](tutorials.md)
