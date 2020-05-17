package models

type Playlist struct {
	Id           uint          `json:"id" gorm:"primary_key"`
	Name         string        `json:"name" gorm:"column:name"`
	Public       bool          `json:"public" gorm:"column:public"`
	UserId       uint          `json:"userId" gorm:"column:user_id"`
	Films        ListsFilm     `json:"films" gorm:"-"`
	Series       ListSeriesArr `json:"series" gorm:"-"`
	IsSubscribed bool          `json:"isSubscribed" gorm:"-"`
}

//easyjson:json
type Playlists []Playlist

type FilmToPlaylist struct {
	Id     uint `gorm:"primary_key"`
	Pid    uint `gorm:"column:playlist_id"`
	FilmId uint `gorm:"column:film_id"`
}

type SeriesToPlaylist struct {
	Id       uint `gorm:"primary_key"`
	Pid      uint `gorm:"column:playlist_id"`
	SeriesId uint `gorm:"column:series_id"`
}
