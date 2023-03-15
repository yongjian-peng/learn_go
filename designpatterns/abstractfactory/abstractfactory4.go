package abstractfactory

import "fmt"

// gruelDrink 粥
type gruelDrink struct {
	gruelName string
}

func (g *gruelDrink) Drunk() string {
	return fmt.Sprintf("%v被喝", g.gruelName)
}

// sodaDrink 汽水
type sodaDrink struct {
	sodaName string
}

func (s *sodaDrink) Drunk() string {
	return fmt.Sprintf("%v被喝", s.sodaName)
}

// soupDrink 汤
type soupDrink struct {
	soupName string
}

func (s *soupDrink) Drunk() string {
	return fmt.Sprintf("%v被喝", s.soupName)
}
