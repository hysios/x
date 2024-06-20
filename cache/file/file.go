package file

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hysios/x/cache"
	"github.com/hysios/x/maps"
)

type FileCache[Key, Value any] struct {
	filename  string
	immediate bool
	m         maps.Map[Key, Value]
}

type FileOpt[Key, Value any] func(*FileCache[Key, Value])

func New[Key, Value any](filename string, opts ...FileOpt[Key, Value]) cache.Cache[Key, Value] {
	c := &FileCache[Key, Value]{
		filename: filename,
	}

	for _, opt := range opts {
		opt(c)
	}

	_ = c.load()
	go c.worker()
	return c
}

// Load
func (f *FileCache[Key, Value]) Load(key Key, opts ...cache.LoadOpt) (Value, bool) {
	var z Value
	if v, ok := f.m.Load(key); ok {
		return v, true
	}
	return z, false
}

// Update
func (f *FileCache[Key, Value]) Update(key Key, val Value, opts ...cache.UpdateOpt) {
	f.m.Store(key, val)
	if f.immediate {
		f.store()
	}
}

// Clear
func (f *FileCache[Key, Value]) Clear(key Key) {
	f.m.Delete(key)
}

// load
func (f *FileCache[Key, Value]) load() error {
	file, err := os.Open(f.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var (
		key Key
		val Value
		s   = bufio.NewScanner(file)
	)

	for s.Scan() {
		line := s.Text()
		ss := strings.Split(line, "=")
		if len(ss) != 2 {
			continue
		}

		switch x := any(key).(type) {
		case string:
			x = ss[0]
			_ = x
		case int:
			i, err := strconv.Atoi(ss[0])
			if err != nil {
				continue
			}
			x = i
		case int64:
			i, err := strconv.ParseInt(ss[0], 10, 64)
			if err != nil {
				continue
			}
			x = i
		case int32:
			i, err := strconv.ParseInt(ss[0], 10, 32)
			if err != nil {
				continue
			}
			x = int32(i)
		case float64:
			i, err := strconv.ParseFloat(ss[0], 64)
			if err != nil {
				continue
			}
			x = i
		default:
			continue
		}

		switch x := any(val).(type) {
		case string:
			x = ss[0]
			_ = x
		case int:
			i, err := strconv.Atoi(ss[0])
			if err != nil {
				continue
			}
			x = i
		case int64:
			i, err := strconv.ParseInt(ss[0], 10, 64)
			if err != nil {
				continue
			}
			x = i
		case int32:
			i, err := strconv.ParseInt(ss[0], 10, 32)
			if err != nil {
				continue
			}
			x = int32(i)
		case float64:
			i, err := strconv.ParseFloat(ss[0], 64)
			if err != nil {
				continue
			}
			x = i
		case bool:
			i, err := strconv.ParseBool(ss[0])
			if err != nil {
				continue
			}
			x = i
		case time.Time:
			i, err := time.Parse(time.RFC3339, ss[0])
			if err != nil {
				continue
			}
			x = i
		default:
			continue
		}

		f.m.Store(key, val)
	}

	return nil
}

// worker
func (f *FileCache[Key, Value]) worker() error {
	var (
		ticker = time.NewTicker(10 * time.Second)
		size   = 0
	)

	for {
		select {
		case <-ticker.C:
			c := 0
			var (
				keys []Key
			)

			f.m.Range(func(key Key, val Value) bool {
				c++
				keys = append(keys, key)
				return true
			})

			if c != size {
				size = c
				sort.Slice(keys, func(i, j int) bool {
					k1 := fmt.Sprint(keys[i])
					k2 := fmt.Sprint(keys[j])
					return k1 < k2
				})

				// write to file
				out, err := os.OpenFile(f.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
				if err != nil {
					return err
				}
				defer out.Close()
				for _, key := range keys {
					val, ok := f.m.Load(key)
					if !ok {
						continue
					}
					_, _ = out.WriteString(fmt.Sprint(key))
					_, _ = out.WriteString("=")
					_, _ = out.WriteString(fmt.Sprint(val))
					_, _ = out.WriteString("\n")
				}
			}
			// do something
		}
	}
}

// store
func (f *FileCache[Key, Value]) store() error {
	var (
		keys []Key
	)

	f.m.Range(func(key Key, val Value) bool {
		keys = append(keys, key)
		return true
	})

	sort.Slice(keys, func(i, j int) bool {
		k1 := fmt.Sprint(keys[i])
		k2 := fmt.Sprint(keys[j])
		return k1 < k2
	})

	// write to file
	out, err := os.OpenFile(f.filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	for _, key := range keys {
		val, ok := f.m.Load(key)
		if !ok {
			continue
		}
		if fmt.Sprint(key) == "" {
			continue
		}

		_, _ = out.WriteString(fmt.Sprint(key))
		_, _ = out.WriteString("=")
		_, _ = out.WriteString(fmt.Sprint(val))
		_, _ = out.WriteString("\n")
	}

	return nil
}
