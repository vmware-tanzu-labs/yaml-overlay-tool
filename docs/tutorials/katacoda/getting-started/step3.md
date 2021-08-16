# YAML Overlay Tool instructions file

An empty Yot instructions file has been provided to get us started with building out our overlays.  Let's open it up in the editor and take a look. `yot.yaml`{{ open }}  

You'll notice that we have two lines that have been commented out.  Let's get started by removing the `# ` from each of those lines in your editor.

---

Let's discuss what these two lines mean.

`commonOverlays` allows us to do something to all `yamlFiles`.  It is a list of things to do.

`yamlFiles` is a list of the YAML files we would like to manipulate.  Each item in the list will have a `path` key that should be a path relative to the instructions file (yot.yaml).

In the next step we'll start to build overlays to address prefixing our Kubernetes labels.