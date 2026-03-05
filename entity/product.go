package entity

type Product struct {
	Id         string  `json:"id" gorm:"primaryKey;unique;not null"`
	Definition string  `json:"definition"`
	Name       string  `json:"name" gorm:"not null;size:64"`
	Price      float32 `json:"price" gorm:"not null;precision:2"`
	Image      string  `json:"image"`
}
