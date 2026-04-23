package main

import "fmt"

func main() {
	ret, err := Simple("return 1;")
	fmt.Printf("%#v\n", ret)
	fmt.Println(err)
}
