package proxies

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/elliotchance/pie/pie"
	"github.com/elliotchance/pie/pie/util"

	"subsrcibe/proxy"
)

// All will return true if all callbacks return true. It follows the same logic
// as the all() function in Python.
//
// If the list is empty then true is always returned.
func (ss ProxyList) All(fn func(value proxy.Proxy) bool) bool {
	for _, value := range ss {
		if !fn(value) {
			return false
		}
	}

	return true
}

// Any will return true if any callbacks return true. It follows the same logic
// as the any() function in Python.
//
// If the list is empty then false is always returned.
func (ss ProxyList) Any(fn func(value proxy.Proxy) bool) bool {
	for _, value := range ss {
		if fn(value) {
			return true
		}
	}

	return false
}

// Append will return a new slice with the elements appended to the end.
//
// It is acceptable to provide zero arguments.
func (ss ProxyList) Append(elements ...proxy.Proxy) ProxyList {
	// Copy ss, to make sure no memory is overlapping between input and
	// output. See issue #97.
	result := append(ProxyList{}, ss...)

	result = append(result, elements...)
	return result
}

// Bottom will return n elements from bottom
//
// that means that elements is taken from the end of the slice
// for this [1,2,3] slice with n == 2 will be returned [3,2]
// if the slice has less elements then n that'll return all elements
// if n < 0 it'll return empty slice.
func (ss ProxyList) Bottom(n int) (top ProxyList) {
	var lastIndex = len(ss) - 1
	for i := lastIndex; i > -1 && n > 0; i-- {
		top = append(top, ss[i])
		n--
	}

	return
}

// Contains returns true if the element exists in the slice.
//
// When using slices of pointers it will only compare by address, not value.
func (ss ProxyList) Contains(lookingFor proxy.Proxy) bool {
	for _, s := range ss {
		if lookingFor == s {
			return true
		}
	}

	return false
}

// Diff returns the elements that needs to be added or removed from the first
// slice to have the same elements in the second slice.
//
// The order of elements is not taken into consideration, so the slices are
// treated sets that allow duplicate items.
//
// The added and removed returned may be blank respectively, or contain upto as
// many elements that exists in the largest slice.
func (ss ProxyList) Diff(against ProxyList) (added, removed ProxyList) {
	// This is probably not the best way to do it. We do an O(n^2) between the
	// slices to see which items are missing in each direction.

	diffOneWay := func(ss1, ss2raw ProxyList) (result ProxyList) {
		ss2 := make(ProxyList, len(ss2raw))
		copy(ss2, ss2raw)

		for _, s := range ss1 {
			found := false

			for i, element := range ss2 {
				if s == element {
					ss2 = append(ss2[:i], ss2[i+1:]...)
					found = true
				}
			}

			if !found {
				result = append(result, s)
			}
		}

		return
	}

	removed = diffOneWay(ss, against)
	added = diffOneWay(against, ss)

	return
}

// DropTop will return the rest slice after dropping the top n elements
// if the slice has less elements then n that'll return empty slice
// if n < 0 it'll return empty slice.
func (ss ProxyList) DropTop(n int) (drop ProxyList) {
	if n < 0 || n >= len(ss) {
		return
	}

	// Copy ss, to make sure no memory is overlapping between input and
	// output. See issue #145.
	drop = make([]proxy.Proxy, len(ss)-n)
	copy(drop, ss[n:])

	return
}

// Each is more condensed version of Transform that allows an action to happen
// on each elements and pass the original slice on.
//
//   cars.Each(func (car *Car) {
//       fmt.Printf("Car color is: %s\n", car.Color)
//   })
//
// Pie will not ensure immutability on items passed in so they can be
// manipulated, if you choose to do it this way, for example:
//
//   // Set all car colors to Red.
//   cars.Each(func (car *Car) {
//       car.Color = "Red"
//   })
//
func (ss ProxyList) Each(fn func(proxy.Proxy)) ProxyList {
	for _, s := range ss {
		fn(s)
	}

	return ss
}

// Equals compare elements from the start to the end,
//
// if they are the same is considered the slices are equal if all elements are the same is considered the slices are equal
// if each slice == nil is considered that they're equal
//
// if element realizes Equals interface it uses that method, in other way uses default compare
func (ss ProxyList) Equals(rhs ProxyList) bool {
	if len(ss) != len(rhs) {
		return false
	}

	for i := range ss {
		if !(ss[i] == rhs[i]) {
			return false
		}
	}

	return true
}

