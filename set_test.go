package interval

import (
	"fmt"
	"testing"
)

func parseOrderedSet(s string) *OrderedSet {
	if s == "" {
		return &OrderedSet{}
	}

	var intervals []Interval
	var begin = -1
	var incBegin bool
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '-', ' ':
			if begin != -1 {
				intervals = append(intervals, Interval{
					Begin:    begin,
					IncBegin: incBegin,
					End:      i - 1,
					IncEnd:   true,
				})
				begin = -1
			}
		case '=':
			if begin == -1 {
				begin = i
				incBegin = true
			}
		case '*':
			if begin == -1 {
				begin = i
				incBegin = false
			} else {
				intervals = append(intervals, Interval{
					Begin:    begin,
					IncBegin: incBegin,
					End:      i,
					IncEnd:   false,
				})
				begin = -1
			}
		case 'f':
			// (i,
			if begin != -1 {
				intervals = append(intervals, Interval{
					Begin:    begin,
					IncBegin: incBegin,
					End:      i - 1,
					IncEnd:   true,
				})
			}
			begin = i
			incBegin = false
		case 'p': // [i,i]
			if begin != -1 {
				intervals = append(intervals, Interval{
					Begin:    begin,
					IncBegin: incBegin,
					End:      i - 1,
					IncEnd:   true,
				})
				begin = -1
			}
			intervals = append(intervals, Interval{
				Begin:    i,
				IncBegin: true,
				End:      i,
				IncEnd:   true,
			})
		case 'e': // , e)(e,
			if begin != -1 {
				intervals = append(intervals, Interval{
					Begin:    begin,
					IncBegin: incBegin,
					End:      i,
					IncEnd:   false,
				})
			}
			begin = i
			incBegin = false
		default:
			panic(fmt.Sprintf("unsupport rune %q", s[i]))
		}
	}
	if begin != -1 {
		intervals = append(intervals, Interval{
			Begin:    begin,
			IncBegin: incBegin,
			End:      len(s) - 1,
			IncEnd:   true,
		})
	}
	return &OrderedSet{intervals: intervals}
}

func TestOrderedSet_Iterator(t *testing.T) {
	var itCases = []struct {
		s string
		b string
		w string
	}{
		{ // 1
			s: "",
			b: "",
			w: "",
		},
		{ // 2
			s: "",
			b: "===",
			w: "",
		},
		{ // 3
			s: "===",
			b: "",
			w: "",
		},
		{ // 4
			s: "===",
			b: "===",
			w: "===",
		},
		{ // 5
			s: "===",
			b: " =",
			w: "===",
		},
		{ // 6
			s: " ===",
			b: "=====",
			w: " ===",
		},
		{ // 7
			s: " ===   === ====      ===     ===  =",
			b: "      =====*",
			w: "       ===",
		},
		{ // 8
			s: " ===   === ====      ===     ===  =",
			b: "      ======",
			w: "       === ====",
		},
		{ // 9
			s: " ===   === ====      ===     ===  =",
			b: "========",
			w: " ===   ===",
		},
		{ // 10
			s: " ===   === ====      ===     ===  =",
			b: "                              ====*",
			w: "                             ===",
		},
		{ // 11
			s: " ===   === ====      ===     ===  =",
			b: "                              ========",
			w: "                             ===  =",
		},
	}
	for n, tc := range itCases {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			s := parseOrderedSet(tc.s)
			b := parseInterval(tc.b)
			w := parseOrderedSet(tc.w)
			it := s.Iterator(b, true)
			var intervals []Interval
			for {
				i := it()
				if i.IsEmpty() {
					break
				}
				intervals = append(intervals, i)
			}
			if !equalIntervals(intervals, w.Intervals()) {
				t.Errorf("want %s.Iterator(%s, true) = %s get get %s", s, b, w.Intervals(), intervals)
			}
		})
	}
}

