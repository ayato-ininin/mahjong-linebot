package models

import "time"

// int64にしているのは、firestoreからくる数値データはintではなく、int64でくるから。
// https://qiita.com/sorarideblog/items/1dd10042be2ec8e3dcaa
type MatchSetting struct {
	RoomId           int64     `firestore:"roomId" json:"roomId"`
	MahjongNumber    string    `firestore:"mahjongNumber" json:"mahjongNumber"`
	Name1            string    `firestore:"name1" json:"name1"`
	Name2            string    `firestore:"name2" json:"name2"`
	Name3            string    `firestore:"name3" json:"name3"`
	Name4            string    `firestore:"name4" json:"name4"`
	Uma              string    `firestore:"uma" json:"uma"`
	Oka              int64     `firestore:"oka" json:"oka"`
	IsYakitori       bool      `firestore:"isYakitori" json:"isYakitori"`
	IsTobishou       bool      `firestore:"isTobishou" json:"isTobishou"`
	TobishouPoint    int64     `firestore:"tobishouPoint" json:"tobishouPoint"`
	Rate             int64     `firestore:"rate" json:"rate"`
	IsTip            bool      `firestore:"isTip" json:"isTip"`
	TipInitialNumber int64     `firestore:"tipInitialNumber" json:"tipInitialNumber"`
	TipRate          int64     `firestore:"tipRate" json:"tipRate"`
	CreateTimestamp  time.Time `firestore:"createTimestamp" json:"createTimestamp"`
	UpdateTimestamp  time.Time `firestore:"updateTimestamp" json:"updateTimestamp"`
}

// type RespMatchSettingModel struct {
// 	Status         int32 `json:"status"`
// 	Message       string `json:"message"`
// 	MatchSetting *MatchSetting `json:"matchSetting"`
// }
