package mathx

// MinInt returns minimum value of provided params
func MinInt(is ...int) int {
	res := 0
	for _, i := range is {
		if i < res {
			res = i
		}
	}
	return res
}

// MaxInt returns maximum value of provided params
func MaxInt(is ...int) int {
	res := 0
	for _, i := range is {
		if i > res {
			res = i
		}
	}
	return res
}
