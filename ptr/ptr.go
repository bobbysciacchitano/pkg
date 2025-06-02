package ptr

func String(v string) *string {
	return &v
}

func Boolean(v bool) *bool {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Float64(v float64) *float64 {
	return &v
}
