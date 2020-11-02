package interval

import (
	"fmt"
	"strings"
	"testing"
)

var testIntervals = []struct {
	i string
	x string

	// i.Before(x)
	a bool
	// x.Before(i)
	b bool

	// i.Cover(x)
	c bool
	// x.Cover(i)
	d bool

	// i.Intersect(x)
	e string
	// x.Intersect(i)
	// == e
	// f string

	// i.Bisect(x)
	g, h string
	// x.Bisect(i)
	j, k string

	// i.Adjoin(x)
	l string
	// x.Adjoin(i)
	// l == m
	// m string

	// i.Encompass(x)
	o string
	// x.Encompass(i)
	// == o
	// p string
}{
	{ // 0
		i: "=====",
		x: "-------=========",
		a: true,
		b: false,
		c: false,
		d: false,
		e: "",

		g: "=====",
		h: "",
		j: "",
		k: "-------=========",

		l: "",

		o: "================",
	},
	{ // 1
		i: "=====",
		x: "------=========",
		a: true,
		b: false,
		c: false,
		d: false,
		e: "",

		g: "=====",
		h: "",
		j: "",
		k: "------=========",

		l: "",

		o: "===============",
	},
	{ // 2
		i: "=====",
		x: "-----*=========",
		a: true,
		b: false,
		c: false,
		d: false,
		e: "",

		g: "=====",
		h: "",
		j: "",
		k: "-----*=========",

		l: "",

		o: "===============",
	},
	{ // 3
		i: "=====",
		x: "-----=========",
		a: true,
		b: false,
		c: false,
		d: false,
		e: "",

		g: "=====",
		h: "",
		j: "",
		k: "-----=========",

		l: "",

		o: "==============",
	},
	{ // 4
		i: "=====",
		x: "----*=========",
		a: true,
		b: false,
		c: false,
		d: false,
		e: "",

		g: "=====",
		h: "",
		j: "",
		k: "----*=========",

		l: "==============",

		o: "==============",
	},
	{ // 5
		i: "=====",
		x: "----=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "----=",

		g: "====*",
		h: "",
		j: "",
		k: "----*========",

		l: "=============",

		o: "=============",
	},
	{ // 6
		i: "=====",
		x: "--=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "--===",

		g: "==*",
		h: "",
		j: "",
		k: "----*======",

		l: "",

		o: "===========",
	},
	{ // 7
		i: "=====",
		x: "-=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "-====",

		g: "=*",
		h: "",
		j: "",
		k: "----*=====",

		l: "",

		o: "==========",
	},
	{ // 8
		i: "=====",
		x: "*=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "*====",

		g: "=",
		h: "",
		j: "",
		k: "----*=====",

		l: "",

		o: "==========",
	},
	{ // 9
		i: "=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: true,
		e: "=====",

		g: "",
		h: "",
		j: "",
		k: "----*====",

		l: "",

		o: "=========",
	},
	{ // 10
		i: "*=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: true,
		e: "*=====",

		g: "",
		h: "",
		j: "=",
		k: "-----*===",

		l: "",

		o: "=========",
	},
	{ // 11
		i: "-=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: true,
		e: "-=====",

		g: "",
		h: "",
		j: "=*",
		k: "-----*===",

		l: "",

		o: "=========",
	},
	{ // 12
		i: "--=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: true,
		e: "--=====",

		g: "",
		h: "",
		j: "==*",
		k: "------*==",

		l: "",

		o: "=========",
	},
	{ // 13
		i: "---=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: true,
		e: "---=====",

		g: "",
		h: "",
		j: "===*",
		k: "-------*=",

		l: "",

		o: "=========",
	},
	{ // 14
		i: "---=====*",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: true,
		e: "---=====*",

		g: "",
		h: "",
		j: "===*",
		k: "--------=",

		l: "",

		o: "=========",
	},
	{ // 15
		i: "----=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: true,
		e: "----=====",

		g: "",
		h: "",
		j: "====*",
		k: "",

		l: "",

		o: "=========",
	},
	{ // 16
		i: "----=====*",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "----=====",

		g: "",
		h: "--------**",
		j: "====*",
		k: "",

		l: "",

		o: "=========*",
	},
	{ // 17
		i: "-----=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "-----====",

		g: "",
		h: "--------*=",
		j: "=====*",
		k: "",

		l: "",

		o: "==========",
	},
	{ // 18
		i: "------=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "------===",

		g: "",
		h: "--------*==",
		j: "======*",
		k: "",

		l: "",

		o: "===========",
	},
	{ // 19
		i: "-------=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "-------==",

		g: "",
		h: "--------*===",
		j: "=======*",
		k: "",

		l: "",

		o: "============",
	},
	{ // 20
		i: "--------=====",
		x: "=========",
		a: false,
		b: false,
		c: false,
		d: false,
		e: "--------=",

		g: "",
		h: "--------*====",
		j: "========*",
		k: "",

		l: "=============",

		o: "=============",
	},
	{ // 21
		i: "--------*=====",
		x: "=========",
		a: false,
		b: true,
		c: false,
		d: false,
		e: "",

		g: "",
		h: "--------*=====",
		j: "=========",
		k: "",

		l: "==============",

		o: "==============",
	},
	{ // 22
		i: "---------=====",
		x: "=========",
		a: false,
		b: true,
		c: false,
		d: false,
		e: "",

		g: "",
		h: "---------=====",
		j: "=========",
		k: "",

		l: "",

		o: "==============",
	},
	{ // 23
		i: "---------*=====",
		x: "=========",
		a: false,
		b: true,
		c: false,
		d: false,
		e: "",

		g: "",
		h: "---------*=====",
		j: "=========",
		k: "",

		l: "",

		o: "===============",
	},
	{ // 24
		i: "----------=====",
		x: "=========",
		a: false,
		b: true,
		c: false,
		d: false,
		e: "",

		g: "",
		h: "----------=====",
		j: "=========",
		k: "",

		l: "",

		o: "===============",
	},
	{ // 25
		i: "-----------=====",
		x: "=========",
		a: false,
		b: true,
		c: false,
		d: false,
		e: "",

		g: "",
		h: "-----------=====",
		j: "=========",
		k: "",

		l: "",

		o: "================",
	},
}

