package main

import (
	"EcoNode/route"
)

func main() {
	engine := route.SetupEngine()
	err := engine.Run()
	if err != nil {
		return
	}
}
