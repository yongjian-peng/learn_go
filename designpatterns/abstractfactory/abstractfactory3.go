package abstractfactory

import "fmt"

// cakeFood 蛋糕
type cakeFood struct {
	cakeName string
}

func (c *cakeFood) Eaten() string {
	return fmt.Sprintf("%v被吃", c.cakeName)
}

// dishFood 菜肴
type dishFood struct {
	dishName string
}

func (d *dishFood) Eaten() string {
	return fmt.Sprintf("%v被吃", d.dishName)
}

// noodleFood 面条
type noodleFood struct {
	noodleName string
}

func (n *noodleFood) Eaten() string {
	return fmt.Sprintf("%v被吃", n.noodleName)
}