// Extend will return a new slice with the slices of elements appended to the
// end.
//
// It is acceptable to provide zero arguments.
func (ss ProxyList) Extend(slices ...ProxyList) (ss2 ProxyList) {
	ss2 = ss

	for _, slice := range slices {
		ss2 = ss2.Append(slice...)
	}

	return ss2
}

// Filter will return a new slice containing only the elements that return
// true from the condition. The returned slice may contain zero elements (nil).
//
// FilterNot works in the opposite way of Filter.
func (ss ProxyList) Filter(condition func(proxy.Proxy) bool) (ss2 ProxyList) {
	for _, s := range ss {
		if condition(s) {
			ss2 = append(ss2, s)
		}
	}
	return
}

// FilterNot works the same as Filter, with a negated condition. That is, it will
// return a new slice only containing the elements that returned false from the
// condition. The returned slice may contain zero elements (nil).
func (ss ProxyList) FilterNot(condition func(proxy.Proxy) bool) (ss2 ProxyList) {
	for _, s := range ss {
		if !condition(s) {
			ss2 = append(ss2, s)
		}
	}

	return
}

// FindFirstUsing will return the index of the first element when the callback returns true or -1 if no element is found.
// It follows the same logic as the findIndex() function in Javascript.
//
// If the list is empty then -1 is always returned.
func (ss ProxyList) FindFirstUsing(fn func(value proxy.Proxy) bool) int {
	for idx, value := range ss {
		if fn(value) {
			return idx
		}
	}

	return -1
}

// First returns the first element, or zero. Also see FirstOr().
func (ss ProxyList) First() proxy.Proxy {
	return ss.FirstOr(proxy.Proxy{})
}

// FirstOr returns the first element or a default value if there are no
// elements.
func (ss ProxyList) FirstOr(defaultValue proxy.Proxy) proxy.Proxy {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[0]
}

// Float64s transforms each element to a float64.
func (ss ProxyList) Float64s() pie.Float64s {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(pie.Float64s, l)
	for i := 0; i < l; i++ {
		mightBeString := ss[i]
		result[i], _ = strconv.ParseFloat(fmt.Sprintf("%v", mightBeString), 64)
	}

	return result
}

// Insert a value at an index
func (ss ProxyList) Insert(index int, values ...proxy.Proxy) ProxyList {
	if index >= ss.Len() {
		return ProxyList.Extend(ss, ProxyList(values))
	}

	return ProxyList.Extend(ss[:index], ProxyList(values), ss[index:])
}

// Ints transforms each element to an integer.
func (ss ProxyList) Ints() pie.Ints {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(pie.Ints, l)
	for i := 0; i < l; i++ {
		mightBeString := ss[i]
		f, _ := strconv.ParseFloat(fmt.Sprintf("%v", mightBeString), 64)
		result[i] = int(f)
	}

	return result
}

// Join returns a string from joining each of the elements.
func (ss ProxyList) Join(glue string) (s string) {
	var slice interface{} = []proxy.Proxy(ss)

	if y, ok := slice.([]string); ok {
		// The stdlib is efficient for type []string
		return strings.Join(y, glue)
	} else {
		// General case
		parts := make([]string, len(ss))
		for i, element := range ss {
			mightBeString := element
			parts[i] = fmt.Sprintf("%v", mightBeString)
		}
		return strings.Join(parts, glue)
	}
}

// JSONBytes returns the JSON encoded array as bytes.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array.
func (ss ProxyList) JSONBytes() []byte {
	if ss == nil {
		return []byte("[]")
	}

	// An error should not be possible.
	data, _ := json.Marshal(ss)

	return data
}

// JSONBytesIndent returns the JSON encoded array as bytes with indent applied.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array. See
// json.MarshalIndent for details.
func (ss ProxyList) JSONBytesIndent(prefix, indent string) []byte {
	if ss == nil {
		return []byte("[]")
	}

	// An error should not be possible.
	data, _ := json.MarshalIndent(ss, prefix, indent)

	return data
}

// JSONString returns the JSON encoded array as a string.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array.
func (ss ProxyList) JSONString() string {
	if ss == nil {
		return "[]"
	}

	// An error should not be possible.
	data, _ := json.Marshal(ss)

	return string(data)
}

