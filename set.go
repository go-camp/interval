package interval

import (
	"sort"
	"strings"
)

// OrderedSet is a set of ordered and non-overlapping interval objects.
type OrderedSet struct {
	intervals []Interval
}

// Copy returns a copy of a ordered set that without affecting the original.
func (s OrderedSet) Copy() OrderedSet {
	return OrderedSet{append([]Interval(nil), s.intervals...)}
}

// Len returns length of intervals in this ordered set.
func (s OrderedSet) Len() int {
	return len(s.intervals)
}

// IsEmpty returns true if no intervals in this ordered set.
func (s OrderedSet) IsEmpty() bool {
	return len(s.intervals) == 0
}

func (s OrderedSet) Equal(x OrderedSet) bool {
	return equalIntervals(s.intervals, x.intervals)
}

func equalIntervals(s1, s2 []Interval) bool {
	if len(s1) != len(s2) {
		return false
	}
	for n := 0; n < len(s1); n++ {
		if !s1[n].Equal(s2[n]) {
			return false
		}
	}
	return true
}

func (s OrderedSet) String() string {
	n := len(s.intervals)
	switch n {
	case 0:
		return "{}"
	default:
		var b strings.Builder
		b.WriteByte('{')
		b.WriteString(s.intervals[0].String())
		for _, i := range s.intervals[1:] {
			b.WriteString(", ")
			b.WriteString(i.String())
		}
		b.WriteByte('}')
		return b.String()
	}
}

// Bound returns the Interval defined by the minimum and maximum values of this ordered set.
func (s OrderedSet) Bound() Interval {
	n := len(s.intervals)
	switch n {
	case 0:
		return Interval{}
	case 1:
		return s.intervals[0]
	default:
		return s.intervals[0].Encompass(s.intervals[n-1])
	}
}

//                 (low)                 (high)
//       0    1      2         3           4        5    6   7
//      === ===== ======== =========== ========= ======= == ====
//                   ================
//                        (x)
//
// searchLow returns the first index in s.intervals that is not before x.
// if not found, searchLow returns len(s.intervals).
func (s *OrderedSet) searchLow(x Interval) int {
	return sort.Search(len(s.intervals), func(i int) bool {
		return !s.intervals[i].LtBeginOf(x)
	})
}

// searchHigh returns the index of the first interval in s.intervals that is
// entirely after x.
// if not found, searchHigh returns len(s.intervals).
func (s *OrderedSet) searchHigh(x Interval) int {
	return sort.Search(len(s.intervals), func(i int) bool {
		return x.LtBeginOf(s.intervals[i])
	})
}

// Contains returns true if x interval is completely covered by this ordered set.
func (s OrderedSet) Contains(x Interval) bool {
	idx := s.searchLow(x)
	if idx == len(s.intervals) {
		return false
	}
	return s.intervals[idx].Contains(x)
}

// Intervals returns a copy of intervals in this ordered set.
func (s OrderedSet) Intervals() []Interval {
	return append([]Interval(nil), s.intervals...)
}

// Iterator returns a iterator that iterates over all the intervals both in
// this ordered set and bound.
// If iterator returns empty Interval, the iteration is over.
// If forward is true, the iteration from left to right.
func (s OrderedSet) Iterator(bound Interval, forward bool) func() Interval {
	if bound.IsEmpty() {
		return emptyIterator
	}

	low, high := s.searchLow(bound), s.searchHigh(bound)-1
	idx, stride := low, 1
	if !forward {
		idx, stride = high, -1
	}
	return func() Interval {
		if idx < low || idx > high {
			return Interval{}
		}
		x := s.intervals[idx]
		idx += stride
		return x
	}
}

func emptyIterator() Interval { return Interval{} }

func adjoinOrAppend(intervals []Interval, x Interval) []Interval {
	n := len(intervals)
	switch n {
	case 0:
		return append(intervals, x)
	default:
		n--
		ad := intervals[n].Adjoin(x)
		if ad.IsEmpty() {
			return append(intervals, x)
		}
		intervals[n] = ad
		return intervals
	}
}

