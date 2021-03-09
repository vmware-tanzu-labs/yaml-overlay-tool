# Philosophy Behind Management of YAML with 'YAML Overlay Tool'

There are many ways to manage YAML files today.  In fact so many, we're not going to name them.  This is not a new problem, and one could say that there are already great ways of managing YAML files when being used for the purpose of configuration or a desired state.  

Probably the most common tool that comes to mind in a Kubernetes context is Helm. Helm is great at what it does, but what becomes a problem in the world of templating YAML manifests, is the long-term management of those templates.  Another common problem is readability or the complexity that can be involved, and over time, this can be a serious problem to manage.

What if there was another way to accomplish this?  What if we decided to treat manifests as source code?  What would this look like?
