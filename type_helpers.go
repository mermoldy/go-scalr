package scalr

// Bool returns a pointer to the given bool
func Bool(v bool) *bool {
	return &v
}

// Category returns a pointer to the given category type.
func Category(v CategoryType) *CategoryType {
	return &v
}

// Int returns a pointer to the given int.
func Int(v int) *int {
	return &v
}

// Int64 returns a pointer to the given int64.
func Int64(v int64) *int64 {
	return &v
}

// UInt16 returns a pointer to the given uint16.
func UInt16(v uint16) *uint16 {
	return &v
}

// String returns a pointer to the given string.
func String(v string) *string {
	return &v
}
