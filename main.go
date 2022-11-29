/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/xuender/comic/cmd"
)

func main() {
	cmd.Execute()
}
