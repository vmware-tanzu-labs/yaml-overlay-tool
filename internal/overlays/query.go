// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package overlays

import (
	"fmt"
	"strings"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
)

type Query struct {
	yamlPath    *yamlpath.Path
	queryString string
}

type Queries []Query

func (q *Query) UnmarshalYAML(value *yaml.Node) error {
	var s string

	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("%w at line %d column %d", err, value.Line, value.Column)
	}

	p, err := yamlpath.NewPath(s)
	if err != nil {
		return fmt.Errorf("failed to parse the query path %s due to %w at line %d column %d", s, err, value.Line, value.Column)
	}

	q.yamlPath = p
	q.queryString = s

	return nil
}

func (q Query) String() string {
	return q.queryString
}

func (q *Query) Find(node *yaml.Node) []*yaml.Node {
	results, _ := q.yamlPath.Find(node)

	return results
}

func (mq *Queries) UnmarshalYAML(value *yaml.Node) error {
	var q Query

	if value.Kind == yaml.ScalarNode {
		if err := value.Decode(&q); err != nil {
			return fmt.Errorf("%w at line %d column %d", err, value.Line, value.Column)
		}

		*mq = []Query{q}

		return nil
	}

	type qq []Query

	if err := value.Decode((*qq)(mq)); err != nil {
		return fmt.Errorf("%w at line %d column %d", err, value.Line, value.Column)
	}

	return nil
}

func (mq *Queries) Paths() []string {
	s := make([]string, len(*mq))
	for i, q := range *mq {
		s[i] = q.String()
	}

	return s
}

func (mq Queries) String() string {
	return strings.Join(mq.Paths(), ",")
}

func (mq *Queries) Find(node *yaml.Node) []*yaml.Node {
	var results []*yaml.Node

	for _, q := range *mq {
		log.Debugf("searching path %s\n", q)

		result := q.Find(node)

		results = append(results, result...)
	}

	return results
}
