package utils

func StrPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
