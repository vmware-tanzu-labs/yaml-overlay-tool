[Back to Project README](../README.md)


## Introduction to YAML Overlay Tool

YAML Overlay Tool (Yot) is not a traditional text-based templating tool. It is an overlay tool that takes fragments of YAML configuration to apply or inject over the top of an existing YAML configuration.     

Each overlay operation is performed with a JSONPath `query`, a desired `value`, and an `action`.  If a JSONPath `query` returns no results, a desired value is either ignored (default behavior), or injected, (`onMissing`). It can also provide a specific path or set of paths (`injectPath`) to inject the value, if the initial JSONPath query was not a fully-qualified JSONPath.  

## Table of Contents
1. [Installation and Setup](sections/setup.md)
1. [Command Line Usage and Overview](sections/usage.md)
1. [Instructions File YAML Schema](sections/instructionsFile.md)
1. [Overlay Actions](sections/actions.md)
1. [Overlay Qualifiers](sections/qualifiers.md)
1. [Comment Preservation and Injection](sections/comments.md)
1. [Order of Operations/Processing](sections/orderOfOperations.md)
1. [Kubernetes Common Use-Cases and Patterns](sections/kubernetesUseCases.md)
1. [Interactive Tutorials and Learning Paths](sections/tutorials.md)


[Back to Project README](../README.md)
