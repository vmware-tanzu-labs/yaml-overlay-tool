# outputPath usage

There are times when you may want to reorganize how your files are output.  

Under each listing in the `yamlFiles` section of the Yot instructions file, an `outputPath` key is available.

Let's say we want to move each component from our source-manifests directory and move them to sub-directories.

Copy the following code block by clicking the copy icon.

```yaml
    outputPath: frontend/
  - name: Frontend service
    path: source-manifests/frontend-service.yaml
    outputPath: frontend/
  - name: Redis follower service
    path: source-manifests/redis-follower-service.yaml
    outputPath: redis-follower/
  - name: Redis follower deployment
    path: source-manifests/redis-follower-deployment.yaml
    outputPath: redis-follower/
  - name: Redis leader service
    path: source-manifests/redis-leader-service.yaml
    outputPath: redis-leader/
  - name: Redis leader deployment
    path: source-manifests/redis-leader-deployment.yaml
    outputPath: redis-leader/
```{{ copy }}


Paste the copied code on a new line under the last line of your `yamlFiles` section.

Your `yamlFiles` section should now look like this:

```yaml
yamlFiles:
  - name: source manifests
    path: source-manifests
  - name: frontend-deployment
    path: source-manifests/frontend-deployment.yaml
    overlays:
      - name: Ensure frontend-deployment has labels
        query: metadata.labels
        action: merge
        value:
          app.kubernetes.io/app: guesbook
          app.kubernetes.io/tier: frontend
        onMissing:
          action: inject
    documents:
      - path: 0
        overlays:
          - name: add another new label
            query: metadata.labels
            action: merge
            value:
              app.kubernetes.io/owner: me
    outputPath: frontend/
  - name: Frontend service
    path: source-manifests/frontend-service.yaml
    outputPath: frontend/
  - name: Redis follower service
    path: source-manifests/redis-follower-service.yaml
    outputPath: redis-follower/
  - name: Redis follower deployment
    path: source-manifests/redis-follower-deployment.yaml
    outputPath: redis-follower/
  - name: Redis leader service
    path: source-manifests/redis-leader-service.yaml
    outputPath: redis-leader/
  - name: Redis leader deployment
    path: source-manifests/redis-leader-deployment.yaml
    outputPath: redis-leader/
```

Run the following command to generate a newly organized copy of our source manifests:

`yot -i yot.yaml -o /tmp/organized`{{ execute }}

Let's take a look at the directory structure:

`tree /tmp/organized`{{ execute }}

Notice that our guestbook-namespace.yaml is in the root of the /tmp/organized directory.  This is because we did not redirect it with an `outputPath`.  Rather than just using directories, new filenames may be specified as well, or a combination thereof.

> **NOTE:** The output path is ignored when the `-s` option for stdout is specified.  Additionally, if the `outputPath` is not a fully-qualified path such as `/tmp/somedir/`, then the new directory structure will be relative to the output directory specified with the `-o` < output path > command-line option.