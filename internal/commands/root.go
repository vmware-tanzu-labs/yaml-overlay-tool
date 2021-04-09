package commands

import (
	"errors"
	"fmt"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/models"
	"gopkg.in/yaml.v3"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "yot [-h] [-d DEFAULT_VALUES_FILE] [-v VALUES_PATH] -i INSTRUCTION_FILE [-o OUTPUT_DIRECTORY] [-s] [-r] [-l LOG_FILE] [-V]",
	Short: "yot (YAML overlay tool) is a yaml overlay tool which allows for the templating of overlay instruction data with jinja2",
	Long: `yot (YAML overlay tool) is a yaml overlay tool which allows for the templating 
of overlay instruction data with jinja2, and the application of rendered 
overlays "over the top" of a yaml file. yot only produces valid yaml 
documents on output.`,
	Version: "yaml overlay tool v0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("heyyy")
		helloYot()
	},
	// PreRun check for default usage requirements and run example
	// PreRunE: func(cmd *cobra.Command, args []string) error {
	// 	return CheckRequiredFlags(cmd.Flags())
	// },
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	// TODO: this needs to be fixed so that it passed by ref and NOT to use the OS obejct.
	// 	//configPath := os.Args[2]
	// 	if err := instructFile(configPath); err != nil {
	// 		return err
	// 	}
	// 	os.Exit(0)
	// 	return nil
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP(
		"config",
		"c",
		"~/.yot",
		"config file (default is $HOME/.yot)",
	)

	// Define with Cobra
	rootCmd.PersistentFlags().StringP(
		"default-values-file",
		"d",
		"",
		`Path to your default values file. If not set, you must 
pass a values file of defaults.yaml or 
defaults.yml within a path from the -v option. 
Takes multiple default values files in case you would 
like to separate out some of the values. After the 
first default values file, each subsequent file passed 
with -d will be merged with the values from the 
first. If a defaults.yaml or defaults.yml file is 
discovered in one of your -v paths, it will be 
merged with these values last.`,
	)

	rootCmd.PersistentFlags().StringSliceP(
		"values-path",
		"v",
		[]string{},
		`Values file path. May be a path to a file or directory 
containing value files ending in either .yml or .yaml. 
This option can be provided multiple times as required. 
A file named defaults.yaml or defaults.yml is required 
within the path(s) if not using the -d option, and you
may have only 1 default value file in that scenario. 
Additional values files are merged over the defaults.yaml
file values. Each values file is treated as a unique site
and will render your instructions differently based on its
values`,
	)

	rootCmd.PersistentFlags().StringP(
		"instruction-file",
		"i",
		"",
		"Instruction file path. Defaults to ./instructions.yaml (required)",
	)

	rootCmd.MarkPersistentFlagRequired("instruction-file")

	rootCmd.PersistentFlags().StringP(
		"output-directory",
		"o",
		"./output",
		`Path to directory to write the overlayed yaml files to.
If value files were supplied in addition to a 
defaults.yaml/.yml then the rendered templates will land
in <output dir>/<addl value file name>.`,
	)

	rootCmd.PersistentFlags().StringP(
		"stdout",
		"s",
		"",
		`Render output to stdout. Templated instructions files 
will still be output to the --output-directory.`,
	)

	rootCmd.PersistentFlags().StringP(
		"dump-rendered-instructions",
		"r",
		"",
		`If using a templated instructions file, you can dump 
the rendered instructions to stdout to allow for 
reviewing how they were rendered prior to a full run 
of yot. Equivalent to a dry-run. Exits with return
code 0 prior to processing instructions`,
	)

	// Bind w/ viper
	viper.BindPFlag("default-values-file", rootCmd.PersistentFlags().Lookup("default-values-file"))
	viper.BindPFlag("values-path", rootCmd.PersistentFlags().Lookup("values-path"))
	viper.BindPFlag("instruction-file", rootCmd.PersistentFlags().Lookup("instruction-file"))
	viper.BindPFlag("output-directory", rootCmd.PersistentFlags().Lookup("output-directory"))
	viper.BindPFlag("stdout", rootCmd.PersistentFlags().Lookup("stdout"))
	viper.BindPFlag("dump-rendered-instructions", rootCmd.PersistentFlags().Lookup("dump-rendered-instructions"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".yot" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".yot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// Check for required flags identified by MarkPersistentFlagRequired.
func CheckRequiredFlags(flags *pflag.FlagSet) error {
	requiredError := false
	flagName := ""

	flags.VisitAll(func(flag *pflag.Flag) {
		requiredAnnotation := flag.Annotations[cobra.BashCompOneRequiredFlag]
		if len(requiredAnnotation) == 0 {
			return
		}

		flagRequired := requiredAnnotation[0] == "true"

		if flagRequired && !flag.Changed {
			requiredError = true
			flagName = flag.Name
		}
	})

	if requiredError {
		return errors.New("yot: error: the following arguments are required: `" + flagName + "` has not been set")
	}

	return nil
}

func helloYot() {
	var data = `
---
# the commonOverlays apply to all yamlFiles listed out in 'yamlFiles' and are processed first on each file
commonOverlays:
  - name: "add a label to certain yaml documents with refined criteria"
    query: metadata.labels
    value: {'namespace': 'tanzu-dns'}
    action: merge
    # qualifier to further refine when this overlay is applied
    documentQuery:
      # default operator behavior is 'and' and has been omitted as an example of
      ## this behavior
      # all of the 'and' operator queries must match or any one of the 'or'
      ## operator queries
    - conditions:
      - key: kind
        value: Service
      - key: metadata.labels.'app.kubernetes.io/name'
        value: external-dns
    - conditions:
      - key: metadata.name
        value: pvc-var-cache-bind
  - name: "add a common label to everything"
    query: metadata.labels
    value: {'cool_label': 'bro'}
    action: merge
yamlFiles: # what to overlay onto
  - name: "some arbitrary descriptor" # Name is Optional
    path: "examples/manifests/test.yaml"
    overlays: # if multi-doc yaml file, applies to all docs, gets applied first
    - name: "delete all annotations"
      query: metadata.annotations
      value: {}
      action: "delete"
    - name: "add in a new label"
      query: metadata.labels
      value: {'some': 'thing'}
      action: "merge"
      onMissing:
        action: "inject" # inject | ignore
    - name: "Change the apiVersion to v2alpha1"
      query: apiVersion
      value: v2alpha1
      action: replace
    # on the following 2 items, notice that the onMissing is not set
    ## these will only affect the yaml docs that have matches, otherwise ignore
    - name: "Merge in a list item"
      query: spec.ports
      value:
        - name: dns-tcp
          port: 53
          protocol: TCP
          targetPort: dns-tcp
      action: merge
    # not really a real-world example, but showing off functionality
    - name: "now replace the merged list with just the new port"
      query: spec.ports
      value:
        - name: dns-tcp
          port: 53
          protocol: TCP
          targetPort: dns-tcp
      action: replace
    # next one shouldn't do anything because no onMissing = implicit ignore
    - query: status
      value: {}
      action: "merge"
    - name: "Demo the need for an inject path"
      query: fake.key1.*
      value: {'fake': 'content1'}
      action: "merge"
      onMissing:
        action: "inject"
    # same as previous, but with an injectPath (actually does this one)
    - name: "Show same example but with an injectPath"
      query: fake.key2.*
      value: {'fake': 'content2'}
      action: "merge"
      onMissing:
        action: "inject"
        injectPath: fake2.key2
      # qualifier to only apply to the first doc in the yaml file
      documentIndex:
        - 0
    documents: # optional and only used for multi-doc yaml files
    # need to refer to them by their index
    - name: the manifest that does something
      path: 0
      overlays:
        - query: a.b.c.d
          value: {'foo': 'bar'}
          action: merge
          onMissing:
            action: "inject"
        - query: metadata.labels
          value: {'some': 'one'}
          action: merge
          onMissing:
            action: "inject"
        # demos multiple inject paths on missing
        - query: x.*
          value: {'x': 'x'}
          action: merge
          onMissing:
            action: "inject"
            injectPath:
              - x
              - y
              - z
  # demoing application of 'commonOverlays' without a 'overlays' or 'documents' key
  - name: "another file"
    path: "examples/manifests/another.yaml"
    # uncomment the following 3 lines to see this affect 2 of 3 docs in 'another.yaml' with commonOverlays {'cool_label': 'bro'}
    # documents:
    #   - path: 0
    #   - path: 2
`

	var t yaml.Node
	var ts []models.Overlay

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Printf("error: %v", err)
	}
	//fmt.Printf("--- t:\n%v\n\n", t)
	fmt.Printf("%+v\n", t.Content[0].Content[1])
	t.Content[0].Content[1].Decode(&ts)
	fmt.Printf("%+v", ts[0].Name)
}
