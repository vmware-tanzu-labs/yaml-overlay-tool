[Back to Table of Contents](../documentation.md)  


# Output Directory Structure

Since Yot's Instructions File allows for operating on YAML files from many different paths on the file system, and additionally (**v0.2.0**) templated Instructions Files, Yot will output files to the file system with certain considerations.

1. To allow for multi-site templating (**v0.2.0**), Yot outputs its files to the output directory specified at run-time and within a sub-directory called `yamlFiles`, and rendered (templated) instructions files to a sub-directory called `renderedInstructions`.

1. To accomodate files that ***could*** have an identical basename, Yot calculates the least common path, and will recreate directory structures from there as needed.  If this were not the case, a file name `~/yamls/test.yaml` and `/tmp/manifests/test.yaml` would overwrite one another if specified in the Instructions File.

## Output Directory Structure Example

```bash
output
└── yamlFiles
    ├── Users
    │   └── ahuffman
    │       └── yamls
    │           └── test.yaml
    └── tmp
        └── manifests
            └── test.yaml
```


[Back to Table of Contents](../documentation.md)  
[Next Up: Kubernetes Common Use-Cases and Patterns](useCasesForKubernetes.md)