package repos

type Pinyin interface {
	Pinyin() string
}

type DashFormat string

func (n *DashFormat) SetNoFmt(format string) {
	*n = DashFormat(format)
}

// Format
func (n *DashFormat) NoFmt() string {
	return string(*n)
}

type DashFormatter interface {
	DashFmt() string
	SetDashFmt(string)
}

type Nor interface {
	SetNo(string)
}
