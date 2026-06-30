package gutils

// Must 检查错误是否为nil，若为nil则panic
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// MustAny 检查错误是否为nil，若为nil则panic，否则返回v
func MustAny(v any, err error) any {
	if err != nil {
		panic(err)
	}
	return v
}

// MustT 检查错误是否为nil，若为nil则panic，否则返回v
func MustT[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
