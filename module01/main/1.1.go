package main

import "fmt"

/*
给定一个字符串数组
[“I”,“am”,“stupid”,“and”,“weak”]
用 for 循环遍历该数组并修改为
[“I”,“am”,“smart”,“and”,“strong”]
*/
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
