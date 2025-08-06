package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"neomovies-api/pkg/models"
)

type AuthService struct {
	db           *mongo.Database
	jwtSecret    string
	emailService *EmailService
}

func NewAuthService(db *mongo.Database, jwtSecret string, emailService *EmailService) *AuthService {
	service := &AuthService{
		db:           db,
		jwtSecret:    jwtSecret,
		emailService: emailService,
	}
	
	// Проверяем подключение к базе данных при инициализации
	go service.checkDatabaseConnection()
	
	return service
}

// checkDatabaseConnection проверяет подключение к базе данных и выводит диагностическую информацию
func (s *AuthService) checkDatabaseConnection() {
	ctx := context.Background()
	
	// Проверяем подключение
	err := s.db.Client().Ping(ctx, nil)
	if err != nil {
		fmt.Printf("ERROR: Database connection failed: %v\n", err)
		return
	}
	
	fmt.Printf("INFO: Database connection successful\n")
	fmt.Printf("INFO: Database name: %s\n", s.db.Name())
	
	// Получаем список всех коллекций
	collections, err := s.db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		fmt.Printf("ERROR: Failed to list collections: %v\n", err)
		return
	}
	
	fmt.Printf("INFO: Available collections: %v\n", collections)
	
	// Проверяем коллекцию users
	collection := s.db.Collection("users")
	
	// Подсчитываем количество документов
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Printf("ERROR: Failed to count users: %v\n", err)
		return
	}
	
	fmt.Printf("INFO: Total users in database: %d\n", count)
	
	// Если пользователей нет, создаем тестового пользователя
	if count == 0 {
		fmt.Printf("INFO: No users found, creating test user...\n")
		s.createTestUser()
	} else {
		// Показываем первых несколько пользователей
		cursor, err := collection.Find(ctx, bson.M{})
		if err == nil {
			defer cursor.Close(ctx)
			var users []bson.M
			if err := cursor.All(ctx, &users); err == nil {
				fmt.Printf("INFO: First %d users in database:\n", len(users))
				for i, user := range users {
					if i < 3 { // Показываем только первые 3
						fmt.Printf("  User %d: email=%s, name=%s\n", i+1, user["email"], user["name"])
					}
				}
			}
		}
	}
}

// createTestUser создает тестового пользователя для диагностики
func (s *AuthService) createTestUser() {
	collection := s.db.Collection("users")
	
	// Проверяем, существует ли уже тестовый пользователь
	var existingUser bson.M
	err := collection.FindOne(context.Background(), bson.M{"email": "test@example.com"}).Decode(&existingUser)
	if err == nil {
		fmt.Printf("INFO: Test user already exists\n")
		return
	}
	
	// Создаем тестового пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("ERROR: Failed to hash password: %v\n", err)
		return
	}
	
	testUser := bson.M{
		"email":     "test@example.com",
		"password":  string(hashedPassword),
		"name":      "Test User",
		"verified":  true,
		"favorites": []string{},
		"isAdmin":   false,
		"adminVerified": false,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}
	
	_, err = collection.InsertOne(context.Background(), testUser)
	if err != nil {
		fmt.Printf("ERROR: Failed to create test user: %v\n", err)
		return
	}
	
	fmt.Printf("INFO: Test user created successfully\n")
	
	// Также создаем пользователя из примера
	s.createExampleUser()
}

// createExampleUser создает пользователя из примера
func (s *AuthService) createExampleUser() {
	collection := s.db.Collection("users")
	
	// Проверяем, существует ли уже пользователь
	var existingUser bson.M
	err := collection.FindOne(context.Background(), bson.M{"email": "fenixoffc@gmail.com"}).Decode(&existingUser)
	if err == nil {
		fmt.Printf("INFO: Example user already exists\n")
		return
	}
	
	// Создаем пользователя из примера
	hashedPassword := "$2a$12$pc9PdvyI5LFOZ9fvIbKhZ.tM7dt9YC0.RRxLIT21xR6GCrijry8Zy"
	
	exampleUser := bson.M{
		"email":     "fenixoffc@gmail.com",
		"password":  hashedPassword,
		"name":      "Foxix",
		"verified":  true,
		"favorites": []string{},
		"isAdmin":   false,
		"adminVerified": false,
		"createdAt": time.Date(2024, 12, 21, 9, 37, 4, 363000000, time.UTC),
		"updatedAt": time.Date(2024, 12, 21, 9, 37, 4, 363000000, time.UTC),
	}
	
	_, err = collection.InsertOne(context.Background(), exampleUser)
	if err != nil {
		fmt.Printf("ERROR: Failed to create example user: %v\n", err)
		return
	}
	
	fmt.Printf("INFO: Example user created successfully\n")
}

