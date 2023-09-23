package plugin

import (
	"fmt"

	"github.com/hysios/x/maps"
)

type Desc struct {
	Name       string
	Descrption string
}

type Plugin interface {
	isPlugin()
}

type Host[T Plugin] struct {
	plugins []T
	// idxs    maps.Map[string, Plugin]
	descs maps.Map[string, PluginDesc[T]]
}

// Create a new host
func Create[T Plugin]() *Host[T] {
	return &Host[T]{
		plugins: make([]T, 0),
		descs:   maps.Map[string, PluginDesc[T]]{},
	}
}

type PluginDesc[T Plugin] struct {
	Desc   Desc
	plugin T
}

// Install
func (h *Host[T]) Install(desc Desc, plugin T) error {
	if _, ok := h.descs.Load(desc.Name); ok {
		return fmt.Errorf("plugin %s already exists", desc.Name)
	}

	h.plugins = append(h.plugins, plugin)
	h.descs.Store(desc.Name, PluginDesc[T]{
		Desc:   desc,
		plugin: plugin,
	})

	return nil
}

// Uninstall plugin
// func (h *Host[T]) Uninstall(name string) error {
// 	if _, ok := h.descs.Load(name); !ok {
// 		return fmt.Errorf("plugin %s not exists", name)
// 	}

// 	var (
// 		plugins = make([]T, 0, len(h.plugins))
// 		descs   = maps.Map[string, PluginDesc[T]]{}
// 	)

// 	for _, plugin := range h.plugins {
// 		if desc, ok := h.descs.Load(name); ok {
// 			if plugin != desc.plugin {
// 				plugins = append(plugins, plugin)
// 				descs.Store(desc.Desc.Name, desc)
// 			}
// 		}
// 	}

// 	h.plugins = plugins
// 	h.descs = descs

// 	return nil
// }

// MakeAgent
func MakeAgent[T Plugin](host *Host[T]) (T, error) {
	panic("nonimplement")
	// return nil, nil
}

// Invoke

// Range