// JSONStringIndent returns the JSON encoded array as a string with indent applied.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array. See
// json.MarshalIndent for details.
func (ss ProxyList) JSONStringIndent(prefix, indent string) string {
	if ss == nil {
		return "[]"
	}

	// An error should not be possible.
	data, _ := json.MarshalIndent(ss, prefix, indent)

	return string(data)
}

// Last returns the last element, or zero. Also see LastOr().
func (ss ProxyList) Last() proxy.Proxy {
	return ss.LastOr(proxy.Proxy{})
}

// LastOr returns the last element or a default value if there are no elements.
func (ss ProxyList) LastOr(defaultValue proxy.Proxy) proxy.Proxy {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[len(ss)-1]
}

// Len returns the number of elements.
func (ss ProxyList) Len() int {
	return len(ss)
}

// Map will return a new slice where each element has been mapped (transformed).
// The number of elements returned will always be the same as the input.
//
// Be careful when using this with slices of pointers. If you modify the input
// value it will affect the original slice. Be sure to return a new allocated
// object or deep copy the existing one.
func (ss ProxyList) Map(fn func(proxy.Proxy) proxy.Proxy) (ss2 ProxyList) {
	if ss == nil {
		return nil
	}

	ss2 = make([]proxy.Proxy, len(ss))
	for i, s := range ss {
		ss2[i] = fn(s)
	}

	return
}

// Mode returns a new slice containing the most frequently occuring values.
//
// The number of items returned may be the same as the input or less. It will
// never return zero items unless the input slice has zero items.
func (ss ProxyList) Mode() ProxyList {
	if len(ss) == 0 {
		return nil
	}
	values := make(map[proxy.Proxy]int)
	for _, s := range ss {
		values[s]++
	}

	var maxFrequency int
	for _, v := range values {
		if v > maxFrequency {
			maxFrequency = v
		}
	}

	var maxValues ProxyList
	for k, v := range values {
		if v == maxFrequency {
			maxValues = append(maxValues, k)
		}
	}

	return maxValues
}

// Pop the first element of the slice
//
// Usage Example:
//
//   type knownGreetings []string
//   greetings := knownGreetings{"ciao", "hello", "hola"}
//   for greeting := greetings.Pop(); greeting != nil; greeting = greetings.Pop() {
//       fmt.Println(*greeting)
//   }
func (ss *ProxyList) Pop() (popped *proxy.Proxy) {

	if len(*ss) == 0 {
		return
	}

	popped = &(*ss)[0]
	*ss = (*ss)[1:]
	return
}

// Random returns a random element by your rand.Source, or zero
func (ss ProxyList) Random(source rand.Source) proxy.Proxy {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 1 {
		return proxy.Proxy{}
	}
	if n < 2 {
		return ss[0]
	}
	rnd := rand.New(source)
	i := rnd.Intn(n)
	return ss[i]
}

// Reverse returns a new copy of the slice with the elements ordered in reverse.
// This is useful when combined with Sort to get a descending sort order:
//
//   ss.Sort().Reverse()
//
func (ss ProxyList) Reverse() ProxyList {
	// Avoid the allocation. If there is one element or less it is already
	// reversed.
	if len(ss) < 2 {
		return ss
	}

	sorted := make([]proxy.Proxy, len(ss))
	for i := 0; i < len(ss); i++ {
		sorted[i] = ss[len(ss)-i-1]
	}

	return sorted
}

// Send sends elements to channel
// in normal act it sends all elements but if func canceled it can be less
//
// it locks execution of gorutine
// it doesn't close channel after work
// returns sended elements if len(this) != len(old) considered func was canceled
func (ss ProxyList) Send(ctx context.Context, ch chan<- proxy.Proxy) ProxyList {
	for i, s := range ss {
		select {
		case <-ctx.Done():
			return ss[:i]
		default:
			ch <- s
		}
	}

	return ss
}

// SequenceUsing generates slice in range using creator function
//
// There are 3 variations to generate:
// 		1. [0, n).
//		2. [min, max).
//		3. [min, max) with step.
//
// if len(params) == 1 considered that will be returned slice between 0 and n,
// where n is the first param, [0, n).
// if len(params) == 2 considered that will be returned slice between min and max,
// where min is the first param, max is the second, [min, max).
// if len(params) > 2 considered that will be returned slice between min and max with step,
// where min is the first param, max is the second, step is the third one, [min, max) with step,
// others params will be ignored
func (ss ProxyList) SequenceUsing(creator func(int) proxy.Proxy, params ...int) ProxyList {
	var seq = func(min, max, step int) (seq ProxyList) {
		lenght := int(util.Round(float64(max-min) / float64(step)))
		if lenght < 1 {
			return
		}

		seq = make(ProxyList, lenght)
		for i := 0; i < lenght; min += step {
			seq[i] = creator(min)
			i++
		}

		return seq
	}

	if len(params) > 2 {
		return seq(params[0], params[1], params[2])
	} else if len(params) == 2 {
		return seq(params[0], params[1], 1)
	} else if len(params) == 1 {
		return seq(0, params[0], 1)
	} else {
		return nil
	}
}

