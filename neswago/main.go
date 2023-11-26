package main

import (
	_ "embed"
	"flag"
	"fmt"

	"github.com/fogleman/nes/ui"
)

//go:embed static/roms/BattleCity.nes
var nesRom_BattleCity []byte

//go:embed static/roms/GongLuSaiChe.nes
var nesRom_GongLuSaiChe []byte

//go:embed static/roms/HunDouLuo1_S_30.nes
var nesRom_HunDouLuo1_S_30 []byte

//go:embed static/roms/SuperMarioBros.nes
var nesRom_SuperMarioBros []byte

var allGames = []struct {
	Name string
	Data []byte
}{
	{"BattleCity", nesRom_BattleCity},
	{"GongLuSaiChe", nesRom_GongLuSaiChe},
	{"HunDouLuo1_S_30", nesRom_HunDouLuo1_S_30},
	{"SuperMarioBros", nesRom_SuperMarioBros},
}

var (
	flagIndex = flag.Int("n", 0, "set game index(0-3)")
	flagList  = flag.Bool("l", false, "show game list")
)

func main() {
	flag.Parse()

	if *flagList {
		for i, x := range allGames {
			fmt.Printf("[%d] - %s\n", i, x.Name)
		}
		return
	}

	game := allGames[*flagIndex]
	ui.Main(game.Name, game.Data)
}
