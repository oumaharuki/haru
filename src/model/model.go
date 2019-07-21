package model
type Anime struct{
	Id int64 `db:"id" json:"id"`
	SId int64 `db:"s_id" json:"s_id"`
	Name string `db:"name" json:"name"`
	EmNum string `db:"em_num" json:"em_num"`
	Year string `db:"year" json:"year"`
	Area string `db:"area" json:"area"`
	Picture string `db:"picture" json:"picture"`
	Introduction string `db:"introduction" json:"introduction"`
	Form string `db:"form" json:"form"`
	CreateTime string `db:"create_time" json:"create_time"`
}
type Director struct{
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	AnimeId int `db:"anime_id" json:"anime_id"`
	CreateTime string `db:"create_time" json:"create_time"`
}
type Drama struct{
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	PlayUrl string `db:"play_url" json:"play_url"`
	Source string `db:"source" json:"source"`
	AnimeId int `db:"anime_id" json:"anime_id"`
	PlayMethod string `db:"play_method" json:"play_method"`
	CreateTime string `db:"create_time" json:"create_time"`
}
type Star struct{
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	AnimeId int `db:"anime_id" json:"anime_id"`
	CreateTime string `db:"create_time" json:"create_time"`
}