// Shift will return two values: the shifted value and the rest slice.
func (ss ProxyList) Shift() (proxy.Proxy, ProxyList) {
	return ss.First(), ss.DropTop(1)
}

// Shuffle returns shuffled slice by your rand.Source
func (ss ProxyList) Shuffle(source rand.Source) ProxyList {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 2 {
		return ss
	}

	// go 1.10+ provides rnd.Shuffle. However, to support older versions we copy
	// the algorithm directly from the go source: src/math/rand/rand.go below,
	// with some adjustments:
	shuffled := make([]proxy.Proxy, n)
	copy(shuffled, ss)

	rnd := rand.New(source)

	util.Shuffle(rnd, n, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

// SortStableUsing works similar to sort.SliceStable. However, unlike sort.SliceStable the
// slice returned will be reallocated as to not modify the input slice.
func (ss ProxyList) SortStableUsing(less func(a, b proxy.Proxy) bool) ProxyList {
	// Avoid the allocation. If there is one element or less it is already
	// sorted.
	if len(ss) < 2 {
		return ss
	}

	sorted := make(ProxyList, len(ss))
	copy(sorted, ss)
	sort.SliceStable(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})

	return sorted
}

// SortUsing works similar to sort.Slice. However, unlike sort.Slice the
// slice returned will be reallocated as to not modify the input slice.
func (ss ProxyList) SortUsing(less func(a, b proxy.Proxy) bool) ProxyList {
	// Avoid the allocation. If there is one element or less it is already
	// sorted.
	if len(ss) < 2 {
		return ss
	}

	sorted := make(ProxyList, len(ss))
	copy(sorted, ss)
	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})

	return sorted
}

// Strings transforms each element to a string.
//
// If the element type implements fmt.Stringer it will be used. Otherwise it
// will fallback to the result of:
//
//   fmt.Sprintf("%v")
//
func (ss ProxyList) Strings() pie.Strings {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(pie.Strings, l)
	for i := 0; i < l; i++ {
		mightBeString := ss[i]
		result[i] = fmt.Sprintf("%v", mightBeString)
	}

	return result
}

// SubSlice will return the subSlice from start to end(excluded)
//
// Condition 1: If start < 0 or end < 0, nil is returned.
// Condition 2: If start >= end, nil is returned.
// Condition 3: Return all elements that exist in the range provided,
// if start or end is out of bounds, zero items will be placed.
func (ss ProxyList) SubSlice(start int, end int) (subSlice ProxyList) {
	if start < 0 || end < 0 {
		return
	}

	if start >= end {
		return
	}

	length := ss.Len()
	if start < length {
		if end <= length {
			subSlice = ss[start:end]
		} else {
			zeroArray := make([]proxy.Proxy, end-length)
			subSlice = ss[start:length].Append(zeroArray[:]...)
		}
	} else {
		zeroArray := make([]proxy.Proxy, end-start)
		subSlice = zeroArray[:]
	}

	return
}

// Top will return n elements from head of the slice
// if the slice has less elements then n that'll return all elements
// if n < 0 it'll return empty slice.
func (ss ProxyList) Top(n int) (top ProxyList) {
	for i := 0; i < len(ss) && n > 0; i++ {
		top = append(top, ss[i])
		n--
	}

	return
}

// StringsUsing transforms each element to a string.
func (ss ProxyList) StringsUsing(transform func(proxy.Proxy) string) pie.Strings {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(pie.Strings, l)
	for i := 0; i < l; i++ {
		result[i] = transform(ss[i])
	}

	return result
}

// Unshift adds one or more elements to the beginning of the slice
// and returns the new slice.
func (ss ProxyList) Unshift(elements ...proxy.Proxy) (unshift ProxyList) {
	unshift = append(ProxyList{}, elements...)
	unshift = append(unshift, ss...)

	return
}