func TestOrderedSet_Add(t *testing.T) {
	var addCases = []struct {
		s string
		a string
		w string
		c bool
	}{
		{ // 0
			s: "",
			a: "*=====*",
			w: "*=====*",
			c: true,
		},
		{ // 1
			s: "*=====*",
			a: "",
			w: "*=====*",
			c: false,
		},
		{ // 2
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "                                                             *=========*",
			w: "      === ===== ======== =========== ========= ======= == =============*",
			c: true,
		},
		{ // 3
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "                 =====",
			w: "      === ===== ======== =========== ========= ======= == ====",
			c: false,
		},
		{ // 4
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: " *===",
			w: " *=== === ===== ======== =========== ========= ======= == ====",
			c: true,
		},
		{ // 5
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "                       *=*",
			w: "      === ===== ==================== ========= ======= == ====",
			c: true,
		},
		{ // 6
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "                        =",
			w: "      === ===== ========p=========== ========= ======= == ====",
			c: true,
		},
		{ // 7
			s: "      === ===== ========e=========== ========= ======= == ====",
			a: "                        =",
			w: "      === ===== ==================== ========= ======= == ====",
			c: true,
		},
		{ // 8
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "              *=========",
			w: "      === ============== =========== ========= ======= == ====",
			c: true,
		},
		{ // 9
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "              *======================*",
			w: "      === ==================================== ======= == ====",
			c: true,
		},
		{ // 10
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "                                                     *=========*",
			w: "      === ===== ======== =========== ========= ================*",
			c: true,
		},
		{ // 11
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "    *===============",
			w: "    *=================== =========== ========= ======= == ====",
			c: true,
		},
		{ // 12
			s: "      === ===== ======== =========== ========= ======= == ====",
			a: "        ============",
			w: "      ================== =========== ========= ======= == ====",
			c: true,
		},
	}

	for n, tc := range addCases {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			s := parseOrderedSet(tc.s)
			os := s.Copy()
			a := parseInterval(tc.a)
			c := s.Add(a)
			if c != tc.c {
				t.Errorf("want changed is %v but get %v", tc.c, c)
			}
			w := parseOrderedSet(tc.w)
			if !s.Equal(w) {
				t.Errorf("want %s.Add(%s) = %s but get %s", os, a, w, s)
			}
		})
	}
}

