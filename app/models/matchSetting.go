package models

import "time"

type MatchSetting struct {
	RoomId           string    `firestore:"roomId" json:"roomId"`
	MahjongNumber    string    `firestore:"mahjongNumber" json:"mahjongNumber"`
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
