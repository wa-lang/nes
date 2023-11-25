package main

import "github.com/fogleman/nes/ui"

func main() {
	ui.Run([]string{
		"../roms/BattleCity.nes",
		"../roms/GongLuSaiChe.nes",
		"../roms/HunDouLuo1_S_30.nes",
		"../roms/SuperMarioBros.nes",
	})
}
