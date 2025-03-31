package main

import (
	"artio-insight/pkg/scripts"
	"fmt"
)

func main() {
	tags,_ := scripts.GetCurrentNipsTags()

	for nip, tag := range tags {
		fmt.Println(nip, ":", tag)
	}
}