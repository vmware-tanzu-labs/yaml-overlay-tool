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

func (q *Query) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string

	if err := unmarshal(&s); err != nil {
		return err
	}

	p, err := yamlpath.NewPath(s)
	if err != nil {
		return fmt.Errorf("failed to parse the query path %s due to %w", s, err)
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

func (mq *Queries) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var q Query
	if err := unmarshal(&q); err == nil {
		*mq = []Query{q}

		return nil
	} else if !strings.Contains(err.Error(), "cannot unmarshal") {
		return err
	}

	type qq []Query

	return unmarshal((*qq)(mq))
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
