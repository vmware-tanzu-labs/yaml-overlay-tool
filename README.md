# YAML Overlay Tool (yot)

YAML Overlay Tool, or `yot` for short, and often pronounced *yaucht*, is a tool to assist with patching YAML files.  Yot uses JSONPath to query YAML documents within YAML files, and to perform a change.  YAML Overlay Tool operates on YAML nodes, and therefore is able to preserve and inject head, foot, and line comments into the newly output versions of the files you manipulate.


## Why Create Another YAML Tool?

`yot` was designed to be flexible, simple, and familiar with a focus on end-user and developer experience.  Whether you want the ability to use a templating language to transform YAML data, or just change a couple values in a YAML document, `yot` will make it possible.  However, as we port the protype Python version of `yot` to be written in Go, we have temporarily stripped out the templating language functionality, and plan to have a robust set of options available in version 0.2.0.  

Our philosophy is to treat YAML manifests as source code.  We don't want to manage templated YAML files, we want to manage patches (overlays) for YAML files, and furthermore we want to keep any potential templating outside of the source YAML files altogether.  

Templated files become hard to manage over time and can be difficult to read.  `yot` allows us to take YAML documents from multiple sources and transform them through overlays to fit our environment's requirements.  This allows us to take any generic YAML file and manipulate it to suit our purpose, without contaminating the original file.

This practice gets you out of the cycle of updating and managing complex YAML templates.  At the same time `yot`'s instructions file specification serves as documentation-as-code, where you have now essentially documented all the required changes to source YAML files in one place, and what is required to get an application running in your environment.

The use of JSONPath queries and templating give the tool familiar interfaces, making adoption easier, and a more pleasant end-user experience.  The specification, also known as the *instructions file*, is assembled in a declarative way, where we only operate on what has been clearly defined.  We take *actions* based on JSONPath query results.  We provide flexibility by allowing your instructions to be templated if needed (functionality will return in v0.2.0).  Please see the [full documentation](docs/documentation.md), which will help get you moving along with `yot`!


## Author

[Andrew J. Huffman](https://github.com/ahuffman)  
[Jeff Davis](https://github.com/JefeDavis)


## License

[MIT](LICENSE)  
[NOTICE](NOTICE)


## Contributing

Please see our [Contribution Guide](CONTRIBUTING.md)


## Code of Conduct

Please see our project's [Code of Conduct](CODE-OF-CONDUCT.md)


## Communication
### E-Mail

Please join our mailing list on Google Groups: [yaml-overlay-tool-users](https://groups.google.com/g/yaml-overlay-tool-users)


## Full Documentation and User Guide

[docs/documentation.md](docs/documentation.md)
