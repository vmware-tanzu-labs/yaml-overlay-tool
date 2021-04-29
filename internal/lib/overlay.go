package lib

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

func (o *Overlay) process(f *YamlFile, i int) error {
	if ok := o.checkDocumentIndex(i); !ok {
		return nil
	}

	node := f.Nodes[i]

	ok, err := o.checkDocumentQuery(node)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	log.Debugf("%s at %s in file %s on Document %d\n", o.Action, o.Query, f.Path, i)

	yp, err := yamlpath.NewPath(o.Query)
	if err != nil {
		return fmt.Errorf("failed to parse the query path %s due to %w", o.Query, err)
	}

	results, err := yp.Find(node)
	if err != nil {
		return fmt.Errorf("failed to find results for %s, %w", o.Query, err)
	}

	if results == nil {
		log.Debugf("Call OnMissing Here")
		// o.processOnMissing(f, i)
	}

	if err := o.processActions(node, results); err != nil {
		if errors.Is(err, ErrInvalidAction) {
			return fmt.Errorf("%w in instructions file", err)
		}

		return fmt.Errorf("%w in file %s on Document %d", err, f.Path, i)
	}

	return nil
}

func (o *Overlay) processActions(node *yaml.Node, results []*yaml.Node) error {
	for ri := range results {
		b, _ := yaml.Marshal(&results[ri])
		p, _ := yaml.Marshal(o.Value)

		log.Debugf("Current: >>>\n%s\n", b)
		log.Debugf("Proposed: >>>\n%s\n", p)

		// do something with the results based on the provided overlay action
		switch o.Action {
		case "delete":
			actions.Delete(node, results[ri])
		case "replace":
			if err := actions.Replace(results[ri], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping replace for %s query result %d", err, o.Query, ri)
			}
		case "format":
			if err := actions.Format(results[ri], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping format for %s query result %d", err, o.Query, ri)
			}
		case "merge":
			if err := actions.Merge(results[ri], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping merge for %s query result %d", err, o.Query, ri)
			}
		default:
			return fmt.Errorf("%w of type '%s'", ErrInvalidAction, o.Action)
		}
	}

	return nil
}

func (o *Overlay) checkDocumentIndex(current int) bool {
	if o.DocumentIndex != nil {
		for f := range o.DocumentIndex {
			if current == o.DocumentIndex[f] {
				return true
			}
		}

		return false
	}

	return true
}

func (o *Overlay) checkDocumentQuery(node *yaml.Node) (bool, error) {
	log.Debugf("Checking Document Queries for %s", o.Query)

	if o.DocumentQuery == nil {
		log.Debugf("No Document Queries found, continuing")

		return true, nil
	}

	conditionsMet := false

	compareOptions := cmpopts.IgnoreFields(yaml.Node{}, "HeadComment", "LineComment", "FootComment", "Line", "Column", "Style")

	for i := range o.DocumentQuery {
		for ci := range o.DocumentQuery[i].Conditions {
			yp, err := yamlpath.NewPath(o.DocumentQuery[i].Conditions[ci].Key)
			if err != nil {
				return false, fmt.Errorf("failed to parse the documentQuery condition %s due to %w", o.DocumentQuery[i].Conditions[ci].Key, err)
			}

			results, err := yp.Find(node)
			if err != nil {
				return false, fmt.Errorf("failed to find results for %s, %w", o.DocumentQuery[i].Conditions[ci].Key, err)
			}

			for _, result := range results {
				conditionsMet = cmp.Equal(*result, o.DocumentQuery[i].Conditions[ci].Value, compareOptions)
				if !conditionsMet {
					break
				}
			}
		}

		if conditionsMet {
			log.Debugf("Document Query conditions were met, continuing")

			return true, nil
		}
	}

	log.Debugf("Document Query Conditions were not met, skipping")

	return false, nil
}

// func (o *Overlay) processOnMissing(f *YamlFile, i int) error {
// 	switch t := o.OnMissing.InjectPath.(type) {
// 	case []string:
// 		fmt.Println("yo")

// 	case string:
// 		fmt.Println("howdy")

// 		if ok := checkDocIndex(i, o.DocumentIndex); !ok {
// 			return nil
// 		}

// 		node := f.Nodes[i]

// 		yp, err := yamlpath.NewPath(t)
// 		if err != nil {
// 			return fmt.Errorf("failed to parse the query path %s due to %w", t, err)
// 		}

// 		results, err := yp.Find(node)
// 		if results != nil {
// 			// use replace

// } else {
// 	// check if query is dot notation (must be)
// 	m, err := regexp.MatchString(`^(\*|(\[\?\()|\.\.)`, t)
// 	if err != nil {
// 		return fmt.Errorf("error, %w", err)
// 	}
// 	if m {
// 		return fmt.Errorf("injectPath must be a fully qualified dot notation path")
// 	}

// var b map[string]interface{}
// addKey := ""

// build map
// for _, k := range strings.Split(string(t), ".") {
// 	if addKey == "" {
// 		addKey = k
// 		b[k] = ""
// 	} else {
// 		addKey += "." + k
// 		b[addKey] = ""
// 	}
// 	fmt.Println(b)
// }

// add map into original yamlNode

// use replace to insert the value as yamlNode into original yamlNode (which will preserve comments)

// 		}
// 	}

// 	return nil
// }
