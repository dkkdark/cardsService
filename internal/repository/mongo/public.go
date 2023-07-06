package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *ServiceImpl) FindPurchaseByCategory(category string) ([]Food, error) {
	var idStruct IDStruct
	foodCollection := s.client.Database("purchases").Collection("food")
	categoryCollection := s.client.Database("purchases").Collection("category")

	opts := options.FindOne().SetProjection(bson.D{{"_id", 1}})
	err := categoryCollection.FindOne(s.c, bson.D{{"name", category}}, opts).Decode(&idStruct)
	res, err := foodCollection.Find(s.c, bson.D{{"category", idStruct.ID}})
	var food []Food
	err = res.All(s.c, &food)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *ServiceImpl) FindPurchaseByName(name string) ([]Food, error) {
	foodCollection := s.client.Database("purchases").Collection("food")

	res, err := foodCollection.Find(s.c, bson.D{{"name", name}})
	var food []Food
	err = res.All(s.c, &food)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *ServiceImpl) FilterPurchasesByCost(params *CostParams) ([]Food, error) {
	foodCollection := s.client.Database("purchases").Collection("food")

	res, err := foodCollection.Find(s.c, bson.D{{"cost", bson.D{{"$gt", params.LowCost}, {"$lt", params.HighCost}}}})
	var food []Food
	err = res.All(s.c, &food)
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (s *ServiceImpl) AddPurchaseToCart(params *AddToCartParams) error {
	clientId, err := primitive.ObjectIDFromHex(params.UserID)
	fId, err := primitive.ObjectIDFromHex(params.PurchaseID)

	filter := bson.D{{"_id", clientId}}
	update := bson.D{
		{"$push",
			bson.D{{"cart", fId}}},
	}

	clientCollection := s.client.Database("purchases").Collection("clients")

	_, err = clientCollection.UpdateOne(s.c, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServiceImpl) CountCartCost(id string) (*CartCount, error) {
	clientId, err := primitive.ObjectIDFromHex(id)

	matchStage := bson.D{{"$match", bson.D{{"_id", clientId}}}}
	unwindStage1 := bson.D{{"$unwind", bson.D{{"path", "$cart"}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "cart"}, {"foreignField", "_id"}, {"as", "food"}}}}
	unwindStage2 := bson.D{{"$unwind", bson.D{{"path", "$food"}}}}
	projectStage1 := bson.D{{"$project", bson.D{{"cost", "$food.cost"}}}}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", nil},
			{"sum", bson.D{{"$sum", "$cost"}}},
		}}}
	projectStage2 := bson.D{{"$project", bson.D{{"_id", 0}}}}

	clientCollection := s.client.Database("purchases").Collection("clients")

	res, err := clientCollection.Aggregate(s.c, mongo.Pipeline{matchStage, unwindStage1, lookupStage, unwindStage2, projectStage1, groupStage, projectStage2})
	if err != nil {
		return nil, err
	}
	var cartCount []CartCount
	err = res.All(s.c, &cartCount)

	return &CartCount{
		cartCount[0].Sum,
	}, nil
}
