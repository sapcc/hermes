package storage

// RemoveDuplicates removes duplicates from a slice of strings while preserving the order.
func RemoveDuplicates(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	var result []string
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