func TestOrderedSet_Remove(t *testing.T) {
	var removeCases = []struct {
		s string
		r string
		w string
		c bool
	}{
		{ // 0
			s: "",
			r: "*====*",
			w: "",
			c: false,
		},
		{ // 1
			s: "*====*",
			r: "",
			w: "*====*",
			c: false,
		},
		{ // 2
			s: "*====*",
			r: "*====*",
			w: "",
			c: true,
		},
		{ // 3
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                                                             *=========*",
			w: "      === ===== ======== =========== ========= ======= == ====",
			c: false,
		},
		{ // 4
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                 =====",
			w: "      === ===== =*   *== =========== ========= ======= == ====",
			c: true,
		},
		{ // 5
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                                                           ==",
			w: "      === ===== ======== =========== ========= ======= == =**=",
			c: true,
		},
		{ // 6
			s: "      === ===== ======== =========== ========= ======= === ====",
			r: "                                                        =",
			w: "      === ===== ======== =========== ========= ======= =e= ====",
			c: true,
		},
		{ // 7
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "              *=========",
			w: "      === =====          =========== ========= ======= == ====",
			c: true,
		},
		{ // 8
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: " *===",
			w: "      === ===== ======== =========== ========= ======= == ====",
			c: false,
		},
		{ // 9
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                       *=*",
			w: "      === ===== ======== =========== ========= ======= == ====",
			c: false,
		},
		{ // 10
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                        =",
			w: "      === ===== ======== =========== ========= ======= == ====",
			c: false,
		},
		{ // 11
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                 =======*",
			w: "      === ===== =*       =========== ========= ======= == ====",
			c: true,
		},
		{ // 12
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                 ===========",
			w: "      === ===== =*         *======== ========= ======= == ====",
			c: true,
		},
		{ // 13
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "               *===========",
			w: "      === =====           *========= ========= ======= == ====",
			c: true,
		},
		{ // 14
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "              *======================*",
			w: "      === =====                      ========= ======= == ====",
			c: true,
		},
		{ // 15
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                                                     *=========*",
			w: "      === ===== ======== =========== ========= =======",
			c: true,
		},
		{ // 16
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "    *===============",
			w: "                   *==== =========== ========= ======= == ====",
			c: true,
		},
		{ // 17
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: " *=====",
			w: "      *== ===== ======== =========== ========= ======= == ====",
			c: true,
		},
		{ // 18
			s: "      === ===== ======== =========== ========= ======= == ====",
			r: "                 =======",
			w: "      === ===== =*       =========== ========= ======= == ====",
			c: true,
		},
	}

	for n, tc := range removeCases {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			s := parseOrderedSet(tc.s)
			os := s.Copy()
			r := parseInterval(tc.r)
			c := s.Remove(r)
			if c != tc.c {
				t.Errorf("want changed is %v but get %v", tc.c, c)
			}
			w := parseOrderedSet(tc.w)
			if !s.Equal(w) {
				t.Errorf("want %s.Remove(%s) = %s but get %s", os, r, w, s)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	var unionCases = []struct {
		a string
		b string

		w string
	}{
		{ // 0
			a: "",
			b: "",

			w: "",
		},
		{ // 1
			a: "===",
			b: "",

			w: "===",
		},
		{ //2
			a: "",
			b: "==",

			w: "==",
		},
		{ // 3
			a: "  ===",
			b: "=",

			w: "= ===",
		},
		{ // 4
			a: "  ===",
			b: "=*",

			w: "=*===",
		},
		{ // 5
			a: "  ===",
			b: "==*",

			w: "=====",
		},
		{ // 6
			a: "  *==",
			b: "==*",

			w: "==e==",
		},
		{ // 7
			a: "  ===",
			b: "===",

			w: "=====",
		},
		{ // 8
			a: "  ===",
			b: "====*",

			w: "=====",
		},
		{ // 9
			a: "  ===",
			b: "=====",

			w: "=====",
		},
		{ // 10
			a: "  ===",
			b: "=====*",

			w: "=====*",
		},
		{ // 11
			a: "  ===",
			b: "======",

			w: "======",
		},
		{ // 12
			a: "  === *====  ==",
			b: "   ==========* ",

			w: "  =============",
		},
		{ // 13
			a: "  === *====  =====    =====",
			b: "   ==========*   *===",

			w: "  =================== =====",
		},
		{ // 14
			a: "  === *====  =====    =====",
			b: "   ==========*   ====",

			w: "  =================== =====",
		},
		{ // 15
			a: "  === *====  =====    =====",
			b: "   ==========*    *==         **",

			w: "  ================f== =====   **",
		},
		{ // 16
			a: "  === *====  =====*    =====",
			b: "   ==========*     ==         **",

			w: "  ================*==  =====  **",
		},
		{ // 17
			a: "  === *====  =====*    =====",
			b: "   ==========*    *==         **",

			w: "  ================e==  =====  **",
		},
		{ // 18
			a: "  === *====  =====*    =====",
			b: "=  ==========*    *==         **",

			w: "= ================e==  =====  **",
		},
		{ // 19
			a: "=e=",
			b: " p",

			w: "===",
		},
	}
	for n, tc := range unionCases {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			a := parseOrderedSet(tc.a)
			b := parseOrderedSet(tc.b)
			w := parseOrderedSet(tc.w)
			s := Union(a, b)
			if !s.Equal(w) {
				t.Errorf("want Union(%s, %s) = %s but get %s", a, b, w, s)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	var intersectCases = []struct {
		a string
		b string

		w string
	}{
		{ // 0
			a: "",
			b: "",
			w: "",
		},
		{ // 1
			a: "===",
			b: "",
			w: "",
		},
		{ // 2
			a: "===",
			b: "===",
			w: "===",
		},
		{ // 3
			a: "===",
			b: " =",
			w: " =",
		},
		{ // 4
			a: "=e=",
			b: " p",
			w: "",
		},
		{ // 5
			a: "=e=",
			b: "p",
			w: "p",
		},
		{ // 6
			a: "  ===  ====   =====    ======",
			b: "=     =     ==     ===        ====",
			w: "",
		},
		{ // 7
			a: "  ===  ====   =====    ======",
			b: "==* *==*  ====*   *=====    ====",
			w: "          =            =    =",
		},
		{ // 8
			a: "  ===  ====   =====    ======",
			b: "===* ====* *=====*   *=====**====",
			w: "  =*   ==*    ===*     ====*",
		},
	}
	for n, tc := range intersectCases {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			a := parseOrderedSet(tc.a)
			b := parseOrderedSet(tc.b)
			w := parseOrderedSet(tc.w)
			s := Intersect(a, b)
			if !s.Equal(w) {
				t.Errorf("want Intersect(%s, %s) = %s but get %s", a, b, w, s)
			}
			s = Intersect(b, a)
			if !s.Equal(w) {
				t.Errorf("want Intersect(%s, %s) = %s but get %s", b, a, w, s)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	var subtractCases = []struct {
		a string
		b string

		w string
	}{
		{ // 0
			a: "",
			b: "",

			w: "",
		},
		{ // 1
			a: "==",
			b: "",

			w: "==",
		},
		{ // 2
			a: "",
			b: "==",

			w: "",
		},
		{ // 3
			a: "=e=",
			b: " p",

			w: "=e=",
		},
		{ // 4
			a: "   ===        ==== ",
			b: "==     ===         ==",

			w: "   ===        ====",
		},
		{ // 5
			a: "   ===        ==== ",
			b: "==*   *======*    *==",

			w: "   ===        ====",
		},
		{ // 6
			a: "   ===        ==== ",
			b: "===* *========*  *===",

			w: "   ===        ====",
		},
		{ // 7
			a: "   ===        ==== ",
			b: "====**==========**====",

			w: "    ==          ==",
		},
	}
	for n, tc := range subtractCases {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			a := parseOrderedSet(tc.a)
			b := parseOrderedSet(tc.b)
			w := parseOrderedSet(tc.w)
			s := Subtract(a, b)
			if !s.Equal(w) {
				t.Errorf("want Subtract(%s, %s) = %s but get %s", a, b, w, s)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	var differenceCases = []struct {
		a string
		b string

		w string
	}{
		{ // 0
			a: "",
			b: "",

			w: "",
		},
		{ // 1
			a: "",
			b: "==",

			w: "==",
		},
		{ // 2
			a: "=e=",
			b: " p",

			w: "===",
		},
		{ // 3
			a: "   ===        ==== ",
			b: "==     ===         ==",

			w: "== === ===    ==== ==",
		},
		{ // 4
			a: "   ===        ==== ",
			b: "==*   *======*    *==",

			w: "==*===f======*====f==",
		},
		{ // 5
			a: "   ===        ==== ",
			b: "===* *========*  *===",

			w: "=====================",
		},
		{ // 6
			a: "   ===        ==== ",
			b: "====**==========**====",

			w: "===*==========* ======",
		},
	}
	for n, tc := range differenceCases {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			a := parseOrderedSet(tc.a)
			b := parseOrderedSet(tc.b)
			w := parseOrderedSet(tc.w)
			s := Difference(a, b)
			if !s.Equal(w) {
				t.Errorf("want Difference(%s, %s) = %s but get %s", a, b, w, s)
			}
			s = Difference(b, a)
			if !s.Equal(w) {
				t.Errorf("want Difference(%s, %s) = %s but get %s", b, a, w, s)
			}
		})
	}
}
