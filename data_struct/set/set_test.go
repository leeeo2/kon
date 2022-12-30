package set

import "testing"

func TestSet(t *testing.T) {
	s := New(1, 2, 3, 3, 3, 4, "string")
	if s.Size() != 5 {
		t.Errorf("set.Size() error,want 5 but got %d", s.Size())
	}
	if s.Has(1) != true {
		t.Errorf("set add element failed,set has `1` but set.Has(1) return false")
	}
	s.Del(1)
	if s.Has(1) != false {
		t.Errorf("set delete elmement failed,set deleted `1` but set.Has(1) return true")
	}
	sub := New(1, 2, 3)
	if sub.IsSubOf(s) {
		t.Errorf("set.IsSubOf() error,%v is sub of %v,but return false", sub, s)
	}
	s.Clear()
	if s.Size() != 0 {
		t.Errorf("set.Clear() error,size should be 0,but got %d", s.Size())
	}
	s.Add(1, 2, 3)
	if !s.Equal(sub) {
		t.Errorf("set.Equal() error")
	}
}
