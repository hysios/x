package plugin

import "reflect"

type noformat struct{}

func (n *noformat) Tag(s reflect.StructTag) (interface{}, bool) {
	v, ok := s.Lookup("noformat")
	if !ok {
		return nil, false
	}
	_ = v
	// n.build(v)
	panic("nonimplement")
}
