package repos

import "reflect"

type Plugins struct {
	creates RecordPlugins
	fields  FieldPlugins
}

type RecordPlugin interface {
	Update(record interface{}, env PluginEnv) error
}

type FieldPlugin interface {
	Tag(reflect.StructTag) (interface{}, bool)
}

type (
	RecordPlugins []RecordPlugin
	FieldPlugins  []FieldPlugin
)

type PluginEnv struct {
	Tags map[string]reflect.StructTag
	Vars map[string]interface{}
}

func (p *Plugins) Create(record interface{}) error {
	var env = PluginEnv{
		Tags: make(map[string]reflect.StructTag),
		Vars: make(map[string]interface{}),
	}

	p.loadEnv(record, &env)

	for _, create := range p.creates {
		if err := create.Update(record, env); err != nil {
			return err
		}
	}
	return nil
}

// loadEnv
func (p *Plugins) loadEnv(record interface{}, env *PluginEnv) {
	t := reflect.TypeOf(record).Elem()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		env.Tags[field.Name] = field.Tag

		for _, fieldPlugin := range p.fields {
			val, ok := fieldPlugin.Tag(field.Tag)
			if ok {
				env.Vars[field.Name] = val
			}
		}
	}
}

func (p *RecordPlugins) Add(plugin RecordPlugin) {
	(*p) = append(*p, plugin)
}

func (p *FieldPlugins) Add(plugin FieldPlugin) {
	(*p) = append(*p, plugin)
}

func (p *Plugins) AddField(plugin FieldPlugin) {
	p.fields.Add(plugin)
}

func GetPinyin(record interface{}) string {
	if pinyin, ok := record.(Pinyin); ok {
		return pinyin.Pinyin()
	}
	return ""
}
