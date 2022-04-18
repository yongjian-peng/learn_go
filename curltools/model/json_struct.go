package model

// request josn structure
type PostWithJson struct {
	Appid string `json:"appid"`
	Sn    string `json:"sn"`
	Sign  string `json:"sign"`
}
