package lib

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/path"
	"gopkg.in/yaml.v3"
)

var (
	ErrOnMissingNoInjectAction = errors.New("no matches and no onMissing.action of 'inject'")
	ErrOnMissingNoInjectPath   = errors.New("no matches and no onMissing.injectPath")
	ErrOnMissingInvalidType    = errors.New("invalid type for onMissing.injectPath")
)

func (o *Overlay) process(f *YamlFile, docIndex int) error {
	if ok := o.checkDocumentIndex(docIndex); !ok {
		return nil
	}

	node := f.Nodes[docIndex]

	ok, err := o.checkDocumentQuery(node)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	log.Debugf("%s at %s in file %s on Document %d\n", o.Action, o.Query, f.Path, docIndex)

	results, err := searchPath(o.Query, node)
	if err != nil {
		return err
	}

	if results == nil {
		if err := o.onMissing(f, docIndex); err != nil {
			return err
		}
	}

	if err := o.processActions(node, results); err != nil {
		if errors.Is(err, ErrInvalidAction) {
			return fmt.Errorf("%w in instructions file", err)
		}

		return fmt.Errorf("%w in file %s on Document %d", err, f.Path, docIndex)
	}

	return nil
}

func (o *Overlay) processActions(node *yaml.Node, results []*yaml.Node) error {
	for ri := range results {
		b, _ := yaml.Marshal(&results[ri])
		p, _ := yaml.Marshal(o.Value)

		log.Debugf("Current: >>>\n%s\n", b)
		log.Debugf("Proposed: >>>\n%s\n", p)

		if err := o.doAction(node, results[ri]); err != nil {
			return fmt.Errorf("%w for %s query result %d", err, o.Query, ri)
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

func (o *Overlay) onMissing(f *YamlFile, docIndex int) error {
	// check if the query has a match
	// if no match then we require an inject path
	// we need to then check if each inject path is valid (does it exist)
	// if we had an inject path(s) then we inject the value to those locations
	// if we didn't have an inject path we have an implicit onMissing: ignore and we put out a warning if not stdout option to terminal
	switch o.OnMissing.Action {
	case "ignore", "":
		log.Debugf("ignoring %s at %s in file %s on Document %d due to %s\n", o.Action, o.Query, f.Path, docIndex, ErrOnMissingNoInjectAction)

		return nil
	case "inject":
		if o.OnMissing.InjectPath == nil {
			log.Debugf("ignoring %s at %s in file %s on Document %d due to %s\n", o.Action, o.Query, f.Path, docIndex, ErrOnMissingNoInjectPath)

			return nil
		}

		injectPaths, err := o.OnMissing.getInjectPaths()
		if err != nil {
			return err
		}

		for _, injectPath := range injectPaths {
			if err := o.doInjectPath(injectPath, f.Nodes[docIndex]); err != nil {
				return fmt.Errorf("%w in file %s on Document %d", err, f.Path, docIndex)
			}
		}
	default:
		return fmt.Errorf("%w for onMissing of type '%s'", ErrInvalidAction, o.Action)
	}

	return nil
}

func searchPath(q string, node *yaml.Node) ([]*yaml.Node, error) {
	yp, err := yamlpath.NewPath(q)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the query path %s due to %w", q, err)
	}

	results, err := yp.Find(node)
	if err != nil {
		return nil, fmt.Errorf("failed to find results for %s, %w", q, err)
	}

	return results, nil
}

func (om *OnMissing) getInjectPaths() ([]string, error) {
	var injectPaths []string
	switch injectPath := om.InjectPath.(type) {
	case string:
		injectPaths = append(injectPaths, injectPath)
	case []interface{}:
		for _, v := range injectPath {
			injectPaths = append(injectPaths, v.(string))
		}
	default:
		return nil, fmt.Errorf("%w: %T", ErrOnMissingInvalidType, injectPath)
	}

	return injectPaths, nil
}

func (o *Overlay) doInjectPath(ip string, node *yaml.Node) error {
	y, err := path.Build(ip)
	if err != nil {
		return fmt.Errorf("failed to build inject path %s, %w", ip, err)
	}

	err = actions.Merge(node, y)
	if err != nil {
		return fmt.Errorf("failed to merge injectpath scaffolding %s with document, %w", ip, err)
	}

	results, err := searchPath(ip, node)
	if err != nil {
		return fmt.Errorf("%w, on injectPath %s", err, ip)
	}

	for i := range results {
		if err := actions.Replace(results[i], &o.Value); err != nil {
			if errors.Is(err, ErrInvalidAction) {
				return fmt.Errorf("%w in instructions file", err)
			}

			return fmt.Errorf("%w for onMissing.InjectPath", err)
		}
	}

	return nil
}

func (o *Overlay) doAction(root, node *yaml.Node) error {
	switch o.Action {
	case "delete":
		actions.Delete(root, node)
	case "replace":
		if err := actions.Replace(node, &o.Value); err != nil {
			return fmt.Errorf("%w, skipping replace", err)
		}
	case "format":
		if err := actions.Format(node, &o.Value); err != nil {
			return fmt.Errorf("%w, skipping format", err)
		}
	case "merge":
		if err := actions.Merge(node, &o.Value); err != nil {
			return fmt.Errorf("%w, skipping merge", err)
		}
	default:
		return fmt.Errorf("%w of type '%s'", ErrInvalidAction, o.Action)
	}

	return nil
}
