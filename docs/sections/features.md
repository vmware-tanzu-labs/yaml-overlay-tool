[Back to Table of contents](../index.md)  


## YAML Overlay Tool features

Here's a quick introduction to YAML Overlay Tool's set of features, and why you might want to consider using the tool.


1. JSONPath queries
    * Yot uses JSONPath to lookup data from existing YAML documents within YAML files
    * Provides the user a familiar interface and a means to extend existing skills
    * Allows for convenient data manipulation
1. Ability to [`delete`](overlayActions.md#2-delete) unneeded data from YAML documents
1. Ability to [`inject`](overlayActions.md#2-inject) new data into YAML documents when a JSONPath query returns no results
1. Ability to [`merge`](overlayActions.md#3-merge) new data with existing data in YAML documents
1. Ability to [`replace`](overlayActions.md#4-replace) existing data with new data in YAML documents
1. Ability to [`combine`](overlayActions.md#1-combine) existing values returned from your JSONPath `query` with new values of the same type
    * Allows for integer addition
    * Allows for string concatenation
    * Allows for boolean addition
1. Ability to [remove existing comments](exampleUsage.md#remove-source-yaml-file-comments-prior-to-overlayment) from YAML documents prior to further manipulation of the data
1. Ability to manipulate the value of data returned from JSONPath queries
    * By the use of your JSONPath`query` and [format markers](formatMarkers.md) an end-user can take the returned original data and manipulate it with or without a `sed` implementation
1. Ability to both [*preserve*](comments.md#comment-preservation) and [*inject*](comments.md#comment-injection) new **comments** into a YAML document
    * This is useful for additional tooling that can consume comments
    * Useful when you want to leave end-users of the configuration notes as to *why* something is setup the way it is
1. Ability to perform [one-off overlays from the CLI](exampleUsage.md#use-without-an-instructions-file)
1. [Declarative YAML specification](instructionsFileIntro.md) (instructions file) to direct `yot` on how to manipulate the YAML data
1. Three levels of specificity in regards to how to manipulate YAML documents within YAML files
    1. [`commonOverlays`](instructionsFileSpec.md#top-level-commonoverlays-keys) affect **all** YAML documents with **all** YAML files
    1. [`overlays`](instructionsFileSpec.md#overlays-keys) affect **all** YAML documents within a **specific** YAML file
    1. [`documents`](instructionsFileSpec.md#documents-keys) affect a **specific** YAML document within a **specific** YAML file
1. Overlay [Qualifiers](overlayQualifiers.md) provide a way to *conditionally* apply overlays to a YAML document
    * Groups of `conditions` that must all match
    * Multiple `conditions` groups act as an implicit or, meaning if all the conditions of **one** of the groups match, the overlay will be applied
1. Ability to operate on a *single* YAML file or *entire directories* of YAML files from anywhere on the filesystem
1. Ability to redirect where an overlayed YAML file will be placed within the output directory with the [`outputPath` key](instructionsFileSpec.md#top-level-yamlfiles-keys)
    * When working with directories, a secondary listing of a specific filename can also be made to override where its final output location will be
1. Ability to use [Go templating and Go templating Sprig functions](instructionsFileTemplating.md) within the instructions file

[Back to Table of contents](../index.md)  
[Next Up: Installation and setup](setup.md)