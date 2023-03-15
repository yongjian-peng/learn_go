package abstractfactory

// breakfastCook 早餐厨师
type breakfastCook struct{}

func NewBreakfastCook() *breakfastCook {
	return &breakfastCook{}
}

func (b *breakfastCook) MakeFood() Food {
	return &cakeFood{"切片面包"}
}

func (b *breakfastCook) MakeDrink() Drink {
	return &gruelDrink{"小米粥"}
}

// lunchCook 午餐厨师
type lunchCook struct{}

func NewLunchCook() *lunchCook {
	return &lunchCook{}
}

func (l *lunchCook) MakeFood() Food {
	return &dishFood{"烤全羊"}
}

func (l *lunchCook) MakeDrink() Drink {
	return &sodaDrink{"冰镇可口可乐"}
}

// dinnerCook 晚餐厨师
type dinnerCook struct{}

func NewDinnerCook() *dinnerCook {
	return &dinnerCook{}
}

func (d *dinnerCook) MakeFood() Food {
	return &noodleFood{"大盘鸡拌面"}
}

func (d *dinnerCook) MakeDrink() Drink {
	return &soupDrink{"西红柿鸡蛋汤"}
}