// Генерация 6-значного кода
func (s *AuthService) generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(900000)+100000)
}

func (s *AuthService) Register(req models.RegisterRequest) (map[string]interface{}, error) {
	collection := s.db.Collection("users")

	// Проверяем, не существует ли уже пользователь с таким email
	var existingUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Генерируем код верификации
	code := s.generateVerificationCode()
	codeExpires := time.Now().Add(10 * time.Minute) // 10 минут

	// Создаем нового пользователя (НЕ ВЕРИФИЦИРОВАННОГО)
	user := models.User{
		ID:                 primitive.NewObjectID(),
		Email:              req.Email,
		Password:           string(hashedPassword),
		Name:               req.Name,
		Favorites:          []string{},
		Verified:           false,
		VerificationCode:   code,
		VerificationExpires: codeExpires,
		IsAdmin:            false,
		AdminVerified:      false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	// Отправляем код верификации на email
	if s.emailService != nil {
		go s.emailService.SendVerificationEmail(user.Email, code)
	}

	return map[string]interface{}{
		"success": true,
		"message": "Registered. Check email for verification code.",
	}, nil
}

func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	collection := s.db.Collection("users")

	fmt.Printf("Attempting to find user with email: %s\n", req.Email)
	fmt.Printf("Database name: %s\n", s.db.Name())
	fmt.Printf("Collection name: %s\n", collection.Name())
	
	// Проверяем подключение к базе данных
	err := s.db.Client().Ping(context.Background(), nil)
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		return nil, errors.New("Database connection failed")
	}
	
	// Сначала попробуем найти пользователя без декодирования в структуру
	var rawUser bson.M
	err = collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&rawUser)
	if err != nil {
		fmt.Printf("Login error: user not found for email %s, error: %v\n", req.Email, err)
		
		// Попробуем найти всех пользователей с похожим email для диагностики
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err == nil {
			defer cursor.Close(context.Background())
			var allUsers []bson.M
			if err := cursor.All(context.Background(), &allUsers); err == nil {
				fmt.Printf("All users in database: %d users\n", len(allUsers))
				for i, u := range allUsers {
					if i < 5 { // Показываем только первые 5 пользователей
						fmt.Printf("User %d: %v\n", i+1, u)
					}
				}
			}
		}
		
		// Попробуем найти пользователя по частичному совпадению email
		fmt.Printf("Trying partial email match...\n")
		cursor, err := collection.Find(context.Background(), bson.M{"email": bson.M{"$regex": ".*" + req.Email + ".*", "$options": "i"}})
		if err == nil {
			defer cursor.Close(context.Background())
			var partialUsers []bson.M
			if err := cursor.All(context.Background(), &partialUsers); err == nil {
				fmt.Printf("Partial match users: %d users\n", len(partialUsers))
				for i, u := range partialUsers {
					fmt.Printf("Partial match user %d: %v\n", i+1, u)
				}
			}
		}
		
		return nil, errors.New("User not found")
	}
	
	fmt.Printf("Raw user found: %v\n", rawUser)
	
	// Создаем пользователя вручную из raw данных, чтобы избежать проблем с отсутствующими полями
	user := models.User{
		Email:    rawUser["email"].(string),
		Password: rawUser["password"].(string),
		Name:     rawUser["name"].(string),
		Verified: rawUser["verified"].(bool),
	}
	
	// Обрабатываем ID
	if id, ok := rawUser["_id"].(primitive.ObjectID); ok {
		user.ID = id
	}
	
	// Обрабатываем опциональные поля
	if favorites, ok := rawUser["favorites"].(primitive.A); ok {
		user.Favorites = make([]string, len(favorites))
		for i, fav := range favorites {
			user.Favorites[i] = fav.(string)
		}
	} else {
		user.Favorites = []string{}
	}
	
	if avatar, ok := rawUser["avatar"].(string); ok {
		user.Avatar = avatar
	}
	
	if isAdmin, ok := rawUser["isAdmin"].(bool); ok {
		user.IsAdmin = isAdmin
	}
	
	if adminVerified, ok := rawUser["adminVerified"].(bool); ok {
		user.AdminVerified = adminVerified
	}
	
	// Обрабатываем даты
	if createdAt, ok := rawUser["createdAt"].(primitive.DateTime); ok {
		user.CreatedAt = createdAt.Time()
	} else {
		user.CreatedAt = time.Now() // fallback
	}
	
	if updatedAt, ok := rawUser["updatedAt"].(primitive.DateTime); ok {
		user.UpdatedAt = updatedAt.Time()
	} else {
		user.UpdatedAt = user.CreatedAt // fallback to createdAt
	}
	
	fmt.Printf("User found: ID=%s, Email=%s, Verified=%v\n", user.ID.Hex(), user.Email, user.Verified)

	// Проверяем верификацию email
	if !user.Verified {
		fmt.Printf("Login error: email not verified for %s\n", req.Email)
		return nil, errors.New("Account not activated. Please verify your email.")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Printf("Login error: invalid password for email %s\n", req.Email)
		return nil, errors.New("Invalid password")
	}

	fmt.Printf("Login successful for user %s\n", req.Email)

	// Генерируем JWT токен
	token, err := s.generateJWT(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Сначала получаем raw данные
	var rawUser bson.M
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&rawUser)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя вручную из raw данных
	user := models.User{
		Email:    rawUser["email"].(string),
		Password: rawUser["password"].(string),
		Name:     rawUser["name"].(string),
		Verified: rawUser["verified"].(bool),
	}
	
	// Обрабатываем ID
	if id, ok := rawUser["_id"].(primitive.ObjectID); ok {
		user.ID = id
	}
	
	// Обрабатываем опциональные поля
	if favorites, ok := rawUser["favorites"].(primitive.A); ok {
		user.Favorites = make([]string, len(favorites))
		for i, fav := range favorites {
			user.Favorites[i] = fav.(string)
		}
	} else {
		user.Favorites = []string{}
	}
	
	if avatar, ok := rawUser["avatar"].(string); ok {
		user.Avatar = avatar
	}
	
	if isAdmin, ok := rawUser["isAdmin"].(bool); ok {
		user.IsAdmin = isAdmin
	}
	
	if adminVerified, ok := rawUser["adminVerified"].(bool); ok {
		user.AdminVerified = adminVerified
	}
	
	// Обрабатываем даты
	if createdAt, ok := rawUser["createdAt"].(primitive.DateTime); ok {
		user.CreatedAt = createdAt.Time()
	} else {
		user.CreatedAt = time.Now() // fallback
	}
	
	if updatedAt, ok := rawUser["updatedAt"].(primitive.DateTime); ok {
		user.UpdatedAt = updatedAt.Time()
	} else {
		user.UpdatedAt = user.CreatedAt // fallback to createdAt
	}

	return &user, nil
}

