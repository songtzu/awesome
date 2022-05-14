package mq

func Contains[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}
