package kv

type Item struct {
	Key   []string
	Value interface{}
}

type ItemsProcess interface {
	ProcessOne(Item) error
	Finish() (string, error)
}
