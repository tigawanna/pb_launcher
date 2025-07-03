package utils

func StrPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func Ptr[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}
