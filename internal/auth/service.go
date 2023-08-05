package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/smtp"
	"time"

	"github.com/assbomber/myzone/configs"
	store "github.com/assbomber/myzone/internal/store/sqlc"
	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/assbomber/myzone/pkg/utils"

	"errors"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

// Interface for auth Service
type Service interface {
	Login(context.Context, LoginRequest) (*LoginResponse, error)
	Register(context.Context, RegisterRequest) (*RegisterResponse, error)
	SendVerificationEmail(context.Context, string) error
	ResetPassword(context.Context, ResetPasswordRequest) error
}

// Auth service struct
type authService struct {
	jwtSecret string
	logger    *logger.Logger
	db        *sql.DB
	queries   *store.Queries
	redisIn   *redis.Client
}

// Constructor for auth service
func NewService(logger *logger.Logger, jwtSecret string, db *sql.DB, queries *store.Queries, redisIn *redis.Client) Service {
	return &authService{
		jwtSecret: jwtSecret,
		logger:    logger,
		db:        db,
		queries:   queries,
		redisIn:   redisIn,
	}
}

// Verifies the user creds and returns JWT
func (as *authService) Login(ctx context.Context, args LoginRequest) (*LoginResponse, error) {

	// Fetching user in db and updating lastlogin for him
	user, err := as.queries.UpdateUserLoginTimeByEmail(ctx, store.UpdateUserLoginTimeByEmailParams{
		Email:     args.Email,
		LastLogin: time.Now(),
	})
	if err != nil {
		// No rows, return user not found
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoSuchUser
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	// comparing passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(args.Password))
	if err != nil { // password not match
		return nil, constants.ErrWrongPassword
	}

	// Creating JWT
	token, err := utils.GenerateJWT(user.ID, as.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("error creating JWT: %w", err)
	}
	return &LoginResponse{
		Token: token,
	}, nil
}

// Registers a new user and returns JWT
func (as *authService) Register(ctx context.Context, args RegisterRequest) (*RegisterResponse, error) {

	// Retreiving otp from redis
	redisRes := as.redisIn.Get(ctx, utils.GetOTPRedisKey(args.Email))
	if redisRes.Err() == redis.Nil || redisRes.Val() != fmt.Sprint(args.OTP) {
		return nil, constants.ErrInvalidOTP
	}

	// check if user already exist
	user, err := as.queries.GetUserByEmailOrUsername(ctx, store.GetUserByEmailOrUsernameParams{
		Email:    args.Email,
		Username: args.Username,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	// User already exist
	if user.ID != 0 {
		if user.Username == args.Username {
			return nil, constants.ErrUsernameAlreadyExist
		}
		return nil, constants.ErrEmailAlreadyExist
	}

	// Encryping pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error generating hash password: %w", err)
	}

	// Creating user in db
	user, err = as.queries.CreateUser(ctx, store.CreateUserParams{
		Username:  args.Username,
		Email:     args.Email,
		Password:  string(hashedPassword),
		LastLogin: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Creating JWT
	token, err := utils.GenerateJWT(user.ID, as.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("error creating JWT: %w", err)
	}

	return &RegisterResponse{
		Token: token,
	}, nil
}

// Send the OTP to the provided email for verification
func (as *authService) SendVerificationEmail(ctx context.Context, email string) error {
	from := "myzoneapp01@gmail.com"
	toList := []string{email}

	// Compose the email
	msg, otp := utils.GetOTPEmail()

	smtpHost := configs.GetString("smtpHost")
	smtpPort := configs.GetString("smtpPort")
	smtpSecret := configs.GetString("SMTP_SECRET")

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", from, smtpSecret, smtpHost)

	// Save to redis with expiry of 2 mins
	if err := as.redisIn.SetEx(ctx, utils.GetOTPRedisKey(email), otp, 2*time.Minute).Err(); err != nil {
		return fmt.Errorf("error saving OTP to redis: %w", err)
	}

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toList, []byte(msg))
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}

// Resets the user password
func (as *authService) ResetPassword(ctx context.Context, args ResetPasswordRequest) error {
	// Retreiving otp from redis
	redisRes := as.redisIn.Get(ctx, utils.GetOTPRedisKey(args.Email))
	if redisRes.Err() == redis.Nil || redisRes.Val() != fmt.Sprint(args.OTP) {
		return constants.ErrInvalidOTP
	}

	// Encryping pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error generating pass hash: %w", err)
	}

	// Updating user pass in db
	err = as.queries.UpdateUserPasswordByEmail(ctx, store.UpdateUserPasswordByEmailParams{
		Email:    args.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return fmt.Errorf("error updating user password: %w", err)
	}

	return nil
}
