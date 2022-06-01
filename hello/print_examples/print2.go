package main

import (
	"fmt"
)

func main() {
	// a, b := 3.0, 4.0
	// h := math.Hypot(a, b)

	// Print inserts blanks between arguments when neither is a string.
	// It does not add a newline to the output, so we add one explicitly.
	// fmt.Print("The vector (", a, b, ") has length ", h, ".\n")

	// Println always inserts spaces between its arguments,
	// so it cannot be used to produce the same output as Print in this case;
	// its output has extra spaces.
	// Also, Println always adds a newline to the output.
	// fmt.Println("The vector (", a, b, ") has length", h, ".")

	// Printf provides complete control but is more complex to use.
	// It does not add a newline to the output, so we add one explicitly
	// at the end of the format specifier string.
	// fmt.Printf("The vector (%g %g) has length %g.\n", a, b, h)

	// const name, id = "bueller", 17
	// err := fmt.Errorf("user %q (id %d) not found", name, id)
	// fmt.Println(err.Error())

	// const name, age = "Kim", 22
	// n, err := fmt.Fprint(os.Stdout, name, " is ", age, " years old.\n")

	// // The n and err return values from Fprint are
	// // those returned by the underlying io.Writer.
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Fprint: %v\n", err)
	// }
	// fmt.Print(n, " bytes written.\n")

	// const name, age = "Kim", 22
	// n, err := fmt.Fprintf(os.Stdout, "%s is %d years old.\n", name, age)

	// // The n and err return values from Fprintf are
	// // those returned by the underlying io.Writer.
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Fprintf: %v\n", err)
	// }
	// fmt.Printf("%d bytes written.\n", n)

	// const name, age = "Kim", 22
	// n, err := fmt.Fprintln(os.Stdout, name, "is", age, "years old.")

	// // The n and err return values from Fprintln are
	// // those returned by the underlying io.Writer.
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Fprintln: %v\n", err)
	// }
	// fmt.Println(n, "bytes written.")

	// var (
	// 	i int
	// 	b bool
	// 	s string
	// )
	// r := strings.NewReader("5 true gophers")
	// n, err := fmt.Fscanf(r, "%d %t %s", &i, &b, &s)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Fscanf: %v\n", err)
	// }
	// fmt.Println(i, b, s)
	// fmt.Println(n)
	// 5 true gophers 3

	// s := `dmr 1771 1.61803398875
	// ken2 471828 6.14159
	// ken 271828 3.14159`
	// r := strings.NewReader(s)
	// var a string
	// var b int
	// var c float64
	// for {
	// 	// fmt.Printf("%d: %s, %d, %f\n", r, a, b, c)
	// 	n, err := fmt.Fscanln(r, &a, &b, &c)
	// 	if err == io.EOF {
	// 		fmt.Println("break\n")
	// 		break
	// 	}
	// 	if err != nil {
	// 		fmt.Println("err\n")
	// 		panic(err)
	// 	}
	// 	fmt.Printf("%d: %s, %d, %f\n", n, a, b, c)
	// }

	// const name, age = "Kim", 22
	// fmt.Printf("%s is %d years old.\n", name, age)

	// It is conventional not to worry about any
	// error returned by Printf.

	// const name, age = "Kim", 22
	// fmt.Println(name, "is", age, "years old.")

	// It is conventional not to worry about any
	// error returned by Println.

	// const name, age = "Kim", 22
	// s := fmt.Sprint(name, " is ", age, " years old.\n")

	// io.WriteString(os.Stdout, s) // Ignoring error for simplicity.

	// const name, age = "Kim", 22
	// s := fmt.Sprintf("%s is %d years old.\n", name, age)

	// io.WriteString(os.Stdout, s) // Ignoring error for simplicity.

	// const name, age = "Kim", 22
	// s := fmt.Sprintln(name, "is", age, "years old.")

	// io.WriteString(os.Stdout, s) // Ignoring error for simplicity.

	var name string
	var age int
	n, err := fmt.Sscanf("Kim is 22 years old", "%s is %d years old", &name, &age)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d: %s, %d\n", n, name, age)

}
