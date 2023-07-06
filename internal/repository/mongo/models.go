package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type IDStruct struct {
	ID primitive.ObjectID `bson:"_id"`
}

type CartCount struct {
	Sum int `bson:"sum"`
}

type Food struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Name  string             `bson:"name"`
	Count float32            `bson:"count"`
	Cost  float32            `bson:"cost"`
	Type  string             `bson:"type"`
}

type CartSum struct {
	Sum int
}
