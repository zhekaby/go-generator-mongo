package tests

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

var ctx = context.TODO()
var repo = NewUserRepositoryDefault(ctx)

func TestUsersRepository_Ping(t *testing.T) {
	Convey(t.Name(), t, func() {
		So(repo.Ping(), ShouldBeNil)
	})
}

func TestUsersRepository_InsertOne(t *testing.T) {
	u := genUser()
	Convey(t.Name(), t, func() {
		id, err := repo.InsertOne(ctx, u)
		So(err, ShouldBeNil)
		So(id, ShouldNotBeEmpty)
		Convey("FindOneById", func() {
			user, err := repo.FindOneById(ctx, id.Hex())
			So(err, ShouldBeNil)
			shouldEqual(u, user)
		})
	})
}

func TestUsersRepository_DeleteOne(t *testing.T) {
	u := genUser()
	Convey(t.Name(), t, func() {
		id, err := repo.InsertOne(ctx, u)
		So(err, ShouldBeNil)
		Convey("DeleteOneById", func() {
			isDeleted, err := repo.DeleteOneById(ctx, id.Hex())
			So(err, ShouldBeNil)
			So(isDeleted, ShouldBeTrue)

			Convey("FindOne", func() {
				user, err := repo.FindOne(ctx, bson.M{"_id": id})
				So(user, ShouldBeNil)
				So(err, ShouldBeNil)
			})

			Convey("FindOneById", func() {
				user, err := repo.FindOneById(ctx, id.Hex())
				So(user, ShouldBeNil)
				So(err, ShouldBeNil)
			})

			Convey("FindMany", func() {
				users, err := repo.FindMany(ctx, bson.M{"_id": id}, bson.D{}, 0, 100)
				So(users, ShouldHaveLength, 0)
				So(err, ShouldBeNil)
			})
		})
	})
}

func Test_Find(t *testing.T) {
	u := genUser()
	Convey(t.Name(), t, func() {
		id, err := repo.InsertOne(ctx, u)
		So(err, ShouldBeNil)

		Convey("Find*", func() {
			Convey("FindOne", func() {
				user, err := repo.FindOne(ctx, bson.M{"_id": id})
				So(err, ShouldBeNil)
				shouldEqual(u, user)
			})
			Convey("FindOneById", func() {
				user, err := repo.FindOneById(ctx, id.Hex())
				So(err, ShouldBeNil)
				shouldEqual(u, user)
			})
			Convey("FindMany", func() {
				users, err := repo.FindMany(ctx, bson.M{"_id": id}, bson.D{}, 0, 100)
				So(err, ShouldBeNil)
				So(users, ShouldHaveLength, 1)
				shouldEqual(u, users[0])
			})
		})
	})
}

func Test_Update(t *testing.T) {
	u := genUser()
	Convey(t.Name(), t, func() {
		id, err := repo.InsertOne(ctx, u)
		So(err, ShouldBeNil)
		Convey("Update*", func() {
			Convey("UpdateOneById", func() {
				matched, modified, err := repo.UpdateOneById(ctx, id.Hex(), bson.M{"$set": bson.M{"email": "m1"}})
				So(err, ShouldBeNil)
				So(matched, ShouldBeTrue)
				So(modified, ShouldBeTrue)
				u, _ := repo.FindOneById(ctx, id.Hex())
				So(u.Email, ShouldEqual, "m1")
			})
			Convey("UpdateOne", func() {
				matched, modified, err := repo.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"email": "m2"}})
				So(err, ShouldBeNil)
				So(matched, ShouldBeTrue)
				So(modified, ShouldBeTrue)
				u, _ := repo.FindOneById(ctx, id.Hex())
				So(u.Email, ShouldEqual, "m2")
			})
			Convey("UpdateOneFluent", func() {
				updater := NewUserUpdater()
				updater.SetEmail("m55").SetFinIncome(11).SetAddressCity("any")
				updater.SetProfile(*&Profile{
					FirstName: "fn1",
					LastName:  "fn2",
				})

				matched, modified, err := repo.UpdateOneFluent(ctx, bson.M{"_id": id}, updater)
				So(err, ShouldBeNil)
				So(matched, ShouldBeTrue)
				So(modified, ShouldBeTrue)
				u, _ := repo.FindOneById(ctx, id.Hex())
				So(u.Email, ShouldEqual, "m55")
				So(u.Fin.Income, ShouldEqual, 11)
				So(u.Address.City, ShouldEqual, "any")
				So(u.Profile.FirstName, ShouldEqual, "fn1")
				So(u.Profile.LastName, ShouldEqual, "fn2")
			})
		})
	})
}

func genUser() *User {
	return &User{
		Email: "test@email",
		Profile: Profile{
			FirstName: "first_name",
			LastName:  "last_name",
		},
		Address: struct {
			City string
		}{
			City: "warsaw",
		},
		Fin: &Finance{
			Income: 42,
		},
	}
}

func shouldEqual(u1, u2 *User) {
	So(u1, ShouldNotBeEmpty)
	So(u1.Email, ShouldEqual, u2.Email)
	So(u1.Profile.FirstName, ShouldEqual, u2.Profile.FirstName)
	So(u1.Profile.LastName, ShouldEqual, u2.Profile.LastName)
	So(u1.Fin.Income, ShouldEqual, u2.Fin.Income)
}
