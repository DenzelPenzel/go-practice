package main

// Clearing maps
func clearMaps(m map[string]int) {
	for k := range m {
		delete(m, k)
	}
}

// Increasing the length of a slice
func incSlice() {
	s := make([]int, 100)
	if len(s) < 200 {
		// Need to extend s
		s = append(s, make([]int, 100)...)
	}
}
