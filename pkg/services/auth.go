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
	return &AuthService{
		db:           db,
		jwtSecret:    jwtSecret,
		emailService: emailService,
	}
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
	
	// Сначала попробуем найти пользователя без декодирования в структуру
	var rawUser bson.M
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&rawUser)
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
		
		return nil, errors.New("User not found")
	}
	
	fmt.Printf("Raw user found: %v\n", rawUser)
	
	// Теперь декодируем в структуру
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		fmt.Printf("Error decoding user to struct: %v\n", err)
		return nil, errors.New("User not found")
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

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
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

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Verified {
		return map[string]interface{}{
			"success": true,
			"message": "Email already verified",
		}, nil
	}

	// Проверяем код и срок действия
	if user.VerificationCode != req.Code || user.VerificationExpires.Before(time.Now()) {
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

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Verified {
		return nil, errors.New("email already verified")
	}

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
		go s.emailService.SendVerificationEmail(user.Email, code)
	}

	return map[string]interface{}{
		"success": true,
		"message": "Verification code sent to your email",
	}, nil
}