package nodes

type NodePicker interface {
	PickNode(key string) (NodeGetter, bool)
}

type NodeGetter interface {
	Get(group string, key string) ([]byte, error)
}
