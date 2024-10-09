package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/technoshantoms/mccs-alpha-api/global/constant"
	"github.com/technoshantoms/mccs-alpha-api/internal/app/types"
	"github.com/technoshantoms/mccs-alpha-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type entity struct {
	c *mongo.Collection
}

var Entity = &entity{}

func (en *entity) Register(db *mongo.Database) {
	en.c = db.Collection("entities")
}

func (e *entity) Create(update *types.Entity) (*types.Entity, error) {
	filter := bson.M{"_id": bson.M{"$exists": false}}
	update.Status = constant.Entity.Pending
	update.CreatedAt = time.Now()

	result := e.c.FindOneAndUpdate(
		context.Background(),
		filter,
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	if result.Err() != nil {
		return nil, result.Err()
	}

	created := types.Entity{}
	err := result.Decode(&created)
	if err != nil {
		return nil, err
	}

	// Make sure "offers" and "wants" fields exist so it's much easier to update later on.
	err = e.setDefaultOffersAndWants(created.ID, created.Offers, created.Wants)

	return &created, nil
}

func (e *entity) setDefaultOffersAndWants(entityID primitive.ObjectID, offers []*types.TagField, wants []*types.TagField) error {
	if len(offers) == 0 || len(wants) == 0 {
		filter := bson.M{"_id": entityID}
		update := bson.M{}
		if len(offers) == 0 {
			update["offers"] = []*types.TagField{}
		}
		if len(wants) == 0 {
			update["wants"] = []*types.TagField{}
		}
		_, err := e.c.UpdateOne(
			context.Background(),
			filter,
			bson.M{"$set": update},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *entity) FindByID(id primitive.ObjectID) (*types.Entity, error) {
	ctx := context.Background()
	entity := types.Entity{}
	filter := bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}
	err := e.c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Please enter valid entity id.")
		}
		return nil, err
	}
	return &entity, nil
}

func (e *entity) FindByAccountNumber(accountNumber string) (*types.Entity, error) {
	ctx := context.Background()
	entity := types.Entity{}
	filter := bson.M{
		"accountNumber": accountNumber,
		"deletedAt":     bson.M{"$exists": false},
	}
	err := e.c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Entity not found.")
		}
		return nil, err
	}
	return &entity, nil
}

func (e *entity) FindByEmail(email string) (*types.Entity, error) {
	if email == "" {
		return nil, errors.New("Please specify an email address.")
	}
	email = strings.ToLower(email)
	entity := types.Entity{}
	filter := bson.M{"email": email, "deletedAt": bson.M{"$exists": false}}
	err := e.c.FindOne(context.Background(), filter).Decode(&Entity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("The specified entity could not be found.")
		}
		return nil, err
	}
	return &entity, nil
}

// PATCH /user/entities/{entityID}

