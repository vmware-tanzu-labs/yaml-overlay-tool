# Working with Kubernetes labels

Yot is primarily a generic YAML manipulation tool, however, we will focus on using the tool in the context of Kubernetes for this lab.

## Recommended labels

Per the ["Recommended Labels" section of the Kubernetes documentation](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/), it is recommended to prefix a common shared label prefix of `app.kubernetes.io/` to each label for a Kubernetes component.

Let's take a look at one of the existing sample manifests for our guestbook application. `source-manifests/frontend-service.yaml`{{ open }}

Now that the file is opened in your editor, inspect the value of `metadata.labels` and `spec.selector`.  These are two locations within the manifest that contain labels.  Also notice that they do not contain the recommended `app.kubernetes.io/` prefix.

---

Let's take a look at another manifest.  `source-manifests/frontend-deployment.yaml`{{ open }}  

Inspect the value of `metadata.labels`.  There are not any labels here for some reason.  

Also inspect the value of `spec.selector.matchLabels`, `spec.template.metadata.labels`.  Again, we do not have the recommended `app.kubernetes.io/` prefix on these labels.

Let's learn how to quickly rectify this with Yot!