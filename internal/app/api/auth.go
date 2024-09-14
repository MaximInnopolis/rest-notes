package api

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"rest-notes/internal/app/models"
	"rest-notes/internal/app/repository"
	"rest-notes/internal/app/repository/postgresql"
)

var ErrUserAlreadyExists = errors.New("пользователь уже существует")

// AuthService provides authentication services using user repository
type AuthService struct {
	repo repository.UserRepo
}

// NewAuthService creates new instance of AuthService
func NewAuthService(repo repository.UserRepo) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser creates new user using repository and returns created user
func (as *AuthService) CreateUser(user models.User) error {
	_, err := as.repo.Get(user)
	if err == nil {
		log.Printf("Создание пользователя не удалось: %s уже существует", user.Name)
		return ErrUserAlreadyExists
	}

	if !errors.Is(err, postgresql.ErrNotFound) {
		log.Printf("Ошибка при получении пользователя: %v", err)
		return err
	}

	// Hash user's password before saving
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		log.Printf("Ошибка при хэшировании пароля: %v", err)
		return err
	}

	// Save new user in repository
	err = as.repo.Create(user)
	if err != nil {
		log.Printf("Ошибка при создании пользователя: %v", err)
		return err
	}

	log.Printf("Пользователь успешно создан: %s", user.Name)
	return nil
}

// GenerateToken generates JWT for authenticated user
// It retrieves user from repository and creates signed JWT token
func (as *AuthService) GenerateToken(user models.User) (string, error) {
	dbUser, err := as.repo.Get(user)
	if err != nil {
		log.Printf("Ошибка при получении пользователя для генерации токена: %v", err)
		return "", err
	}

	// Generate JWT token
	token, err := generateJWT(dbUser)
	if err != nil {
		log.Printf("Ошибка при генерации JWT: %v", err)
		return "", err
	}

	log.Printf("JWT успешно сгенерирован для пользователя: %s", dbUser.Name)
	return token, nil
}

// IsTokenValid validates given JWT
// It checks token's signature, claims, and expiration time
func (as *AuthService) IsTokenValid(tokenString string) (bool, jwt.MapClaims, error) {
	// Check token validity
	validToken, claims, err := checkToken(tokenString)
	if err != nil || !validToken {
		log.Printf("Неверный токен: %v", err)
		return false, nil, errors.New("invalid token")
	}

	log.Printf("Токен валиден")
	return true, claims, nil
}

// checkToken parses and validates JWT
// It verifies token's signature and checks expiration claim
func checkToken(tokenString string) (bool, jwt.MapClaims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		log.Printf("Ошибка при разборе токена: %v", err)
		return false, nil, err
	}

	// Check if token is valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Printf("Некорректные claims или токен недействителен")
		return false, nil, nil
	}

	// Check if expiration claim exists and validate it
	expiration, ok := claims["exp"].(float64)
	if !ok {
		log.Printf("Некорректное время истечения токена")
		return false, nil, nil
	}

	if int64(expiration) < time.Now().Unix() {
		log.Printf("Токен истек")
		return false, nil, nil
	}

	return true, claims, nil
}

// generateJWT generates JWT for provided user with 24-hour expiration time
func generateJWT(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set standard claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["sub"] = user.Name

	// Add additional claims
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Printf("Ошибка при подписании токена: %v", err)
		return "", err
	}
	return tokenString, nil
}

// generatePasswordHash hashes user's password using bcrypt
func generatePasswordHash(password string) (string, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Ошибка хэширования пароля: %s", err)
		return "", err
	}
	return string(hashedPassword), nil
}
