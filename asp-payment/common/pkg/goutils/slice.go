package goutils

import (
	"math/rand"
	"time"
)

func DeleteSlice[T comparable](a []T, elem T) []T {
	j := 0
	for _, v := range a {
		if v != elem {
			a[j] = v
			j++
		}
	}
	return a[:j]
}

func InSlice[T comparable](needle T, hyStack []T) bool {
	for _, item := range hyStack {
		if needle == item {
			return true
		}
	}
	return false
}

type reduceType func(interface{}) interface{}
type filterType func(interface{}) bool

// SliceRandList generate an int slice from min to max.
func SliceRandList(min, max int) []int {
	if max < min {
		min, max = max, min
	}
	length := max - min + 1
	t0 := time.Now()
	rand.Seed(int64(t0.Nanosecond()))
	list := rand.Perm(length)
	for index := range list {
		list[index] += min
	}
	return list
}

// SliceMerge merges interface slices to one slice.
func SliceMerge(slice1, slice2 []interface{}) (c []interface{}) {
	c = append(slice1, slice2...)
	return
}

// SliceReduce generates a new slice after parsing every value by reduce function
func SliceReduce(slice []interface{}, a reduceType) (dslice []interface{}) {
	for _, v := range slice {
		dslice = append(dslice, a(v))
	}
	return
}

// SliceRand returns random one from slice.
func SliceRand(a []interface{}) (b interface{}) {
	randNum := rand.Intn(len(a))
	b = a[randNum]
	return
}

// SliceSum sums all values in int64 slice.
func SliceSum(intSlice []int64) (sum int64) {
	for _, v := range intSlice {
		sum += v
	}
	return
}

// SliceFilter generates a new slice after filter function.
func SliceFilter(slice []interface{}, a filterType) (ftSlice []interface{}) {
	for _, v := range slice {
		if a(v) {
			ftSlice = append(ftSlice, v)
		}
	}
	return
}

// SliceDiff returns diff slice of slice1 - slice2.
func SliceDiff[T comparable](slice1, slice2 []T) (diffSlice []T) {
	for _, v := range slice1 {
		if !InSlice(v, slice2) {
			diffSlice = append(diffSlice, v)
		}
	}
	return
}

// SliceIntersect returns slice that are present in all the slice1 and slice2.
func SliceIntersect[T comparable](slice1, slice2 []T) (diffSlice []T) {
	for _, v := range slice1 {
		if InSlice(v, slice2) {
			diffSlice = append(diffSlice, v)
		}
	}
	return
}

// SliceChunk separates one slice to some sized slice.
func SliceChunk(slice []interface{}, size int) (chunkSlice [][]interface{}) {
	if size >= len(slice) {
		chunkSlice = append(chunkSlice, slice)
		return
	}
	end := size
	for i := 0; i <= (len(slice) - size); i += size {
		chunkSlice = append(chunkSlice, slice[i:end])
		end += size
	}
	return
}

// SliceRange generates a new slice from begin to end with step duration of int64 number.
func SliceRange(start, end, step int64) (intSlice []int64) {
	for i := start; i <= end; i += step {
		intSlice = append(intSlice, i)
	}
	return
}

// SlicePad prepends size number of val into slice.
func SlicePad(slice []interface{}, size int, val interface{}) []interface{} {
	if size <= len(slice) {
		return slice
	}
	for i := 0; i < (size - len(slice)); i++ {
		slice = append(slice, val)
	}
	return slice
}

// SliceUnique cleans repeated values in slice.
func SliceUnique[T comparable](slice []T) (uniqueSlice []T) {
	for _, v := range slice {
		if !InSlice(v, uniqueSlice) {
			uniqueSlice = append(uniqueSlice, v)
		}
	}
	return
}

// SliceShuffle shuffles a slice.
func SliceShuffle(slice []interface{}) []interface{} {
	for i := 0; i < len(slice); i++ {
		a := rand.Intn(len(slice))
		b := rand.Intn(len(slice))
		slice[a], slice[b] = slice[b], slice[a]
	}
	return slice
}

func ShuffleIntSlice(list []int) []int {
	dest := make([]int, len(list))
	perm := rand.Perm(len(list))
	for i, v := range perm {
		dest[v] = list[i]
	}
	return dest
}

func ShuffleInt32Slice(list []int32) []int32 {
	dest := make([]int32, len(list))
	perm := rand.Perm(len(list))
	for i, v := range perm {
		dest[v] = list[i]
	}
	return dest
}

func ShuffleStringSlice(list []string) []string {
	dest := make([]string, len(list))
	perm := rand.Perm(len(list))
	for i, v := range perm {
		dest[v] = list[i]
	}
	return dest
}

func ShuffleSlice(list []interface{}) []interface{} {
	dest := make([]interface{}, len(list))
	perm := rand.Perm(len(list))
	for i, v := range perm {
		dest[v] = list[i]
	}
	return dest
}
