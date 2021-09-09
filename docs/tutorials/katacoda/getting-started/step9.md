# Improving our first overlay (continued)

In the previous step we were able to prefix the labels for almost all of the Kubernetes manifests.  However, we would have run into an issue with the `spec.selector` query.  

On a Kubernetes `Deployment` the labels are under `spec.selector.matchLabels`, and on a Kubernetes `Service` the labels are under `spec.selector`.  If we had both of those JSONPath queries listed in our `query`, the `Deployment`'s `spec.selector.matchLabels` key would be transformed to `spec.selector.app.kubernetes.io/matchLabels` and the data would have been lost.

That brings us to the topic of **Overlay qualifiers**.