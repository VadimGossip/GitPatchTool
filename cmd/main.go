package main

import (
	"github.com/VadimGossip/gitPatchTool/internal/app"
)

var configDir = "config"

func main() {
	gpt := app.NewApp("Git Patch Tool", configDir)
	gpt.Run()
}
