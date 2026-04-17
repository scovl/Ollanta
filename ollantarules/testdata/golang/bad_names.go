package testdata

// bad_names.go — triggers go:naming-conventions

// Get_Value uses underscore in exported name.
func Get_Value() int { return 0 }

// MAXSIZE uses ALL_CAPS style.
func MAXSIZE() int { return 100 }

// goodName is fine (unexported, no underscore).
func goodName() int { return 0 }
