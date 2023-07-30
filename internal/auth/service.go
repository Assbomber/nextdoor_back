package auth

import (
	"context"
	"database/sql"
	"time"

	store "github.com/assbomber/myzone/internal/store/sqlc"
	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Interface for auth Service
type Service interface {
	Login(context.Context, LoginRequest) (*LoginResponse, error)
	Register(context.Context, RegisterRequest) (*RegisterResponse, error)
	generateJWT(userID int64) (string, error)
	validateJWT(tokenStr string) (*MyCustomClaims, error)
}

// Struct for JWT custom claims
type MyCustomClaims struct {
	UserID int64
	jwt.RegisteredClaims
}

// Auth service struct
type authService struct {
	jwtSecret string
	logger    *logger.Logger
	queries   *store.Queries
}

// Constructor for auth service
func NewService(logger *logger.Logger, jwtSecret string, queries *store.Queries) Service {
	return &authService{
		jwtSecret: jwtSecret,
		logger:    logger,
		queries:   queries,
	}
}

// Verifies the user creds and returns JWT
func (as *authService) Login(ctx context.Context, args LoginRequest) (*LoginResponse, error) {

	// Fetching user in db
	user, err := as.queries.GetUserByEmail(ctx, args.Email)
	if err != nil {
		// No rows, return user not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoSuchUser
		}
		return nil, errors.Wrap(err, "Error getting user")
	}

	// comparing passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(args.Password))
	if err != nil { // password not match
		return nil, constants.ErrWrongPassword
	}

	// Creating JWT
	token, err := as.generateJWT(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Err creating JWT")
	}
	return &LoginResponse{
		Token: token,
	}, nil
}

// Registers a new user and returns JWT
func (as *authService) Register(ctx context.Context, args RegisterRequest) (*RegisterResponse, error) {

	// Encryping pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "Err generating hash password")
	}

	// Creating user in db
	user, err := as.queries.CreateUser(ctx, store.CreateUserParams{
		Name:     args.Name,
		Email:    args.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Error creating user")
	}

	// Creating JWT
	token, err := as.generateJWT(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Err creating JWT")
	}

	return &RegisterResponse{
		Token: token,
	}, nil
}

// Generates JWT using user id, or else returns err if any.
func (as *authService) generateJWT(userID int64) (string, error) {
	claims := MyCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 6, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// Helps validate JWT using provided secret in JWT_SECRET environment variable.
// If Success, returns MyCustomClaims, else error
func (as *authService) validateJWT(tokenStr string) (*MyCustomClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		// verifing if signing method is same
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrUnexpectedSigningMethod
		}
		return []byte(as.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(MyCustomClaims); ok && token.Valid {
		return &claims, nil
	} else {
		return nil, constants.ErrInvalidJWT
	}
}
