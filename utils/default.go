package utils

func Default[T comparable](val T, def ...T) T {
	var z T
	if val == z {
		for _, d := range def {
			if d != z {
				return d
			}
		}
	}
	return val
}

// DefaultSlice
func DefaultSlice[T any](val []T, def ...[]T) []T {
	if len(val) == 0 {
		for _, d := range def {
			if len(d) > 0 {
				return d
			}
		}
	}
	return val
}
