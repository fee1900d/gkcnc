package main

import "fmt"

func main() {
	var src = [5]string{"I", "am", "stupid", "and", "weak"}

	for i, s := range src {
		if s == "stupid" {
			src[i] = "smart"
		} else if s == "weak" {
			src[i] = "strong"
		}
	}

	fmt.Println(src)
}