func parseInterval(s string) Interval {
	if s == "" {
		return Interval{}
	}
	begin := strings.IndexAny(s, "*=")
	end := strings.LastIndexAny(s, "*=")
	return Interval{
		Begin:    begin,
		IncBegin: s[begin] == '=',
		End:      end,
		IncEnd:   s[end] == '=',
	}
}

func TestInterval(t *testing.T) {
	for n, tc := range testIntervals {
		t.Run(fmt.Sprint(n), func(t *testing.T) {
			i := parseInterval(tc.i)
			x := parseInterval(tc.x)

			a, b := i.LtBeginOf(x), x.LtBeginOf(i)
			if a != tc.a {
				t.Errorf("want %s.LtBeginOf(%s) = %v but get %v", i, x, tc.a, a)
			}
			if b != tc.b {
				t.Errorf("want %s.LtBeginOf(%s) = %v but get %v", x, i, tc.b, b)
			}

			c, d := i.Contains(x), x.Contains(i)
			if c != tc.c {
				t.Errorf("want %s.Cover(%s) = %v but get %v", i, x, tc.c, c)
			}
			if d != tc.d {
				t.Errorf("want %s.Cover(%s) = %v but get %v", x, i, tc.d, d)
			}

			e, f := i.Intersect(x), x.Intersect(i)
			we := parseInterval(tc.e)
			if !e.Equal(we) {
				t.Errorf("want %s.Intersect(%s) = %s but get %s", i, x, we, e)
			}
			if !f.Equal(we) {
				t.Errorf("want %s.Intersect(%s) = %s but get %s", x, i, we, f)
			}

			g, h := i.Bisect(x)
			wg, wh := parseInterval(tc.g), parseInterval(tc.h)
			if !g.Equal(wg) || !h.Equal(wh) {
				t.Errorf("want %s.Bisect(%s) = %s, %s but get %s, %s", i, x, wg, wh, g, h)
			}
			j, k := x.Bisect(i)
			wj, wk := parseInterval(tc.j), parseInterval(tc.k)
			if !j.Equal(wj) || !k.Equal(wk) {
				t.Errorf("want %s.Bisect(%s) = %s, %s but get %s, %s", x, i, wj, wk, k, k)
			}

			l, m := i.Adjoin(x), x.Adjoin(i)
			wl := parseInterval(tc.l)
			if !l.Equal(wl) {
				t.Errorf("want %s.Adjoin(%s) = %s but get %s", i, x, wl, l)
			}
			if !m.Equal(wl) {
				t.Errorf("want %s.Adjoin(%s) = %s but get %s", x, i, wl, m)
			}

			o, p := i.Encompass(x), x.Encompass(i)
			wo := parseInterval(tc.o)
			if !o.Equal(wo) {
				t.Errorf("want %s.Encompass(%s) = %s but get %s", i, x, wo, o)
			}
			if !p.Equal(wo) {
				t.Errorf("want %s.Encompass(%s) = %s but get %s", x, i, wo, p)
			}
		})
	}
}
