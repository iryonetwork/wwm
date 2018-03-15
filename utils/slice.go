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

// SliceContains returns boolean value indicating if string s is in slice
func SliceContains(slice []string, needle string) bool {
	for _, s := range slice {
		if s == needle {
			return true
		}
	}

	return false
}
