[Back to Project README](../README.md)


## Introduction to YAML Overlay Tool

YAML Overlay Tool (Yot) is not a traditional text-based templating tool. It is a YAML overlay/patching tool that takes fragments of YAML configuration to apply or inject over the top of an existing YAML configuration.  

Each overlay operation is performed with a JSONPath `query`, a desired `value`, and an `action`.  If a JSONPath `query` returns no results, a desired value is either ignored (default behavior), or injected, (`onMissing`). It can also provide a specific path or set of paths (`injectPath`) to inject the value, if the initial JSONPath query was not a fully-qualified JSONPath.  

## Table of contents
1. [Installation and setup](sections/setup.md)
1. [Command line interface usage and overview](sections/usage.md)
1. [Instructions file introduction](sections/instructionsFileIntro.md)
    - [Instructions file YAML specification](sections/instructionsFileSpec.md)
    - [Overlay actions](sections/overlayActions.md)
    - [Overlay qualifiers](sections/overlayQualifiers.md)
    - [Format markers](sections/formatMarkers.md)
    - [Comment preservation and injection](sections/comments.md)
    - [Using templating within instructions files](sections/instructionsFileTemplating.md)
1. [Order of operations/Processing order](sections/orderOfOperations.md)
1. [Output directory structure](sections/outputDirStructure.md)
1. [Common use-cases and patterns for Kubernetes](sections/useCasesForKubernetes.md)
1. [Interactive tutorials and learning paths](sections/tutorials.md)


[Back to Project README](../README.md)
