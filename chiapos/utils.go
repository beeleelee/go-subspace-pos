package chiapos

func divCeil(a, b uint) uint {
	res := a / b
	if a%b > 0 {
		res++
	}
	return res
}
