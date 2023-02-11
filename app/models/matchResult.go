package models

import "time"

type MatchResult struct {
	RoomId          int64           `firestore:"roomId" json:"roomId"`
	MatchIndex      int64           `firestore:"matchIndex" json:"matchIndex"`
	PointList       []PointOfPerson `firestore:"pointList" json:"pointList"`
	CreateTimestamp time.Time       `firestore:"createTimestamp" json:"createTimestamp"`
	UpdateTimestamp time.Time       `firestore:"updateTimestamp" json:"updateTimestamp"`
}
