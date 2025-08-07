package main

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	fmt.Println("=== CREATING TEST USER ===")
	
	// Подключаемся к базе данных
	uri := "mongodb+srv://neomoviesmail:Vfhreif1@neo-movies.nz1e2.mongodb.net/database"
	
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
	db := client.Database("database")
	collection := db.Collection("users")
	
	// Проверяем, существует ли уже тестовый пользователь
	var existingUser bson.M
	err = collection.FindOne(ctx, bson.M{"email": "test@example.com"}).Decode(&existingUser)
	if err == nil {
		fmt.Printf("ℹ️ Test user already exists\n")
		return
	}
	
	// Создаем хеш пароля
	password := "test123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}
	
	// Создаем тестового пользователя
	testUser := bson.M{
		"email":     "test@example.com",
		"password":  string(hashedPassword),
		"name":      "Test User",
		"verified":  true,
		"favorites": []string{},
		"isAdmin":   false,
		"adminVerified": false,
		"createdAt": bson.M{"$date": "2024-12-21T09:37:04.363Z"},
		"updatedAt": bson.M{"$date": "2024-12-21T09:37:04.363Z"},
	}
	
	_, err = collection.InsertOne(ctx, testUser)
	if err != nil {
		log.Fatal("Failed to create test user:", err)
	}
	
	fmt.Printf("✅ Test user created successfully\n")
	fmt.Printf("📧 Email: test@example.com\n")
	fmt.Printf("🔑 Password: %s\n", password)
	fmt.Printf("=== END CREATING TEST USER ===\n")
}