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
