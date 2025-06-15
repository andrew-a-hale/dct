package jsonschema2sql

type TreeNode struct {
	Id       string
	Type     string
	Children map[string]*TreeNode
	IsLeaf   bool
	IsArray  bool
}

func NewTreeNode(id string) *TreeNode {
	return &TreeNode{
		Id:       id,
		Children: make(map[string]*TreeNode),
	}
}

func BuildTree(fields []Field) *TreeNode {
	root := NewTreeNode("root")

	for _, field := range fields {
		current := root

		for _, segment := range field.Path {
			if current.Children[segment] == nil {
				current.Children[segment] = NewTreeNode(segment)
			}
			current = current.Children[segment]
		}

		current.Type = field.Type
		current.IsLeaf = true
		current.IsArray = field.Array
	}

	return root
}
