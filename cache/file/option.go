package file

func WithImmediate[Key, Value any]() FileOpt[Key, Value] {
	return func(f *FileCache[Key, Value]) {
		f.immediate = true
	}
}
