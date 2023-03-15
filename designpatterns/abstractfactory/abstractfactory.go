package abstractfactory

// Cook 厨师接口，抽象工厂
type Cook interface {
	// MakeFood 制作主食
	MakeFood() Food
	// MakeDrink 制作饮品
	MakeDrink() Drink
}

// Food 主食接口
type Food interface {
	// Eaten 被吃
	Eaten() string
}

// Drink 饮品接口
type Drink interface {
	// Drunk 被喝
	Drunk() string
}
