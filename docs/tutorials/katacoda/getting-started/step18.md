# Adding a namespace to all of Kubernetes components

To get us ready for deploying our sample guestbook application, we still need to ensure that all of our Kubernetes components are namespaced.  

Copy the following code block by clicking on the copy icon.

```yaml
  - name: Ensure all components are namespaced
    query: metadata.namespace
    action: merge
    value: guestbook-application
    onMissing:
      action: inject
    documentQuery:
      - conditions:
          - query: kind
            value: Deployment
      - conditions:
          - query: kind
            value: Service
```{{ copy }}

On a new line below the last overlay in your yot.yaml `commonOverlays` section, paste the copied code block.

Finally, let's deploy our application to our Kubernetes cluster to complete the lesson:

1. Create the Kubernetes namespace:

`kubectl create ns guestbook-application`{{ execute }}

1. Run our overlay instructions and deploy it directly to Kubernetes:

`yot -i yot.yaml -s | kubectl apply -f -`{{ execute }}

1. Ensure our application has been deployed:

`kubectl get all -n guestbook-application`{{ execute }}