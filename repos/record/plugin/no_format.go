package plugin

type NoFormatter interface {
	NoFmt() string
	SetNoFmt(string)
}

type Pinyin interface {
	Pinyin() string
}

func NoFormat(record interface{}) error {
	py, ok := record.(Pinyin)
	if !ok {
		return nil
	}

	py.Pinyin()

	if noFmt, ok := record.(NoFormatter); ok {
		noFmt.SetNoFmt(py.Pinyin())
	}
	return nil
}
