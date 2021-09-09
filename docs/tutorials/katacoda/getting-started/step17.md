# Comment injection

Yot is fairly unique in its ability to both preserve and inject comments.  You may have already noticed that comments have been preserved from the output of other steps in this lesson.

Comment injection is trivial with Yot, so let's take a look at how to achieve this.

Copy the following code block by clicking on the copy icon.

```yaml
  - name: Add a new label along with a comment
    query: metadata.labels
    action: merge
    value:
      app.kubernetes.io/environment: development  # this label will change when we're ready to productionalize this application
```{{ copy }}

Go ahead and paste the copied code block on a new line below the last overlay in your `commonOverlays` section within your yot.yaml.

Let's go ahead and run Yot:

`yot -i yot.yaml -s`{{ execute }}

As you can see we now have a new label along with a new comment.

You can also remove existing comments prior to making changes to your source manifests with the `--remove-comments` CLI option.

Let's take a look at running Yot with that option:

`yot -i yot.yaml -s --remove-comments`{{ execute }}

As you can see all "source" comments were removed and we are left with only the new comment that was injected.