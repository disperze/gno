package main

var time int

type time string

func main() {
	var t time = "hello"
	println(t)
}

// Error:
// files/redeclaration-global1.go:5:6: time redeclared in this block
//	previous declaration at files/redeclaration-global1.go:3:5
