## Introduction to YAML Overlay Tool

YAML Overlay Tool (`yot`) is not a templating tool in the sense of a traditional text-based templating tool.  `yot` is primarily an overlayment tool, meaning we take fragments of YAML configuration and apply or inject them over the top of an existing YAML configuration.  

Each overlay operation is performed with a JSONPath `query`, a desired `value`, and an `action`.  If a JSONPath `query` returns no results, a desired value can be either ignored (default behavior) or injected (`onMissing`) and even provide a specific path or set of paths (`injectPath`) to inject the value if the initial JSONPath query was not a fully-qualified JSONPath.  

## Table of Contents
1. [Installation and Setup](sections/setup.md)
2. [Command Line Usage and Overview](sections/usage.md)
3. [Instructions File YAML Schema](sections/instructionsFile.md)
4. [Order of Operations/Processing](sections/orderOfOperations.md)
5. [Kubernetes Common Use Cases and Patterns](sections/kubernetesUseCases.md)
6. [Additional Examples and Use Cases](sections/examples.md)
7. [Interactive Tutorials and Learning Paths](sections/tutorials.md)