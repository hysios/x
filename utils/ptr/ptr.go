package ptr

func Type[T any](p *T) T {
	var z T
	if p == nil {
		return z
	}

	return *p
}
