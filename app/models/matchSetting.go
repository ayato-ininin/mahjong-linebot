package models

import "time"

type MatchSetting struct {
	RoomId           string    `firestore:"roomId"`
	MahjongNumber    string    `firestore:"mahjongNumber"`
	Uma              string    `firestore:"uma"`
	Oka              int       `firestore:"oka"`
	IsYakitori       bool      `firestore:"isYakitori"`
	IsTobishou       bool      `firestore:"isTobishou"`
	TobishouPoint    int       `firestore:"tobishouPoint"`
	Rate             int       `firestore:"rate"`
	IsTip            bool      `firestore:"isTip"`
	TipInitialNumber int       `firestore:"tipInitialNumber"`
	TipRate          int       `firestore:"tipRate"`
	CreateTimestamp  time.Time `firestore:"createTimestamp"`
	UpdateTimestamp  time.Time `firestore:"updateTimestamp"`
}

// type RespMatchSettingModel struct {
// 	Status         int32 `json:"status"`
// 	Message       string `json:"message"`
// 	MatchSetting *MatchSetting `json:"matchSetting"`
// }
