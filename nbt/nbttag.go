package nbt

import "fmt"

type NBTTag struct {
	Type  byte        `json:"tagType"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (n *NBTTag) BlockID() string {
	//	fmt.Printf("%+v\n", n)
	if vs, ok := n.Value.([]interface{}); ok {
		for _, t := range vs {
			if tMap, ok := t.(map[string]interface{}); ok {
				if tMap["name"] == "name" {
					return tMap["value"].(string)
				}
			}
		}
	} else {
		fmt.Println("failed to convert to NBTTag")
	}

	return ""
}
