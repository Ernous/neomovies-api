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

	"neomovies-api/internal/models"
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

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Проверяем верификацию email
	if !user.Verified {
		return nil, errors.New("email not verified. Check your email for verification code.")
	}

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