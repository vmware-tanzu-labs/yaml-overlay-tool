# Lab information

In your lab environment, you will find a collection of source Kubernetes manifests that deploy a redis based guestbook application.  These example manifests were taken from the Kubernetes documentation here: [https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook](https://cloud.google.com/kubernetes-engine/docs/tutorials/guestbook)

These manifests can be found within `/root/source-manifests`.  

Have a look at the original files: `ls -1 /root/source-manifests`{{execute}}

---
You will also find that `yot`, the YAML Overlay Tool binary, has been installed to the environment.

You can see the available command usage by running `yot --help`{{execute}}.

The lab envioronment also has access to a Kubernetes cluster, so we will be available to deploy any of the manifests we manipulate later on with Yot.
