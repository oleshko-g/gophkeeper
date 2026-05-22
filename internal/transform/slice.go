package transform

// SliceToSlice transforms a T1 slice into T2 slice of the same size
func SliceToSlice[T1 any, T2 any](sliceOf []T1, transform func(T1) (to T2)) (dst []T2) {
	dst = make([]T2, len(sliceOf))
	for i, v := range sliceOf {
		dst[i] = transform(v)
	}

	return dst
}
