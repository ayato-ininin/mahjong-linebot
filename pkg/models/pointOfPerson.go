package models

type PointOfPerson struct {
	NameIndex  int64 `firestore:"nameIndex" json:"nameIndex"`
	Point      int64 `firestore:"point" json:"point"`
	IsYakitori bool  `firestore:"isYakitori" json:"isYakitori"`
	IsTobishou bool  `firestore:"isTobishou" json:"isTobishou"`
}
