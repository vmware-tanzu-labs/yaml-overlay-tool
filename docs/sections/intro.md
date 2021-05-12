## Introduction to YAML Overlay Tool

`yot` is not a templating tool in the sense of a traditional text-based templating tool.  `yot` is primarily an overlayment tool, meaning we take fragments of YAML configuration and apply or inject them over the top of an existing YAML configuration.  

Each overlay operation is performed with a JSONPath `query`, a desired `value`, and an `action`.  If a JSONPath `query` returns no results, a desired value can be either ignored (default behavior) or injected (`onMissing`) and even provide a specific path or set of paths (`injectPath`) to inject the value if the initial JSONPath query was not a fully-qualified JSONPath.  
