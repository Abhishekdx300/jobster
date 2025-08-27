package helpers

func ValidateSize(maxSize int, sizes ...int) bool {
	for _, size := range sizes {
		if size > maxSize {
			return false
		}
	}
	return true
}
