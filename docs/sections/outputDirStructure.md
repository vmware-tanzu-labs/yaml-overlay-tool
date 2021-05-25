[Back to Table of Contents](../documentation.md)  


# Output Directory Structure

Since Yot's Instructions File allows for operating on YAML files from many different paths on the file system, and additionally (**v0.2.0**) templated Instructions Files, Yot will output files to the file system with certain considerations.

1. To accomodate files that ***could*** have an identical basename, Yot calculates the least common path, and will recreate directory structures from there as needed.  If this were not the case, a file name `~/yamls/test.yaml` and `/tmp/manifests/test.yaml` would overwrite one another if specified in the Instructions File.

## Output Directory Structure Example

```bash
output
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