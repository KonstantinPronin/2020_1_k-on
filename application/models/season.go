package models

type Season struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	SeriesID    uint   `json:"seriesId" gorm:"column:seriesid"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	TrailerLink string `json:"trailerLink" gorm:"column:trailerlink"`
	Description string `json:"description"`
	Year        int    `json:"year"`
	Image       string `json:"image,omitempty"`
}

//easyjson:json
type Seasons []Season

func (s *Season) TableName() string {
	return "seasons"
}
