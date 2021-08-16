# Improving our first overlay (continued)

Now that we have all instances of the label `app` replaced with `app.kubernetes.io/app` wouldn't it be nice to do the same for all of the other labels in one step?

Yes, we can do it all in one step!

We will be introducing another format marker, `%k`, which represents the existing key when JSONPath query returns a map/dictionary of key/value pairs.  The `%k` format marker can only be used on maps/dictionaries.

Once again, let's update our first `query`.

Click the copy icon to copy the updated query code block:

```yaml
    query:
      - metadata.labels
      - spec.selector.matchLabels
      - spec.template.metadata.labels
      - spec.selector
```{{ copy }}

Replace the existing query by pasting it into the yot.yaml.  

Notice that we've dropped `.app~` from each of the JSONPath queries in the list.

We need to make one more change our overlay.  This is where we'll see the `%k` format marker come into play.

Click the copy icon to copy the updated `value` code block:

```yaml
    value:
      app.kubernetes.io/%k: "%v"
```{{ copy }}

Replace the existing `value` line by pasting what you copied above.


Notice how the YAML for `value` is no longer on a single line.  This represents the updated value will be a map/dictionary.  

Furthermore, the `%k` format marker will be substituted with existing key, and the `%v` will be substituted with the existing value.  Since the data returned from our queries may have 1 or more keys, this update will apply to each of them automatically.

Your `commonOverlays` should now look like this:
```yaml
commonOverlays:
  - name: prefix all labels
    query:
      - metadata.labels
      - spec.selector.matchLabels
      - spec.template.metadata.labels
      - spec.selector
    action: merge
    value:
      app.kubernetes.io/%k: "%v"
```

Go ahead and try it out now:
`yot -i yot.yaml -s`{{ execute }}

All of your labels should now be prefixed!