func (s *AuthService) UpdateUser(userID string, updates bson.M) (*models.User, error) {
	collection := s.db.Collection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	updates["updated_at"] = time.Now()

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return nil, err
	}

	return s.GetUserByID(userID)
}

func (s *AuthService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 дней
		"iat":     time.Now().Unix(),
		"jti":     uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// Верификация email
func (s *AuthService) VerifyEmail(req models.VerifyEmailRequest) (map[string]interface{}, error) {
	collection := s.db.Collection("users")

	// Получаем raw данные пользователя
	var rawUser bson.M
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&rawUser)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Проверяем верификацию
	verified, _ := rawUser["verified"].(bool)
	if verified {
		return map[string]interface{}{
			"success": true,
			"message": "Email already verified",
		}, nil
	}

	// Проверяем код и срок действия
	verificationCode, _ := rawUser["verificationCode"].(string)
	verificationExpires, _ := rawUser["verificationExpires"].(primitive.DateTime)
	
	if verificationCode != req.Code || verificationExpires.Time().Before(time.Now()) {
		return nil, errors.New("invalid or expired verification code")
	}

	// Верифицируем пользователя
	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"email": req.Email},
		bson.M{
			"$set": bson.M{"verified": true},
			"$unset": bson.M{
				"verificationCode":    "",
				"verificationExpires": "",
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"message": "Email verified successfully",
	}, nil
}

// Повторная отправка кода верификации
func (s *AuthService) ResendVerificationCode(req models.ResendCodeRequest) (map[string]interface{}, error) {
	collection := s.db.Collection("users")

	// Получаем raw данные пользователя
	var rawUser bson.M
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&rawUser)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Проверяем верификацию
	verified, _ := rawUser["verified"].(bool)
	if verified {
		return nil, errors.New("email already verified")
	}

	// Получаем email для отправки
	email, _ := rawUser["email"].(string)

	// Генерируем новый код
	code := s.generateVerificationCode()
	codeExpires := time.Now().Add(10 * time.Minute)

	// Обновляем код в базе
	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"email": req.Email},
		bson.M{
			"$set": bson.M{
				"verificationCode":    code,
				"verificationExpires": codeExpires,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	// Отправляем новый код на email
	if s.emailService != nil {
		go s.emailService.SendVerificationEmail(email, code)
	}

	return map[string]interface{}{
		"success": true,
		"message": "Verification code sent to your email",
	}, nil
}