// Add adds x interval to this ordered set.
// Add returns true if this ordered set changed.
func (s *OrderedSet) Add(x Interval) bool {
	if x.IsEmpty() {
		return false
	}

	low := s.searchLow(x)
	if low == len(s.intervals) {
		//                                                             (low)(high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//                                                             *=========*
		//                                                                 (x)
		s.intervals = adjoinOrAppend(s.intervals, x)
		return true
	}

	if s.intervals[low].Contains(x) {
		//                 (low)     (high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//                 =====
		return false
	}

	newIntervals := make([]Interval, 0, len(s.intervals)+1)
	newIntervals = append(newIntervals, s.intervals[:low]...)
	push := func(i Interval) {
		newIntervals = adjoinOrAppend(newIntervals, i)
	}
	if x.LtBeginOf(s.intervals[low]) {
		//   (low)(high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		// *===
		// (x)
		//
		//                         (low)(high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//                       *=*
		//                       (x)
		//
		//                         (low)(high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//                        =
		//                       (x)
		push(x)
		push(s.intervals[low])
		newIntervals = append(newIntervals, s.intervals[low+1:]...)
	} else {
		left, right := x.Bisect(s.intervals[low])
		if !left.IsEmpty() && right.IsEmpty() {
			//                 (low)     (high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= == ====
			//              *=========
			//                  (x)
			push(left)
			push(s.intervals[low])
			newIntervals = append(newIntervals, s.intervals[low+1:]...)
		} else {
			//                 (low)                 (high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= == ====
			//              *======================*
			//                        (x)
			//
			//                                                     (low)       (high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= == ====
			//                                                     *=========*
			//                                                         (x)
			//
			//     (low)                 (high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= == ====
			//    *===============
			//          (x)
			//
			//     (low)                 (high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= == ====
			//        ============
			//          (x)
			high := s.searchHigh(x)
			x = x.Encompass(s.intervals[low])
			x = x.Encompass(s.intervals[high-1])
			push(x)
			if high < len(s.intervals) {
				push(s.intervals[high])
				newIntervals = append(newIntervals, s.intervals[high+1:]...)
			}
		}
	}
	s.intervals = newIntervals
	return true
}

// Remove removes x interval from this ordered set.
// Remove returns true if this ordered set changed.
func (s *OrderedSet) Remove(x Interval) bool {
	if s.IsEmpty() || x.IsEmpty() {
		return false
	}

	low := s.searchLow(x)
	if low == len(s.intervals) {
		//                                                             (low)(high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//                                                             *=========*
		//                                                                 (x)
		return false
	}

	left, right := s.intervals[low].Bisect(x)
	if x.LeEndOf(s.intervals[low]) {
		if left.IsEmpty() {
			if right.IsEmpty() {
				//                 (low)     (high)
				//       0    1      2         3           4        5    6   7
				//      === ===== ======== =========== ========= ======= == ====
				//              *=========
				//                  (x)
				copy(s.intervals[low:], s.intervals[low+1:])
				s.intervals = s.intervals[:len(s.intervals)-1]
			} else {
				if s.intervals[low].Equal(right) {
					//   (low)(high)
					//       0    1      2         3           4        5    6   7
					//      === ===== ======== =========== ========= ======= == ====
					// *===
					// (x)
					//
					//                         (low)(high)
					//       0    1      2         3           4        5    6   7
					//      === ===== ======== =========== ========= ======= == ====
					//                       *=*
					//                       (x)
					//
					//                         (low)(high)
					//       0    1      2         3           4        5    6   7
					//      === ===== ======== =========== ========= ======= == ====
					//                        =
					//                       (x)
					return false
				}
				//   (low)(high)
				//       0    1      2         3           4        5    6   7
				//      === ===== ======== =========== ========= ======= == ====
				// *=====
				// (x)
				s.intervals[low] = right
			}
		} else if right.IsEmpty() {
			//
			//                 (low)     (high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= == ====
			//                 =======
			s.intervals[low] = left
		} else {
			//
			//                 (low)     (high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= == ====
			//                 =====
			//
			//                                                         (low)  (high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= == ====
			//                                                           ==
			//
			//                                                     (low)(high)
			//       0    1      2         3           4        5    6   7
			//      === ===== ======== =========== ========= ======= === ====
			//                                                        =
			if low+2 <= len(s.intervals) {
				s.intervals = append(s.intervals[:low+2], s.intervals[low+1:]...)
				s.intervals[low] = left
				s.intervals[low+1] = right
			} else {
				s.intervals[low] = left
				s.intervals = append(s.intervals, right)
			}
		}
		return true
	}

	high := s.searchHigh(x)
	_, right = s.intervals[high-1].Bisect(x)
	if !left.IsEmpty() {
		//                 (low)     (high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//                 =======*
		//                  (x)
		//
		//                 (low)                 (high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//                 ===========
		//                  (x)
		s.intervals[low] = left
		low++
	}
	if !right.IsEmpty() {
		//                 (low)                 (high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//               *===========
		//                  (x)
		//
		//                 (low)                 (high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//              *======================*
		//                        (x)
		//
		//                                                     (low)       (high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//                                                     *=========*
		//                                                         (x)
		//
		//     (low)                 (high)
		//       0    1      2         3           4        5    6   7
		//      === ===== ======== =========== ========= ======= == ====
		//    *===============
		//          (x)
		s.intervals[low] = right
		low++
	}
	copy(s.intervals[low:], s.intervals[high:])
	s.intervals = s.intervals[:len(s.intervals)-high+low]
	return true
}

