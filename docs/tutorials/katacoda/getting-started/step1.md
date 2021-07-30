# Lab information

In your lab environment, you will find a collection of source Kubernetes manifests that deploy a redis based guestbook application.  These manifests were taken from the Kubernetes documentation here: [https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook](https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook)

These manifests can be found within `~/source-manifests`.

You will also find that `yot`, the YAML Overlay Tool binary, has been installed to the environment.

You can see the available Yot usage by running `yot --help`{{execute}}.

The lab envioronment also has access to a Kubernetes cluster, so we will be available to deploy any of the manifests we manipulate later on with Yot.
