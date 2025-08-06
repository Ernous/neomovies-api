package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Подключение к MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("neomovies")
	collection := db.Collection("users")

	// Находим всех пользователей без поля updatedAt
	filter := bson.M{
		"$or": []bson.M{
			{"updatedAt": bson.M{"$exists": false}},
			{"favorites": bson.M{"$exists": false}},
			{"isAdmin": bson.M{"$exists": false}},
			{"adminVerified": bson.M{"$exists": false}},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var users []bson.M
	if err := cursor.All(context.Background(), &users); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d users to migrate\n", len(users))

	for _, user := range users {
		email := user["email"].(string)
		fmt.Printf("Migrating user: %s\n", email)

		// Подготавливаем обновления
		updates := bson.M{}

		// Добавляем updatedAt если его нет
		if _, exists := user["updatedAt"]; !exists {
			createdAt, ok := user["createdAt"].(time.Time)
			if !ok {
				createdAt = time.Now()
			}
			updates["updatedAt"] = createdAt
		}

		// Добавляем favorites если его нет
		if _, exists := user["favorites"]; !exists {
			updates["favorites"] = []string{}
		}

		// Добавляем isAdmin если его нет
		if _, exists := user["isAdmin"]; !exists {
			updates["isAdmin"] = false
		}

		// Добавляем adminVerified если его нет
		if _, exists := user["adminVerified"]; !exists {
			updates["adminVerified"] = false
		}

		// Добавляем avatar если его нет
		if _, exists := user["avatar"]; !exists {
			updates["avatar"] = ""
		}

		// Обновляем пользователя
		if len(updates) > 0 {
			_, err := collection.UpdateOne(
				context.Background(),
				bson.M{"email": email},
				bson.M{"$set": updates},
			)
			if err != nil {
				fmt.Printf("Error updating user %s: %v\n", email, err)
			} else {
				fmt.Printf("Successfully migrated user: %s\n", email)
			}
		}
	}

	fmt.Println("Migration completed!")
}