package utils

// Ptr возвращает указатель на значение: нужен для необязательных полей запроса, где 0 — допустимое значение
func Ptr[T any](value T) *T {
	return &value
}
