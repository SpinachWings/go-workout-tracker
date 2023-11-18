package utils

func RemoveZerosFromSliceOfUint(allUints []uint) []uint {
	var relevantUints []uint
	for _, id := range allUints {
		if id != 0 {
			relevantUints = append(relevantUints, id)
		}
	}
	return relevantUints
}

func CombineTwoSlicesOfStringNoDuplicates(s1 []string, s2 []string) []string {
	combinedSlice := s1
	for _, valOfS2 := range s2 {
		duplicate := false
		for _, valOfS1 := range s1 {
			if valOfS1 == valOfS2 {
				duplicate = true
			}
		}
		if !duplicate {
			combinedSlice = append(combinedSlice, valOfS2)
		}
	}
	return combinedSlice
}

func SliceOfStringContainsDuplicates(slice []string) bool {
	seen := make(map[string]bool)
	for _, value := range slice {
		if seen[value] {
			return true
		}
		seen[value] = true
	}
	return false
}
