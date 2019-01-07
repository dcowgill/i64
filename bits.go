/*

Package i64 implements a bit field that consists of a single 64-bit word.

It can therefore efficiently represent an integer set, but only for values
between 0 and 63, inclusive.

*/
package i64

import (
	"math/bits"
	"strconv"
	"strings"
)

// Bits is a field of 64 bits.
//
// Unless otherwise specified, methods that accept a bit position, such as Set
// and Test, do not check their arguments, so invoking them with values outside
// [0, 63] will either corrupt the bit field or return an incorrect answer.
//
// Bits is defined as a uint64; therefore, it can be copied and compared for
// equality like any built-in integer value.
type Bits uint64

// Of returns a bit field with the specified bits set.
// Any bits outside [0, 63] are ignored.
func Of(bits ...int) Bits {
	var b Bits
	for _, n := range bits {
		if n >= 0 && n < 64 {
			b = b.Set(n)
		}
	}
	return b
}

// Range returns a bit field with the bits in the specified range set.
// Any bits outside [0, 63] are ignored.
func Range(low, high, step int) Bits {
	var b Bits
	if low < 0 {
		low = 0
	}
	if high > 63 {
		high = 63
	}
	for n := low; n <= high; n += step {
		b = b.Set(n)
	}
	return b
}

// Set returns a copy of the bit field that has the nth bit set.
func (b Bits) Set(n int) Bits {
	return b | (1 << uint64(n))
}

// Unset returns a copy of the bit field that has the nth bit unset.
func (b Bits) Unset(n int) Bits {
	return b & ^(1 << uint64(n))
}

// Test reports whether the nth bit in the field is set.
func (b Bits) Test(n int) bool {
	return b&(1<<uint64(n)) != 0
}

// Empty reports whether the bit field is empty, i.e. has zero bits set.
func (b Bits) Empty() bool {
	return b == 0
}

// Count reports the number of bits in the field that are set.
func (b Bits) Count() int {
	return bits.OnesCount64(uint64(b))
}

// Singular reports whether the bit field has exactly one set bit.
func (b Bits) Singular() bool {
	return b != 0 && (b&(b-1)) == 0
}

// Least returns the least significant set bit in the field.
// If the field has no set bits, returns -1.
func (b Bits) Least() int {
	if b == 0 {
		return -1 // empty
	}
	return bits.TrailingZeros64(uint64(b))
}

// Most returns the most significant set bit in the field.
// If the field has no set bits, returns -1.
func (b Bits) Most() int {
	if b == 0 {
		return -1 // empty
	}
	return 63 - bits.LeadingZeros64(uint64(b))
}

// String implements the Stringer interface. It returns a string containing the
// set bits in the field, in ascending order, separated by spaces. For example,
// Bits(0).Set(1).Set(3).Set(5).String() returns "1 3 5".
func (b Bits) String() string {
	var sb strings.Builder
	var sep string
	it := b.Iter()
	for x := it.Next(); x >= 0; x = it.Next() {
		sb.WriteString(sep)
		sb.WriteString(strconv.Itoa(x))
		sep = " "
	}
	return sb.String()
}

// Iter returns an iterator over the bits in the field.
func (b Bits) Iter() Iter {
	return Iter(b)
}

// Iter iterates over the set bits in a bit field.
//
// Example usage:
//
//		var b i64.Bits
//		b.Set(5)
//		b.Set(10)
//		it := b.Iter()
//		for x := it.Next(); x >= 0; x = it.Next() {
//			fmt.Println(x)
//		}
//
// The above example will print "5" and "10" to stdout.
//
// Note that Iter makes a copy of the bit field; subsequent changes to the field
// cannot affect the iterator.
type Iter uint64

// Next returns the next bit in the field.
// If the iterator is exhausted, returns -1.
func (it *Iter) Next() int {
	b := uint64(*it)
	if b == 0 {
		return -1
	}
	n := bits.TrailingZeros64(b)
	*it = Iter(b & (b - 1))
	return n
}
