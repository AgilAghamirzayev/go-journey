package main

func main() {
	//println(intToRoman(1994))
	//println(intToRoman(58))
	println(intToRoman(3749))
}

func intToRoman(num int) string {
	intValues := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	romanValues := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	i := 0
	res := ""

	for num > 0 {
		for num >= intValues[i] {
			res += romanValues[i]
			num -= intValues[i]
		}
		i++
	}

	return res
}
