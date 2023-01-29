package models

import "time"

type MatchSetting struct {
	RoomId           string    `firestore:"roomId" json:"roomId"`
	MahjongNumber    string    `firestore:"mahjongNumber" json:"mahjongNumber"`
	Name1            string    `firestore:"name1" json:"name1"`
	Name2            string    `firestore:"name2" json:"name2"`
	Name3            string    `firestore:"name3" json:"name3"`
	Name4            string    `firestore:"name4" json:"name4"`
	Uma              string    `firestore:"uma" json:"uma"`
	Oka              int       `firestore:"oka" json:"oka"`
	IsYakitori       bool      `firestore:"isYakitori" json:"isYakitori"`
	IsTobishou       bool      `firestore:"isTobishou" json:"isTobishou"`
	TobishouPoint    int       `firestore:"tobishouPoint" json:"tobishouPoint"`
	Rate             int       `firestore:"rate" json:"rate"`
	IsTip            bool      `firestore:"isTip" json:"isTip"`
	TipInitialNumber int       `firestore:"tipInitialNumber" json:"tipInitialNumber"`
	TipRate          int       `firestore:"tipRate" json:"tipRate"`
	CreateTimestamp  time.Time `firestore:"createTimestamp" json:"createTimestamp"`
	UpdateTimestamp  time.Time `firestore:"updateTimestamp" json:"updateTimestamp"`
}

// type RespMatchSettingModel struct {
// 	Status         int32 `json:"status"`
// 	Message       string `json:"message"`
// 	MatchSetting *MatchSetting `json:"matchSetting"`
// }
