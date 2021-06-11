// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package commands

var version = "unstable"

const (
	YotShort = "Yot (YAML Overlay Tool) is a YAML overlay tool."

	YotLong = `Yot (YAML Overlay Tool) is a YAML overlay tool which uses a templatable YAML schema to define overlay 
operations on a set of YAML documents. Yot only produces valid YAML documents on output, 
and can preserve and inject comments.`

	helpValueFile = `Path to a values file for use with templating an instructions file.
Takes multiple values files in case you would like to better organize the values. 
Each subsequent file passed with -f will be merged over the values 
from the previous. Values are applied to your instructions file when using templating.`
	/*
			helpDefaultValuesFileDeprecated = `--default-values-file argument is deprecated use --common-values instead`

			helpCommonValues = `Path to your common values file. If not set, you must
		pass a values file of common.yaml or
		common.yml within a path from the -f option.
		Takes multiple common values files in case you would
		like to separate out some of the values. After the
		first common values file, each subsequent file passed
		with -d will be merged with the values from the
		first. If a common.yaml or common.yml file is
		discovered in one of your -f paths, it will be
		merged with these values last.`

			helpValuesPath = `Values file path. May be a path to a file or directory
		containing value files ending in either .yml or .yaml.
		This option can be provided multiple times as required.
		A file named defaults.yaml or defaults.yml is required
		within the path(s) if not using the -d option, and you
		may have only 1 default value file in that scenario.
		Additional values files are merged over the defaults.yaml
		file values. Each values file is treated as a unique site
		and will render your instructions differently based on its
		values`
	*/

	helpUsageExample = "yot -i instructions.yaml -o /tmp/output"

	helpLogLevel = `Log-level to display to stdout, one of: 
	CRITICAL: {"critical", "crit", "c"},
	ERROR:    {"error", "err", "e"},
	WARNING:  {"warning", "warn", "w"},
	NOTICE:   {"notice", "note", "n"},
	INFO:     {"info", "i"},
	DEBUG:    {"debug", "d", "verbose", "v"} * used if no argument is provided
`

	helpInstructionsFile = "Path to the instructions file (required)"

	helpOutputDirectory = `Path to a directory for writing the YAML files which were operated on by Yot`
	/*	`If value files were supplied in addition to a
		defaults.yaml/.yml then the rendered templates will land
		in <output dir>/<addl value file name>.`
	*/

	helpRenderStdOut = `Output YAML files which were operated on by Yot to stdout`
	// `Templated instructions files will still be output to the --output-directory.`.

	/*	helpDumpRenderedInstructions = `If using a templated instructions file, you can dump
		the rendered instructions to stdout to allow for
		reviewing how they were rendered prior to a full run
		of yot. Equivalent to a dry-run. Exits with return
		code 0 prior to processing instructions`
	*/

	helpQuery = `define extra JSONPath query or JSONPath fully-qualified (dot-notation) path to value you would like to manipulate.
these will get added to the common overlays in conjunction with action and value`

	helpValue = "The desired value to take action with if `query` is found."

	helpAction = "The action to take when the JSONPath expression is found in the YAML document. Can be one of delete, merge, or replace."

	helpPath = `add extra yaml files outside of the instructions file to process against, theses will receive the common overlays defined,
using a path of "-" will read from stdIn`

	helpRemoveComments = `Remove all comments from the source YAML files prior to overlayment`

	helpIndentLevel = `Number of spaces to be used for indenting YAML output (min: 2, max: 9)`

	helpOutputStyle = `style to be used for rendering final documents.
multiple values can be provided to achieve the desired result, valid values are:
	NORMAL:       {"normal", "n"},
	TAGGED:       {"tagged", "tag", "t"},
	DOUBLEQUOTED: {"doubleQuoted", "doubleQuote", "double", "dq"},
	SINGLEQUOTED: {"singleQuoted", "singleQoute", "single", "sq"},
	LITERAL:      {"literal", "l"},
	FOLDED:       {"folded", "fold", "fo"},
	FLOW:         {"flow", "fl"}
`

	completionUse   = "completion [bash|zsh|fish|powershell]"
	completionShort = "Generate shell auto-completion scripts"

	completionLong = `To load completions:

Bash:

  $ source <(yot completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ yot completion bash > /etc/bash_completion.d/yot
  # macOS:
  $ yot completion bash > /usr/local/etc/bash_completion.d/yot

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ yot completion zsh > "${fpath[1]}/_yot"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ yot completion fish | source

  # To load completions for each session, execute once:
  $ yot completion fish > ~/.config/fish/completions/yot.fish

PowerShell:

  PS> yot completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> yot completion powershell > yot.ps1
  # and source this file from your PowerShell profile.
`
)
