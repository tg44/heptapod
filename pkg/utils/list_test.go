package utils

import (
	"testing"
)

func TestListPrepend(t *testing.T) {
	l := &List{"a", nil}
	k := l.Prepend([]string{"b", "c", "d"})
	if l.Size() != 1 {
		t.Errorf("l.size failed")
	}
	if l.Size() != 1 {
		t.Errorf("l.size failed second time")
	}
	if k.Size() != 4 {
		t.Errorf("k.size failed")
	}
	if k.Size() != 4 {
		t.Errorf("k.size failed second time")
	}
}
