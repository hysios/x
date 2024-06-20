package file

import (
	"testing"
)

func TestNew(t *testing.T) {
	var cache = New[string, string]("/tmp/file_cache", WithImmediate[string, string]())
	if cache == nil {
		t.Error("New failed")
	}

	cache.Update("key", "value")
	val, ok := cache.Load("key")
	if !ok {
		t.Error("Load failed")
	}

	if val != "value" {
		t.Error("Load failed")
	}
}
