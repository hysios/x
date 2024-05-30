package utils

// Default
func Default[T comparable](val T, def T) T {
	var z T
	if val == z {
		return def
	}
	return val
}
