package models

import "time"

type MatchResult struct {
	RoomId          int64           `firestore:"roomId" json:"roomId"`
	PointList       []interface{} `firestore:"pointList" json:"pointList"`//[]models.PointOfPersonにするとfirebaseから受取時、型判定でpanicになるので一旦。
	CreateTimestamp time.Time       `firestore:"createTimestamp" json:"createTimestamp"`
	UpdateTimestamp time.Time       `firestore:"updateTimestamp" json:"updateTimestamp"`
}
