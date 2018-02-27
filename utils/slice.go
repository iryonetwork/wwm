package utils

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
