package i64

import (
	"reflect"
	"testing"
)

func TestBits(t *testing.T) {
	var (
		b Bits

		// If got != want, fails the test. Assumes got was returned by "method".
		check = func(method string, got, want interface{}) {
			t.Helper()
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("Bits(%s).%s returned %v, want %v", b, method, got, want)
			}
		}

		// Verifies the string representation of b.
		checkstring = func(want string) {
			t.Helper()
			if s := b.String(); s != want {
				t.Fatalf("Bits.String() returned %q, want %q", s, want)
			}
		}

		// Verifies the result of iterating over b.
		checkiter = func(want ...int) {
			t.Helper()
			var xs []int
			it := b.Iter()
			for x := it.Next(); x >= 0; x = it.Next() {
				xs = append(xs, x)
			}
			if !reflect.DeepEqual(xs, want) {
				t.Fatalf("iterating over Bits(%s) returned %+v, want %+v", b, xs, want)
			}
		}
	)

	checkstring("")
	checkiter()
	check("Count()", b.Count(), 0)
	check("Singular()", b.Singular(), false)
	check("Empty()", b.Empty(), true)

	b = b.Set(43)
	checkiter(43)
	checkstring("43")
	check("Count()", b.Count(), 1)
	check("Singular()", b.Singular(), true)
	check("Empty()", b.Empty(), false)
	check("Least()", b.Least(), 43)
	check("Most()", b.Most(), 43)

	b = b.Set(41)
	checkiter(41, 43)
	checkstring("41 43")
	check("Count()", b.Count(), 2)
	check("Singular()", b.Singular(), false)
	check("Empty()", b.Empty(), false)
	check("Least()", b.Least(), 41)
	check("Most()", b.Most(), 43)
	check("Test(41)", b.Test(41), true)
	check("Test(42)", b.Test(42), false)
	check("Test(43)", b.Test(43), true)

	b = b.Unset(41).Unset(43)
	checkiter()
	checkstring("")
	check("Count()", b.Count(), 0)
	check("Singular()", b.Singular(), false)
	check("Empty()", b.Empty(), true)

	b = b.Set(2).Set(4).Set(5).Set(0).Set(12).Set(63)
	checkiter(0, 2, 4, 5, 12, 63)
	checkstring("0 2 4 5 12 63")
	check("Count()", b.Count(), 6)
	check("Singular()", b.Singular(), false)
	check("Empty()", b.Empty(), false)
	check("Least()", b.Least(), 0)
	check("Most()", b.Most(), 63)

	b = b.Unset(5).Unset(63).Unset(0)
	checkiter(2, 4, 12)
	checkstring("2 4 12")
	check("Count()", b.Count(), 3)
	check("Singular()", b.Singular(), false)
	check("Empty()", b.Empty(), false)
	check("Least()", b.Least(), 2)
	check("Most()", b.Most(), 12)
}
