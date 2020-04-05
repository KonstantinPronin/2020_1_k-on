package models

type Person struct {
	Id         uint          `json:"id" gorm:"primary_key"`
	Name       string        `json:"name" gorm:"column:name"`
	Occupation string        `json:"occupation" gorm:"column:occupation"`
	BirthDate  string        `json:"birthDate" gorm:"column:birth_date"`
	BirthPlace string        `json:"birthPlace" gorm:"column:birth_place"`
	Image      string        `json:"image" gorm:"column:image"`
	Films      ListsFilm     `json:"films" gorm:"-"`
	Series     ListSeriesArr `json:"series" gorm:"-"`
}

type ListPerson struct {
	Id   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"column:name"`
}

//easyjson:json
type ListPersonArr []ListPerson
