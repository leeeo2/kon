package set

type Set struct {
	m map[interface{}]struct{}
}

func New(elems ...interface{}) *Set {
	s := &Set{}
	s.m = make(map[interface{}]struct{})
	s.Add(elems...)
	return s
}

func (s *Set) Has(elem interface{}) bool {
	_, ok := s.m[elem]
	return ok
}

func (s *Set) Add(elems ...interface{}) {
	for _, elem := range elems {
		s.m[elem] = struct{}{}
	}
}

func (s *Set) Del(elem interface{}) {
	delete(s.m, elem)
}

func (s *Set) Clear() {
	s.m = make(map[interface{}]struct{})
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) Equal(right *Set) bool {
	if s.Size() != right.Size() {
		return false
	}
	for key := range s.m {
		if !right.Has(key) {
			return false
		}
	}
	return true
}

func (s *Set) IsSubOf(other *Set) bool {
	if s.Size() > other.Size() {
		return false
	}
	for key := range s.m {
		if !other.Has(key) {
			return false
		}
	}
	return true
}
