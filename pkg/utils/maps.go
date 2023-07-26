package utils

func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}
