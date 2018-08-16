package utils

// DiffSlice returns slice that contains elements that are in slice 1 but are not in slice 2
func DiffSlice(slice1, slice2 []string) []string {
	var result []string
	for _, s1 := range slice1 {
		inSlice2 := false
		for _, s2 := range slice2 {
			if s2 == s1 {
				inSlice2 = true
			}
		}
		if !inSlice2 {
			result = append(result, s1)
		}
	}

	return result
}

// IntersectionSlice returns slice that contains elements that are in both slices
func IntersectionSlice(slice1, slice2 []string) []string {
	var result []string
	for _, s1 := range slice1 {
		inSlice2 := false
		for _, s2 := range slice2 {
			if s2 == s1 {
				inSlice2 = true
			}
		}
		if inSlice2 {
			result = append(result, s1)
		}
	}

	return result
}

// SliceContains returns boolean value indicating if string needle is in slice
func SliceContains(slice []string, needle string) bool {
	for _, s := range slice {
		if s == needle {
			return true
		}
	}

	return false
}

// SliceContainsAny returns boolean value indicating if any of the string in needles is in slice
func SliceContainsAny(slice []string, needles []string) bool {
	for _, needle := range needles {
		if SliceContains(slice, needle) {
			return true
		}
	}

	return false
}

// SliceToMap converts strings slice to map with keys being elements of the slice and value set to true
func SliceToMap(slice []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range slice {
		m[s] = true
	}

	return m
}
