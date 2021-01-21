package main

import "github.com/fatih/color"

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	g      = color.New(color.FgGreen, color.Bold).SprintFunc()
	b      = color.New(color.FgBlue, color.Bold).SprintFunc()
	token  = ""
	//C *SNBot
	C *SNBot
	//Skip bool
	Skip bool
)
