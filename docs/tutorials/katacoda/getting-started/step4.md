# Creating our first overlay

Make sure you have `yot.yaml` {{ open }} in your editor.  

```yaml
  - name: prefix all of our labels
    query: metadata.labels
    action: merge
    value: app.kubernetes.io/%v
```{{ copy }}

