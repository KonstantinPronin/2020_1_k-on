package models

type Genre struct {
	ID        uint   `json:"ID" gorm:"primary_key"`
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

//easyjson:json
type Genres []Genre

func (g *Genre) TableName() string {
	return "genres"
}