func (e *entity) FindOneAndUpdate(req *types.UpdateUserEntityReq) (*types.Entity, error) {
	updates := []bson.M{}

	update := bson.M{"updatedAt": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Telephone != "" {
		update["telephone"] = req.Telephone
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.IncType != "" {
		update["incType"] = req.IncType
	}
	if req.CompanyNumber != "" {
		update["companyNumber"] = req.CompanyNumber
	}
	if req.Website != "" {
		update["website"] = req.Website
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.DeclaredTurnover != nil {
		update["declaredTurnover"] = *req.DeclaredTurnover
	}
	if req.Address != "" {
		update["address"] = req.Address
	}
	if req.City != "" {
		update["city"] = req.City
	}
	if req.Region != "" {
		update["region"] = req.Region
	}
	if req.Country != "" {
		update["country"] = req.Country
	}
	if req.PostalCode != "" {
		update["postalCode"] = req.PostalCode
	}
	if req.ReceiveDailyMatchNotificationEmail != nil {
		update["receiveDailyMatchNotificationEmail"] = *req.ReceiveDailyMatchNotificationEmail
	}
	if req.ShowTagsMatchedSinceLastLogin != nil {
		update["showTagsMatchedSinceLastLogin"] = *req.ShowTagsMatchedSinceLastLogin
	}
	updates = append(updates, bson.M{"$set": update})

	push := bson.M{}
	if len(req.AddedOffers) != 0 {
		push["offers"] = bson.M{"$each": types.ToTagFields(req.AddedOffers)}
	}
	if len(req.AddedWants) != 0 {
		push["wants"] = bson.M{"$each": types.ToTagFields(req.AddedWants)}
	}
	if len(push) != 0 {
		updates = append(updates, bson.M{"$push": push})
	}

	pull := bson.M{}
	if len(req.RemovedOffers) != 0 {
		pull["offers"] = bson.M{"name": bson.M{"$in": req.RemovedOffers}}
	}
	if len(req.RemovedWants) != 0 {
		pull["wants"] = bson.M{"name": bson.M{"$in": req.RemovedWants}}
	}
	if len(pull) != 0 {
		updates = append(updates, bson.M{"$pull": pull})
	}

	filter := bson.M{"_id": req.OriginEntity.ID, "deletedAt": bson.M{"$exists": false}}
	var writes []mongo.WriteModel
	for _, update := range updates {
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		writes = append(writes, model)
	}

	_, err := e.c.BulkWrite(context.Background(), writes)
	if err != nil {
		return nil, err
	}

	entity, err := e.FindByID(req.OriginEntity.ID)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// PATCH /admin/entities/{entityID}

func (e *entity) AdminFindOneAndUpdate(req *types.AdminUpdateEntityReq) (*types.Entity, error) {
	updates := []bson.M{}

	update := bson.M{"updatedAt": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Telephone != "" {
		update["telephone"] = req.Telephone
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.IncType != "" {
		update["incType"] = req.IncType
	}
	if req.CompanyNumber != "" {
		update["companyNumber"] = req.CompanyNumber
	}
	if req.Website != "" {
		update["website"] = req.Website
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.DeclaredTurnover != nil {
		update["declaredTurnover"] = *req.DeclaredTurnover
	}
	if req.Address != "" {
		update["address"] = req.Address
	}
	if req.City != "" {
		update["city"] = req.City
	}
	if req.Region != "" {
		update["region"] = req.Region
	}
	if req.Country != "" {
		update["country"] = req.Country
	}
	if req.PostalCode != "" {
		update["postalCode"] = req.PostalCode
	}
	if req.Categories != nil {
		update["categories"] = util.FormatTags(*req.Categories)
	}
	if req.Users != nil {
		update["users"] = util.ToObjectIDs(*req.Users)
	}
	if req.Status != "" {
		update["status"] = req.Status
	}
	if req.ReceiveDailyMatchNotificationEmail != nil {
		update["receiveDailyMatchNotificationEmail"] = *req.ReceiveDailyMatchNotificationEmail
	}
	if req.ShowTagsMatchedSinceLastLogin != nil {
		update["showTagsMatchedSinceLastLogin"] = *req.ShowTagsMatchedSinceLastLogin
	}
	// TODO
	// This is a trick to prevent setting nothing for the entity.
	// If we don't do this then it will throw this error:
	// (FailedToParse) '$set' is empty. You must specify a field like so: {$set: {<field>: ...}}
	if req.Name == "" {
		update["name"] = req.OriginEntity.Name
	}
	updates = append(updates, bson.M{"$set": update})

	push := bson.M{}
	if len(req.AddedOffers) != 0 {
		push["offers"] = bson.M{"$each": types.ToTagFields(req.AddedOffers)}
	}
	if len(req.AddedWants) != 0 {
		push["wants"] = bson.M{"$each": types.ToTagFields(req.AddedWants)}
	}
	if len(push) != 0 {
		updates = append(updates, bson.M{"$push": push})
	}

	pull := bson.M{}
	if len(req.RemovedOffers) != 0 {
		pull["offers"] = bson.M{"name": bson.M{"$in": req.RemovedOffers}}
	}
	if len(req.RemovedWants) != 0 {
		pull["wants"] = bson.M{"name": bson.M{"$in": req.RemovedWants}}
	}
	if len(pull) != 0 {
		updates = append(updates, bson.M{"$pull": pull})
	}

	filter := bson.M{"_id": req.OriginEntity.ID, "deletedAt": bson.M{"$exists": false}}
	var writes []mongo.WriteModel
	for _, update := range updates {
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		writes = append(writes, model)
	}

	_, err := e.c.BulkWrite(context.Background(), writes)
	if err != nil {
		return nil, err
	}

	err = User.AssociateEntity(req.AddedUsers, req.OriginEntity.ID)
	if err != nil {
		return nil, err
	}
	err = User.removeAssociatedEntity(req.RemovedUsers, req.OriginEntity.ID)
	if err != nil {
		return nil, err
	}

	entity, err := e.FindByID(req.OriginEntity.ID)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// DELETE /admin/entities/{entityID}

func (e *entity) AdminFindOneAndDelete(id primitive.ObjectID) (*types.Entity, error) {
	filter := bson.M{"_id": id, "deletedAt": bson.M{"$exists": false}}

	result := e.c.FindOneAndUpdate(
		context.Background(),
		filter,
		bson.M{"$set": bson.M{
			"users":     []primitive.ObjectID{},
			"deletedAt": time.Now(),
			"updatedAt": time.Now(),
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	entity := types.Entity{}
	err := result.Decode(&entity)
	if err != nil {
		return nil, result.Err()
	}

	err = User.removeAssociatedEntity(entity.Users, entity.ID)
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (e *entity) AddToFavoriteEntities(req *types.AddToFavoriteReq) error {
	addToEntityID, _ := primitive.ObjectIDFromHex(req.AddToEntityID)
	favoriteEntityID, _ := primitive.ObjectIDFromHex(req.FavoriteEntityID)

	filter := bson.M{"_id": addToEntityID}
	update := bson.M{}
	if *req.Favorite {
		update["$addToSet"] = bson.M{"favoriteEntities": favoriteEntityID}
	} else {
		update["$pull"] = bson.M{"favoriteEntities": favoriteEntityID}
	}
	_, err := e.c.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) FindByStringIDs(ids []string) ([]*types.Entity, error) {
	var results []*types.Entity

	objectIDs, err := toObjectIDs(ids)
	if err != nil {
		return nil, err
	}

	pipeline := newFindByIDsPipeline(objectIDs)
	cur, err := e.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem types.Entity
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return results, nil
}

func (e *entity) FindByIDs(ids []primitive.ObjectID) ([]*types.Entity, error) {
	var results []*types.Entity

	pipeline := newFindByIDsPipeline(ids)
	cur, err := e.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem types.Entity
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return results, nil
}

// PATCH /admin/entities/{entityID}

func (e *entity) UpdateAllTagsCreatedAt(id primitive.ObjectID, t time.Time) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"offers.$[].createdAt": t,
		"wants.$[].createdAt":  t,
	}}
	_, err := e.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) SetMemberStartedAt(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"memberStartedAt": time.Now(),
	}}
	_, err := e.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) RenameCategory(old string, new string) error {
	// Push the new tag tag name.
	filter := bson.M{"categories": old}
	update := bson.M{
		"$push": bson.M{"categories": new},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err := e.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	// Delete the old tag name.
	filter = bson.M{"categories": old}
	update = bson.M{
		"$pull": bson.M{"categories": old},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err = e.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) DeleteCategory(name string) error {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"categories": name},
		},
	}
	update := bson.M{
		"$pull": bson.M{
			"categories": name,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	_, err := e.c.UpdateMany(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) RenameTag(old string, new string) error {
	err := e.updateOffers(old, new)
	if err != nil {
		return err
	}
	err = e.updateWants(old, new)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) DeleteTag(name string) error {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"offers.name": name},
			bson.M{"wants.name": name},
		},
	}
	update := bson.M{
		"$pull": bson.M{
			"offers": bson.M{"name": name},
			"wants":  bson.M{"name": name},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	_, err := e.c.UpdateMany(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) DeleteByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"deletedAt": time.Now()}}
	_, err := e.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) updateOffers(old string, new string) error {
	filter := bson.M{"offers.name": old}
	update := bson.M{
		"$set": bson.M{
			"offers.$.name": new,
			"updatedAt":     time.Now(),
		},
	}
	_, err := e.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (e *entity) updateWants(old string, new string) error {
	filter := bson.M{"wants.name": old}
	update := bson.M{
		"$set": bson.M{
			"wants.$.name": new,
			"updatedAt":    time.Now(),
		},
	}
	_, err := e.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

// PATCH /admin/users/{userID}

func (en *entity) AssociateUser(entityIDs []primitive.ObjectID, UserID primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": entityIDs}}
	updates := []bson.M{
		bson.M{"$addToSet": bson.M{"users": UserID}},
		bson.M{"$set": bson.M{"updatedAt": time.Now()}},
	}

	var writes []mongo.WriteModel
	for _, update := range updates {
		model := mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update)
		writes = append(writes, model)
	}

	_, err := en.c.BulkWrite(context.Background(), writes)
	if err != nil {
		return err
	}
	return nil
}

// PATCH /admin/users/{userID}

func (e *entity) removeAssociatedUser(entityIDs []primitive.ObjectID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": entityIDs}}
	updates := []bson.M{
		bson.M{"$pull": bson.M{"users": userID}},
		bson.M{"$set": bson.M{"updatedAt": time.Now()}},
	}

	var writes []mongo.WriteModel
	for _, update := range updates {
		model := mongo.NewUpdateManyModel().SetFilter(filter).SetUpdate(update)
		writes = append(writes, model)
	}

	_, err := e.c.BulkWrite(context.Background(), writes)
	if err != nil {
		return err
	}
	return nil
}

// daily_email_schedule

func (e *entity) FindByDailyNotification() ([]*types.Entity, error) {
	filter := bson.M{
		"receiveDailyMatchNotificationEmail": true,
		"deletedAt":                          bson.M{"$exists": false},
	}
	cur, err := e.c.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var entities []*types.Entity
	for cur.Next(context.TODO()) {
		var elem types.Entity
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.TODO())

	return entities, nil
}

// daily_email_schedule

func (e *entity) UpdateLastNotificationSentDate(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"lastNotificationSentDate": time.Now(),
		"updatedAt":                time.Now(),
	}}
	_, err := e.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return err
	}
	return nil
}
