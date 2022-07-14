package utils

//ContainsUint checks if a value is present in a slice
func ContainsUint(slice []uint64, value uint64) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}
