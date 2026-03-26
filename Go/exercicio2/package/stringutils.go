package stringutils

func ReverseString(inputString string) string {
	res := ""

	for i := len(inputString) - 1; i >= 0; i-- {
		res = res + string(inputString[i])
	}
	return res
}
