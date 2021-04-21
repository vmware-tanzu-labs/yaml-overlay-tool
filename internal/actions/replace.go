package actions

import "gopkg.in/yaml.v3"

func Replace(original, replaceValue *yaml.Node) {
	*original = *replaceValue
}
