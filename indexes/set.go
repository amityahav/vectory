package indexes

import (
	"fmt"
	"strings"
)

type Set[T comparable] struct {
	m map[T]struct{}
}

// NewSet returns a set containing elements.
func NewSet[T comparable](elements ...T) *Set[T] {
	ns := &Set[T]{m: make(map[T]struct{})}
	for _, e := range elements {
		ns.Add(e)
	}

	return ns

}

// Clone returns a clone of the calling set.
func (s *Set[T]) Clone() *Set[T] {
	ns := NewSet[T]()
	for v := range s.m {
		ns.Add(v)
	}
	return ns
}

// ToSlice returns a slice containing the same elements of the calling set.
func (s *Set[T]) ToSlice() []T {
	ns := make([]T, 0, s.Cardinality())
	for v := range s.m {
		ns = append(ns, v)
	}
	return ns
}

// Equals returns true if the calling set equals s2.
func (s *Set[T]) Equals(s2 *Set[T]) bool {
	return s.Cardinality() == s2.Cardinality() && s.IsSubset(s2)
}

// Cardinality returns the length of the set.
func (s *Set[T]) Cardinality() int {
	return len(s.m)
}

// Add adds elem to the calling set.
func (s *Set[T]) Add(elem T) *Set[T] {
	s.m[elem] = struct{}{}
	return s
}

// Remove removes elem from the calling set.
func (s *Set[T]) Remove(elem T) *Set[T] {
	delete(s.m, elem)
	return s
}

// Has returns true if the calling set contains elem.
func (s *Set[T]) Has(elem T) bool {
	_, ok := s.m[elem]
	return ok
}

// IsSubset returns True if the calling set is subset of s2.
func (s *Set[T]) IsSubset(s2 *Set[T]) bool {
	for e := range s.m {
		if !s2.Has(e) {
			return false
		}
	}
	return true
}

// Diff returns a new set containing all elements which are present in the calling set and not in s2.
func (s *Set[T]) Diff(s2 *Set[T]) *Set[T] {
	diff := NewSet[T]()
	for v := range s.m {
		if !s2.Has(v) {
			diff.Add(v)
		}
	}
	return diff
}

// SymDiff returns a set containing all elements in the calling set which are not in s2 and all elements in s2 which
// are not in the calling set.
func (s *Set[T]) SymDiff(s2 *Set[T]) *Set[T] {
	sd := NewSet[T]()

	for v := range s.m {
		if !s2.Has(v) {
			sd.Add(v)
		}
	}

	for v := range s2.m {
		if !s.Has(v) {
			sd.Add(v)
		}
	}

	return sd
}

// Union returns a new set containing all elements from the calling set and s2.
func (s *Set[T]) Union(s2 *Set[T]) *Set[T] {
	return Union(s, s2)
}

// Intersect returns a new set containing all elements which are present in the calling set and s2.
func (s *Set[T]) Intersect(s2 *Set[T]) *Set[T] {
	return Intersection(s, s2)
}

// All returns true if all the elements in the calling set satisfy the predicate.
func (s *Set[T]) All(pred func(T) bool) bool {
	for e := range s.m {
		if !pred(e) {
			return false
		}
	}
	return true
}

// Some returns true if at least one of the elements in the calling set satisfy the predicate.
func (s *Set[T]) Some(pred func(T) bool) bool {
	for e := range s.m {
		if !pred(e) {
			return false
		}
	}
	return true
}

// Filter returns a new set containing all elements from the calling set which satisfy the predicate.
func (s *Set[T]) Filter(pred func(T) bool) *Set[T] {
	fs := NewSet[T]()

	for e := range s.m {
		if pred(e) {
			fs.Add(e)
		}
	}
	return fs
}

// Clear deletes all elements in the calling
func (s *Set[T]) Clear() {
	s.m = make(map[T]struct{})
}

// String returns a string representation of the calling set.
func (s *Set[T]) String() string {
	elements := make([]string, 0, s.Cardinality())
	for v := range s.m {
		elements = append(elements, fmt.Sprintf("%v", v))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(elements, ", "))
}

// Union returns a set containing all the elements from the given sets.
func Union[T comparable](sets ...*Set[T]) *Set[T] {
	union := NewSet[T]()
	for _, s := range sets {
		for elem := range s.m {
			union.Add(elem)
		}
	}
	return union
}

// Intersection returns a set containing only the elements which are present in all sets.
func Intersection[T comparable](sets ...*Set[T]) *Set[T] {
	if len(sets) == 0 {
		return NewSet[T]()
	}
	intersection := sets[0]
	for i := 1; i < len(sets); i++ {
		// Minor optimization - intersection with the empty set is the empty set, we can return.
		if intersection.Cardinality() == 0 || sets[i].Cardinality() == 0 {
			break
		}
		intersection = twoSetsIntersection[T](intersection, sets[i])

	}

	return intersection
}

// Helpers

func twoSetsIntersection[T comparable](s1, s2 *Set[T]) *Set[T] {
	var (
		ss *Set[T]
		ls *Set[T]
	)
	intersection := NewSet[T]()

	if s1.Cardinality() < s2.Cardinality() {
		ss = s1
		ls = s2
	} else {
		ss = s2
		ls = s1
	}

	for v := range ss.m {
		if ls.Has(v) {
			intersection.Add(v)
		}
	}
	return intersection
}
