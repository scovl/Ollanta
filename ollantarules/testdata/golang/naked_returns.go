package testdata

// naked_returns.go — triggers go:no-naked-returns

func nakedReturn() (result int, err error) {
	result = 1
	_ = err
	// line 7
	// line 8
	// line 9
	// line 10
	return // naked return in a function with >5 lines
}

func shortNaked() (result int) {
	result = 1
	return // short function — should NOT be flagged
}

func explicitReturn() (result int, err error) {
	result = 42
	_ = err
	// padding
	// padding
	// padding
	// padding
	// padding
	return result, nil // explicit — should NOT be flagged
}
