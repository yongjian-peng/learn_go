package check

var checker *Check

func ValidatorInit() *Check {
	checker = NewCheck()
	return checker
}
