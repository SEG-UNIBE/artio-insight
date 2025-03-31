package main

import (
	"artio-insight/pkg/scripts"
	"fmt"
)

func main() {
	tags := scripts.GetCurrentNipsTags()

	for nip, tag := range tags {
		fmt.Println(nip, ":", tag)
	}
}