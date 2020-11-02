package interval

import (
	"fmt"
	"strings"
)

type Interval struct {
	// begin of this interval.
	Begin int
	// if IncBegin is true, this interval is inclusive of the Begin point.
	IncBegin bool

	// end of this interval.
	End int
	// if IncEnd is true, this interval is inclusive of the End point.
	IncEnd bool
}

func (i Interval) String() string {
	var b strings.Builder
	if i.IncBegin {
		b.WriteByte('[')
	} else {

		b.WriteByte('(')
	}
	fmt.Fprintf(&b, "%d", i.Begin)
	b.WriteString(", ")
	fmt.Fprintf(&b, "%d", i.End)
	if i.IncEnd {
		b.WriteByte(']')
	} else {
		b.WriteByte(')')
	}
	return b.String()
}

// Equal returns true if receiver interval is equals x interval.
func (i Interval) Equal(x Interval) bool {
	return (i.Begin == x.Begin &&
		i.End == x.End &&
		i.IncBegin == x.IncBegin &&
		i.IncEnd == x.IncEnd) || (x.IsEmpty() && i.IsEmpty())
}

// IsEmpty returns true if receiver interval has no value.
func (i Interval) IsEmpty() bool {
	if i.Begin < i.End {
		return false
	} else if i.Begin == i.End {
		return !i.IncBegin || !i.IncEnd
	}
	return true
}

// LtBeginOf returns true if receiver interval is less than begin of x interval.
func (i Interval) LtBeginOf(x Interval) bool {
	if x.IsEmpty() {
		return false
	}
	if i.IsEmpty() {
		return false
	}
	if i.End < x.Begin {
		return true
	} else if i.End == x.Begin {
		return !i.IncEnd || !x.IncBegin
	}
	return false
}

// LeEndOf returns true if receiver interval is less than or euqal to end of x interval.
func (i Interval) LeEndOf(x Interval) bool {
	if x.IsEmpty() {
		return false
	}
	if i.IsEmpty() {
		return false
	}
	if i.End < x.End {
		return true
	} else if i.End == x.End {
		return !i.IncEnd || x.IncEnd
	}
	return false
}

// Contains returns true if x interval is completely covered by receiver interval.
func (i Interval) Contains(x Interval) bool {
	if x.IsEmpty() {
		return true
	}
	if i.IsEmpty() {
		return false
	}
	if i.Begin > x.Begin {
		return false
	}
	if i.End < x.End {
		return false
	}
	if i.Begin < x.Begin && i.End > x.End {
		return true
	}
	if i.Begin == x.Begin && (i.IncBegin || !x.IncBegin) {
		return true
	}
	return i.End == x.End && (i.IncEnd || !x.IncEnd)
}

// Intersect returns the intersection of receiver interval with x interval.
func (i Interval) Intersect(x Interval) Interval {
	if x.IsEmpty() || i.IsEmpty() {
		return Interval{}
	}
	if i.Begin > x.Begin {
		x.Begin = i.Begin
		x.IncBegin = i.IncBegin
	} else if i.Begin == x.Begin && !i.IncBegin {
		x.IncBegin = false
	}
	if i.End < x.End {
		x.End = i.End
		x.IncEnd = i.IncEnd
	} else if i.End == x.End && !i.IncEnd {
		x.IncEnd = false
	}
	return maybeEmpty(x)
}

func maybeEmpty(x Interval) Interval {
	if x.IsEmpty() {
		return Interval{}
	}
	return x
}

// Move returns an interval that adds number x to begin and end of receiver interval.
func (i Interval) Move(x int) Interval {
	if i.IsEmpty() {
		return Interval{}
	}
	return Interval{
		Begin:    i.Begin + x,
		IncBegin: i.IncBegin,
		End:      i.End + x,
		IncEnd:   i.IncEnd,
	}
}

// Bisect returns two intervals, one on the before of x and one on the
// after of x, corresponding to the subtraction of x from the receiver
// interval. The returned intervals are always within the range of the
// receiver interval.
func (i Interval) Bisect(x Interval) (Interval, Interval) {
	in := i.Intersect(x)
	if in.IsEmpty() {
		if i.LtBeginOf(x) {
			return i, Interval{}
		}
		return Interval{}, i
	}
	return maybeEmpty(Interval{
			Begin:    i.Begin,
			IncBegin: i.IncBegin,
			End:      in.Begin,
			IncEnd:   !in.IncBegin,
		}), maybeEmpty(Interval{
			Begin:    in.End,
			IncBegin: !in.IncEnd,
			End:      i.End,
			IncEnd:   i.IncEnd,
		})
}

// Adjoin returns the union of two intervals, if the intervals are exactly
// adjacent, or the zero interval if they are not.
func (i Interval) Adjoin(x Interval) Interval {
	if x.IsEmpty() || i.IsEmpty() {
		return Interval{}
	}
	if i.Begin == x.End && (i.IncBegin || x.IncEnd) {
		x.End = i.End
		x.IncEnd = i.IncEnd
		return x
	}
	if i.End == x.Begin && (i.IncEnd || x.IncBegin) {
		x.Begin = i.Begin
		x.IncBegin = i.IncBegin
		return x
	}
	return Interval{}
}

// Encompass returns an interval that covers the exact extents of two intervals.
func (i Interval) Encompass(x Interval) Interval {
	if x.IsEmpty() {
		return i
	}
	if i.IsEmpty() {
		return x
	}
	if i.Begin < x.Begin {
		x.Begin = i.Begin
		x.IncBegin = i.IncBegin
	} else if i.Begin == x.Begin && i.IncBegin {
		x.IncBegin = true
	}
	if i.End > x.End {
		x.End = i.End
		x.IncEnd = i.IncEnd
	} else if i.End == x.End && i.IncEnd {
		x.IncEnd = true
	}
	return x
}
