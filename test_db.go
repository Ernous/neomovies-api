package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("=== DATABASE CONNECTION TEST ===")
	
	// Подключаемся к базе данных
	uri := "mongodb+srv://neomoviesmail:Vfhreif1@neo-movies.nz1e2.mongodb.net/neomovies"
	
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer client.Disconnect(ctx)
	
	// Проверяем подключение
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	fmt.Printf("✅ Database connection successful\n")
	
	// Получаем базу данных
	db := client.Database("neomovies")
	fmt.Printf("📊 Database name: %s\n", db.Name())
	
	// Получаем список всех коллекций
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to list collections:", err)
	}
	
	fmt.Printf("📁 Available collections: %v\n", collections)
	
	// Проверяем коллекцию users
	collection := db.Collection("users")
	
	// Подсчитываем количество документов
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Fatal("Failed to count users:", err)
	}
	
	fmt.Printf("👥 Total users in database: %d\n", count)
	
	if count > 0 {
		// Показываем всех пользователей
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal("Failed to find users:", err)
		}
		defer cursor.Close(ctx)
		
		var users []bson.M
		if err := cursor.All(ctx, &users); err != nil {
			log.Fatal("Failed to decode users:", err)
		}
		
		fmt.Printf("📋 All users in database:\n")
		for i, user := range users {
			fmt.Printf("  %d. Email: %s, Name: %s, Verified: %v\n", 
				i+1, 
				user["email"], 
				user["name"], 
				user["verified"])
		}
		
		// Тестируем поиск конкретного пользователя
		fmt.Printf("\n🔍 Testing specific user search:\n")
		testEmails := []string{"neo.movies.mail@gmail.com", "fenixoffc@gmail.com", "test@example.com"}
		
		for _, email := range testEmails {
			var user bson.M
			err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
			if err != nil {
				fmt.Printf("  ❌ User %s: NOT FOUND (%v)\n", email, err)
			} else {
				fmt.Printf("  ✅ User %s: FOUND (Name: %s, Verified: %v)\n", 
					email, 
					user["name"], 
					user["verified"])
			}
		}
	}
	
	fmt.Println("=== END DATABASE TEST ===")
}