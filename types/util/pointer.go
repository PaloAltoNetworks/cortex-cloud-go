package types

// ToPointer takes any type and returns a pointer for that type.
func ToPointer[T any](d T) *T {
	return &d
}
