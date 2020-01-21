package main

import (
	"fmt"
	"os/exec"
)

func main() {
	c := exec.Command("py", "./sentiment.py")

	out, err := c.Output()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(string(out))
}
