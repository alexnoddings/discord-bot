package main

import (
	"fmt"
	"kingsgate/internal/kingsgate"
)

func main() {
	err := kingsgate.Run()
	if err != nil {
		fmt.Println(err)
	}
}
