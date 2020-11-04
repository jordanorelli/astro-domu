package math

func Abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func Max(a, b int, more ...int) int {
	var v int
	if a > b {
		v = a
	} else {
		v = b
	}
	for _, next := range more {
		if next > v {
			v = next
		}
	}
	return v
}

func Min(a, b int, more ...int) int {
	var v int
	if a < b {
		v = a
	} else {
		v = b
	}
	for _, next := range more {
		if next < v {
			v = next
		}
	}
	return v
}
