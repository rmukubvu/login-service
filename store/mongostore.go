package store

import (
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"time"
)

func init() {
	var err error
	client,err = mongo.NewClient(databaseUri)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := getContext()
	err = client.Connect(ctx)
	if err != nil {
		closeDb()
		log.Fatal(err)
	}
	db = client.Database(schemaName)
}

func (u *User) InsertRecord() error {
	collection := db.Collection(loginCollectionName)
	ctx, _ := getContext()
	//hash the password
	u.Password = getHashedPassword(u.Password)
	//then insert
	_, err := collection.InsertOne(ctx, u)
	return err
}

func FetchRecord(userName string) (User,error) {
	collection := db.Collection(loginCollectionName)
	singleResult := collection.FindOne(context.Background(), User{UserName: userName})
	var user User
	err := singleResult.Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user,nil
}

func (u *User) UpdateRecord() error {
	collection := db.Collection(loginCollectionName)
	//hash the password
	u.Password = getHashedPassword(u.Password)
	//update
	_, err := collection.ReplaceOne(context.Background(), User{UserName: u.UserName}, u)
	return err
}

func fetchAllRecords() []User {
	var users []User
	collection := db.Collection(loginCollectionName)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil
	}
	for cursor.Next(context.Background()) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	if err := cursor.Err(); err == nil {
		return nil
	}
	return users
}

func closeDb() {
	if client != nil {
		ctx, _ := getContext()
		client.Disconnect(ctx)
	}
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(),10 * time.Second)
}
