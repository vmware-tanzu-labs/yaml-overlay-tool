// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rwtodd/Go.Sed/sed"
	"gopkg.in/yaml.v3"
)

func format(formatStr string, v ...interface{}) string {
	v = sanitizeValues(v...)

	return doFormat(formatStr, v...)
}

func sanitizeNode(n ...*yaml.Node) {
	for _, nv := range n {
		if nv == nil {
			continue
		}

		switch nv.Kind {
		case yaml.DocumentNode, yaml.MappingNode, yaml.SequenceNode, yaml.AliasNode:
			sanitizeNode(nv.Content...)
		case yaml.ScalarNode:
		}

		values := []*string{&nv.Value, &nv.HeadComment, &nv.LineComment, &nv.FootComment}
		for _, v := range values {
			*v = format(*v, "", "", "", "", "")
		}
	}
}

func sanitizeValues(v ...interface{}) []interface{} {
	for _, t := range []string{
		ValueMarker,
		LineCommentMarker,
		HeadCommentMarker,
		FootCommentMarker,
		KeyMarker,
	} {
		for j := range v {
			v[j] = strings.ReplaceAll(v[j].(string), t, "")
		}
	}

	return v
}

func doFormat(s string, values ...interface{}) string {
	markers := []string{
		ValueMarker,
		LineCommentMarker,
		HeadCommentMarker,
		FootCommentMarker,
		KeyMarker,
	}

	valueMap := map[string]string{}

	for i, v := range markers {
		if sv, ok := values[i].(string); ok {
			valueMap[v] = sv
		}
	}

	re := regexp.MustCompile(`(?P<marker>%[vklfh])(?P<format>{(?P<command>.*?)})?`)

	matches := re.FindAllStringSubmatch(s, -1)

	for _, match := range matches {
		result := make(map[string]string)

		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		if result["marker"] != "" {
			marker := valueMap[result["marker"]]

			if result["command"] != "" {
				var err error

				marker, err = processSedCommand(result["command"], valueMap[result["marker"]])
				if err != nil {
					log.Warningf("Skipping additional format on [%s] due to invalid sed exp [%s], %s", s, result["command"], err)
				}
			}

			s = strings.Replace(s, match[0], marker, 1)
		}
	}

	return s
}

func processSedCommand(command, value string) (string, error) {
	engine, err := sed.New(strings.NewReader(command))
	if err != nil {
		return value, fmt.Errorf("sed error: %w", err)
	}

	value, err = engine.RunString(value)
	if err != nil {
		return value, fmt.Errorf("sed error: %w", err)
	}

	value = strings.TrimSuffix(value, "\n")
	value = strings.TrimPrefix(value, "\n")

	return value, nil
}
