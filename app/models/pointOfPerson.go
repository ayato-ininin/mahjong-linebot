package models

type PointOfPerson struct {
	NameIndex  int64 `firestore:"nameIndex" json:"nameIndex"`
	Point      int64 `firestore:"point" json:"point"`
	IsYakitori bool  `firestore:"isYakiroti" json:"isYakiroti"`
	IsTobishou bool  `firestore:"isTobishou" json:"isTobishou"`
}
