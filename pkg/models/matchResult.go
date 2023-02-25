package models

import "time"

type MatchResult struct {
	DocId           string          `firestore:"docId" json:"docId"`
	RoomId          int64           `firestore:"roomId" json:"roomId"`
	PointList       []PointOfPerson `firestore:"pointList" json:"pointList"`
	CreateTimestamp time.Time       `firestore:"createTimestamp" json:"createTimestamp"`
	UpdateTimestamp time.Time       `firestore:"updateTimestamp" json:"updateTimestamp"`
}
