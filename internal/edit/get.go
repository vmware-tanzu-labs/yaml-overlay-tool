package edit

import "gopkg.in/yaml.v3"

func IterateNode(node *yaml.Node, identifier string) *yaml.Node {
	returnNode := false

	for _, n := range node.Content {
		if n.Value == identifier {
			returnNode = true
			continue
		}

		if returnNode {
			return n
		}

		if len(n.Content) > 0 {
			acNode := IterateNode(n, identifier)
			if acNode != nil {
				return acNode
			}
		}
	}

	return nil
}
