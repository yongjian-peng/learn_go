package main

func main() {
	s := []int{1, 2, 3, 4}
	noraml(s)

	bce(s)
}

func noraml(s []int) {
	i := 0
	i += s[0]
	i += s[1]
	i += s[2]
	i += s[3]
	println(i)
}

func bce(s []int) {

	_ = s[3]
	i := 0
	i += s[0]
	i += s[1]
	i += s[2]
	i += s[3]
	println(i)
}

