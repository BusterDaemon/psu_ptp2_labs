package entity

type Product struct {
	Id         string  `json:"id" gorm:"primaryKey;unique;not null"`
	Definition string  `json:"definition"`
	Name       string  `json:"name" gorm:"not null;size:64"`
	Price      float32 `json:"price" gorm:"not null;precision:2"`
	Image      string  `json:"image"`
}

func (p Product) GetID() string {
	return p.Id
}

func (p Product) GetDefinition() string {
	return p.Definition
}

func (p Product) GetName() string {
	return p.Name
}

func (p Product) GetPrice() float32 {
	return p.Price
}

func (p Product) GetImage() string {
	return p.Image
}

func (p *Product) SetDefinition(d string) {
	p.Definition = d
}

func (p *Product) SetName(n string) {
	p.Name = n
}

func (p *Product) SetPrice(pr float32) {
	p.Price = pr
}

func (p *Product) SetImage(i string) {
	p.Image = i
}
