[Back to Project](https://github.com/vmware-tanzu-labs/yaml-overlay-tool)


# Introduction to YAML Overlay Tool

YAML Overlay Tool (Yot) is not a traditional text-based templating tool. Yot operates on structured YAML nodes. Yot is a YAML overlay or patching tool that primarily takes fragments of YAML configuration from within a declaritive specification to modify and manipulate an existing YAML configuration (manifest).  

Each overlay operation is performed with a JSONPath `query`, a desired `value`, and an `action`.  If a JSONPath `query` returns no results, a desired value is either ignored (default behavior), or injected (`onMissing`). Yot can also provide a specific path or set of paths (`injectPath`) to inject the value if the initial JSONPath query was not a fully-qualified JSONPath (e.g. using wildcards in the JSONPath query).  


# Table of contents
1. [YAML Overlay Tool features](sections/features.md)
1. [Installation and setup](sections/setup.md)
    - [Configuration file](sections/configFile.md)
    - [Environment variables](sections/envVars.md)
1. [Command line interface usage and overview](sections/commandUsage.md)
    - [Example CLI usage](sections/exampleUsage.md)
1. [Instructions file introduction](sections/instructionsFileIntro.md)
    - [Instructions file YAML specification](sections/instructionsFileSpec.md)
        - [Instructions file usage example](sections/instructionsFileSpec.md#instructions-file-usage-example)
    - [Overlay actions](sections/overlayActions.md)
    - [Overlay qualifiers](sections/overlayQualifiers.md)
    - [Format markers](sections/formatMarkers.md)
    - [Comment preservation and injection](sections/comments.md)
    - [Using templating within instructions files](sections/instructionsFileTemplating.md)
1. [Order of operations/processing order](sections/orderOfOperations.md)
1. [Output directory structure](sections/outputDirStructure.md)
1. [Common use-cases and patterns for Kubernetes](sections/useCasesForKubernetes.md)
1. [Interactive tutorials and learning paths](sections/tutorials.md)


[Back to Project](https://github.com/vmware-tanzu-labs/yaml-overlay-tool)
