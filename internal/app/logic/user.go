package logic

import (
	"errors"

	"github.com/technoshantoms/mccs-alpha-api/internal/app/repository/es"
	"github.com/technoshantoms/mccs-alpha-api/internal/app/repository/mongo"
	"github.com/technoshantoms/mccs-alpha-api/internal/app/repository/redis"
	"github.com/technoshantoms/mccs-alpha-api/internal/app/types"
	"github.com/technoshantoms/mccs-alpha-api/util"
	"github.com/technoshantoms/mccs-alpha-api/util/bcrypt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct{}

var User = &user{}

func (u *user) Create(user *types.User) (*types.User, error) {
	_, err := mongo.User.FindByEmail(user.Email)
	if err == nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	created, err := mongo.User.Create(user)
	if err != nil {
		return nil, err
	}

	err = es.User.Create(created.ID, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// POST /signup

func (u *user) AssociateEntity(userID, entityID primitive.ObjectID) error {
	err := mongo.User.AssociateEntity([]primitive.ObjectID{userID}, entityID)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) EmailExists(email string) bool {
	_, err := mongo.User.FindByEmail(email)
	if err != nil {
		return false
	}
	return true
}

// POST /login

func (u *user) Login(email string, password string) (*types.User, error) {
	user, err := mongo.User.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	attempts := redis.GetLoginAttempts(email)

	if attempts >= viper.GetInt("login_attempts.limit") {
		return nil, ErrLoginLocked
	}

	err = bcrypt.CompareHash(user.Password, password)
	if err != nil {
		if attempts+1 >= viper.GetInt("login_attempts.limit") {
			return nil, ErrLoginLocked
		}
		return nil, errors.New("Invalid password.")
	}

	redis.ResetLoginAttempts(email)

	return user, nil
}

func (u *user) FindByID(id primitive.ObjectID) (*types.User, error) {
	user, err := mongo.User.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *user) FindByStringID(id string) (*types.User, error) {
	objectID, _ := primitive.ObjectIDFromHex(id)
	user, err := mongo.User.FindByID(objectID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *user) FindByIDs(ids []primitive.ObjectID) ([]*types.User, error) {
	users, err := mongo.User.FindByIDs(ids)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *user) FindByEmail(email string) (*types.User, error) {
	user, err := mongo.User.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *user) FindByEntityID(id primitive.ObjectID) (*types.User, error) {
	user, err := mongo.User.FindByEntityID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *user) UpdateLoginInfo(id primitive.ObjectID, ip string) (*types.LoginInfo, error) {
	info, err := mongo.User.UpdateLoginInfo(id, ip)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (u *user) IncLoginAttempts(email string) error {
	err := redis.IncLoginAttempts(email)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) ResetPassword(email string, newPassword string) error {
	user, err := mongo.User.FindByEmail(email)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.Hash(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	err = mongo.User.UpdatePassword(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *user) FindOneAndUpdate(userID primitive.ObjectID, update *types.User) (*types.User, error) {
	err := es.User.Update(userID, update)
	if err != nil {
		return nil, err
	}
	updated, err := mongo.User.FindOneAndUpdate(userID, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (u *user) FindEntities(userID primitive.ObjectID) ([]*types.Entity, error) {
	user, err := mongo.User.FindByID(userID)
	if err != nil {
		return nil, err
	}
	entities, err := mongo.Entity.FindByStringIDs(util.ToIDStrings(user.Entities))
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (u *user) AdminFindOneAndUpdate(req *types.AdminUpdateUserReq) (*types.User, error) {
	err := es.User.AdminUpdate(req)
	if err != nil {
		return nil, err
	}
	updated, err := mongo.User.AdminFindOneAndUpdate(req)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// DELETE /admin/users/{userID}

func (u *user) AdminFindOneAndDelete(id primitive.ObjectID) (*types.User, error) {
	err := es.User.Delete(id.Hex())
	if err != nil {
		return nil, err
	}
	deleted, err := mongo.User.AdminFindOneAndDelete(id)
	if err != nil {
		return nil, err
	}
	return deleted, nil
}

func (u *user) AdminSearchUser(req *types.AdminSearchUserReq) (*types.SearchUserResult, error) {
	result, err := es.User.AdminSearchUser(req)
	if err != nil {
		return nil, err
	}
	users, err := mongo.User.FindByStringIDs(result.UserIDs)
	if err != nil {
		return nil, err
	}
	return &types.SearchUserResult{
		Users:           users,
		NumberOfResults: result.NumberOfResults,
		TotalPages:      result.TotalPages,
	}, nil
}
