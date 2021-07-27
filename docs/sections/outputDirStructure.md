[Back to Table of contents](../index.md)  


# Output directory structure

Since Yot's Instructions File allows for operating on YAML files from many different paths on the file system, and additionally (**v0.3.0**) templated Instructions Files, Yot will output files to the file system with certain considerations.

1. To accomodate files that ***have*** an identical basename, Yot calculates the least common path, and will recreate directory structures from there as needed within the output directory specified.  If this were not the case, a file name `~/yamls/test.yaml` and `/tmp/manifests/test.yaml` would overwrite one another if specified in the instructions file if outputting to a flat output directory structure.

## Output directory structure example

Based on the scenario laid out above, the two files named `test.yaml` from differing locations on the filesystem would be output as follows:

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


[Back to Table of contents](../index.md)  
[Next Up: Common use-cases and patterns for Kubernetes](useCasesForKubernetes.md)