// Union returns an ordered set containing all intervals in a or b.
func Union(a, b OrderedSet) OrderedSet {
	if a.Len() < b.Len() {
		a, b = b, a
	}
	a = a.Copy()
	it := b.Iterator(b.Bound(), true)
	for {
		x := it()
		if x.IsEmpty() {
			break
		}
		a.Add(x)
	}
	return a
}

// Intersect returns an ordered set containing all intervals of a that also belong to b.
func Intersect(a, b OrderedSet) OrderedSet {
	var intervals []Interval
	xit, yit := a.Iterator(b.Bound(), true), b.Iterator(a.Bound(), true)
	x, y := xit(), yit()
	for !x.IsEmpty() && !y.IsEmpty() {
		if x.LtBeginOf(y) {
			x = xit()
		} else if y.LtBeginOf(x) {
			y = yit()
		} else {
			in := x.Intersect(y)
			if !in.IsEmpty() {
				intervals = append(intervals, in)
				_, right := x.Bisect(y)
				if !right.IsEmpty() {
					x = right
				} else {
					x = xit()
				}
			}
		}
	}
	return OrderedSet{intervals: intervals}
}

// Subtract returns an ordered set containing all intervals in a but not in b.
func Subtract(a, b OrderedSet) OrderedSet {
	var intervals []Interval
	xit, yit := a.Iterator(a.Bound(), true), b.Iterator(a.Bound(), true)
	x, y := xit(), yit()
	for !x.IsEmpty() {
		if y.IsEmpty() {
			intervals = append(intervals, x)
			x = xit()
		} else {
			left, right := x.Bisect(y)
			if !left.IsEmpty() {
				intervals = append(intervals, left)
			}
			if right.IsEmpty() {
				x = xit()
			} else {
				x = right
				y = yit()
			}
		}
	}
	return OrderedSet{intervals: intervals}
}

// Difference returns an ordered set containing all intervals in either of a and b,
// but not in their intersection.
func Difference(a, b OrderedSet) OrderedSet {
	var intervals []Interval
	push := func(x Interval) {
		intervals = adjoinOrAppend(intervals, x)
	}
	xit, yit := a.Iterator(a.Bound(), true), b.Iterator(b.Bound(), true)
	x, y := xit(), yit()
	for {
		if x.IsEmpty() {
			if y.IsEmpty() {
				break
			}
			push(y)
			y = yit()
		} else if y.IsEmpty() {
			push(x)
			x = xit()
		} else {
			//     ======  ===   ======= ======== ==== ======
			//===   ==  *==*     =======   =========     ======
			if x.LtBeginOf(y) {
				push(x)
				x = xit()
			} else if y.LtBeginOf(x) {
				push(y)
				y = yit()
			} else {
				leftx, rightx := x.Bisect(y)
				lefty, righty := y.Bisect(x)
				if !leftx.IsEmpty() {
					push(leftx)
				}
				if rightx.IsEmpty() {
					x = xit()
				} else {
					x = rightx
				}
				if !lefty.IsEmpty() {
					push(lefty)
				}
				if righty.IsEmpty() {
					y = yit()
				} else {
					y = righty
				}
			}
		}
	}
	return OrderedSet{intervals: intervals}
}
