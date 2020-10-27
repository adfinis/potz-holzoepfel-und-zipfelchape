package main

import "github.com/adfinis-sygroup/potz-holzoepfel-und-zipfelchape/cmd"

func main() {
	cmd.Execute()
}

//go:generate statik -f --include=*.html
