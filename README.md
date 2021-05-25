[![Go Reference](https://pkg.go.dev/badge/github.com/vmware-tanzu-labs/yaml-overlay-tool.svg)](https://pkg.go.dev/github.com/vmware-tanzu-labs/yaml-overlay-tool)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/vmware-tanzu-labs/yaml-overlay-tool)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/vmware-tanzu-labs/yaml-overlay-tool)](https://goreportcard.com/report/github.com/vmware-tanzu-labs/yaml-overlay-tool)
[![GitHub](https://img.shields.io/github/license/vmware-tanzu-labs/yaml-overlay-tool)](https://github.com/vmware-tanzu-labs/yaml-overlay-tool/blob/main/LICENSE)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/vmware-tanzu-labs/yaml-overlay-tool)](https://github.com/vmware-tanzu-labs/yaml-overlay-tool/releases)
# YAML Overlay Tool (Yot)

YAML Overlay Tool, or Yot for short, often pronounced */yaucht/*, is a tool to assist with patching YAML files.  Yot uses JSONPath to query YAML documents within YAML files, and to perform a change.  YAML Overlay Tool operates on YAML nodes. It is able to preserve and inject head, foot, and line comments into the new output versions of the files that you manipulate.


## Why Create Another YAML Tool?

Yot is designed to be flexible, simple, and familiar; with a focus on end-user and developer experience.  Whether you want to use a templating language to transform YAML data, or just change a couple values in a YAML document, Yot can make it possible.  

**NOTE** In order to port the protype Python version of Yot to be written in Go, we temporarily stripped out the templating language functionality. A robust set of options is planned for version v0.3.0.  

Our philosophy is to treat YAML manifests as source code. We don't want to manage templated YAML files. We want to manage patches (overlays) and keep all potential templating outside of the source YAML files.  

Templated files are hard to manage over time, making them difficult to read.  Yot allows us to take YAML documents from multiple sources and transform them to fit our environment requirements through overlays. This practice allows us to manipulate a generic YAML file, without contaminating the original file. This practice also stops the cycle of having to update and manage complex YAML templates. 

At the same time, Yot's instructions file specification provides you with documentation-as-code. This is because you have documented all the required changes to source YAML files in one place, including what is required to get an application running in your environment.

The use of JSONPath queries and templating give the tool familiar interfaces, making adoption easier, and providing for a more pleasant end-user experience.  The specification, also known as the *instructions file*, is assembled in a declarative way, where we only operate on what is clearly defined.  We take *actions* based on JSONPath query results.  We provide flexibility by allowing your instructions to be templated if needed (functionality will return in v0.3.0). See the [full documentation](docs/documentation.md), which will help get you moving along with Yot!


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
