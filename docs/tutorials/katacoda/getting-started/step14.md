# Overlay specificity

As we saw in the last few steps, overlays can be applied at differing levels of specificity.  Yot has three levels of specificity where overlays can be applied.

1. `commonOverlays` apply to all files listed under `yamlFiles` and are applied first.

1. `overlays` listed under a file within `yamlFiles` are applied to all YAML documents within a YAML file (a YAML file may have more than one document separated by the `---` sequence).  These are applied after `commonOverlays`.  Think of these as common overlays for the file.

1. A third level of specificity can also be applied with the `documents` keyword under a file listed under the `yamlFiles` section.  This allows us to make changes to a specific YAML document within a YAML file.  These are applied after file specific overlays.

Now that we've seen common overlays and file specific overlays, let's see how to use a document specific overlay in the next step.