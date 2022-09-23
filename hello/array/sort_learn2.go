package main

import (
	"fmt"
	"sort"
)

type HeroKind int

const (
	None = iota
	Tank
	Assassin
	Mage
)

type Hero struct {
	Name string
	Kind HeroKind
}

func main() {
	heros := []Hero{
		{"吕布", Tank},
		{"李白", Assassin},
		{"妲己", Mage},
		{"貂蝉", Assassin},
		{"关羽", Tank},
		{"诸葛亮", Mage},
	}
	sort.Slice(heros, func(i, j int) bool {
		if heros[i].Kind != heros[j].Kind {
			return heros[i].Kind < heros[j].Kind
		}
		return heros[i].Name < heros[j].Name
	})
	for _, v := range heros {
		fmt.Printf("%+v\n", v)
	}
}
