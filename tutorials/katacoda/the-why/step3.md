# YAML Overlay Tool's Philosophy (continued)

What advantages would treating manifests as source code provide?

Some Kubernetes application vendors provide manifests for the installation of their software via static files.  That is to say that there is no templating involved.  These manifests contain all the common configuration to allow for running the software on a vanilla Kubernetes cluster.

Some vendors provide a templated version of their Kubernetes manifests through one or more technologies, and also provide manifests for specific Kubernetes distributions/platforms.

This is all good stuff, but what if we want to further customize them?  What if we have built a custom platform that has more specific requirements to get an application running than a vanilla Kubernetes cluster?  Wouldn't it be nice to see what a standard working and tested configuration should look like without a scripting/templating language intermingled within its configuration?
