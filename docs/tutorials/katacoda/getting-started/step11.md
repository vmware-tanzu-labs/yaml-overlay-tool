# Overlay qualifiers (continued)

Looking at the additional overlay we added to handle selectors on `Deployments`, let's describe what each new field does and how it works.

```yaml
  - name: Prefix labels for Service selectors
    query: spec.selector
    action: merge
    value:
      app.kubernetes.io/%k: "%v"
    documentQuery:
      - conditions:
          - query: kind
            value: Service
```

The addition of the `documentQuery` field opens the door to qualifying when to apply overlays with conditions.

Each item containing the key word of `conditions` allows us to have a grouping of conditions that must be met for an overlay to be applied.

In this instance we only have one condition listed, and that is to say when we `query` the `kind` key on a YAML document it must return a `value` of `Service`.  We could list additional conditions here as well that all must be met under the `conditions` key word.

Additionally, we could add a second listing of `conditions` for other scenarios where we would want to see this change applied.  Each group of `conditions` is treated as an **OR**, while each `query` and `value` under a `conditions` group acts as an **AND**.

Your yot.yaml should be auto-saved and we can now try it out.

Go ahead and run the following command to see all of our labels prefixed with `app.kubernetes.io/`.

`yot -i yot.yaml -s`{{ execute }}
