package model

type Dramas struct {
	Url  string
	Name string
}
type AnimeData struct {
	Title        string
	EmNum        string
	Year         string
	Area         string
	Star         string
	Director     string
	Picture      string
	Introduction string
	Form         string
	Drama        []DramaData
}
type DramaPlay struct {
	Flay   string
	Encry  int
	Link   string
	Name   string
	From   string
	Trysee int
	Url    string
}
type DramaData struct {
	Name       string
	From       string
	Url        string
	PlayMethod string
}
type AnimeInfo struct {
	Anime
	Director string
	Star     string
	Drama    map[string][]Drama
}
