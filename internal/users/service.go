package users

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	store "github.com/assbomber/myzone/internal/store/sqlc"
	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/assbomber/myzone/pkg/utils"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	CreateOnboardingDetails(context.Context, int64, OnboardingRequest) error
	GetUserDetails(context.Context, int64) (*store.GetUserDetailsRow, error)
}

type userService struct {
	logger  *logger.Logger
	db      *sql.DB
	queries *store.Queries
	redisIn *redis.Client
}

// Constructor for user service
func NewService(logger *logger.Logger, db *sql.DB, queries *store.Queries, redisIn *redis.Client) Service {
	return &userService{
		logger:  logger,
		db:      db,
		queries: queries,
		redisIn: redisIn,
	}
}

func (us *userService) CreateOnboardingDetails(ctx context.Context, userID int64, args OnboardingRequest) error {

	tx, err := us.db.Begin()
	if err != nil {
		return fmt.Errorf("error initiating transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := us.queries.WithTx(tx)

	bdDate, err := time.ParseInLocation("2006-01-02", args.BirthDate, utils.GetISTLocation())
	if err != nil {
		return errors.Wrap(err, "err parsing birthdate")
	}
	// Updating user basic details
	err = qtx.UpdateBasicUserDetails(ctx, store.UpdateBasicUserDetailsParams{
		UserID: userID,
		Gender: store.NullGenders{
			Valid:   true,
			Genders: store.GendersFemale,
		},
		BirthDate: sql.NullTime{
			Valid: true,
			Time:  bdDate,
		},
		Name: sql.NullString{
			Valid:  true,
			String: args.Name,
		},
	})
	if err != nil {
		return fmt.Errorf("error updating basic user details: %w", err)
	}

	// Inactivating user previous active location
	err = qtx.InactiveUserLocation(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "error inactivating user location")
	}

	// Creating location and setting it to active mode
	_, err = qtx.CreateUserLocation(ctx, store.CreateUserLocationParams{
		UserID:    userID,
		Active:    true,
		Longitude: args.Location.Longitude,
		Latitude:  args.Location.Latitude,
	})
	if err != nil {
		return errors.Wrap(err, "error updating location user details")
	}

	return tx.Commit()
}

func (us *userService) GetUserDetails(ctx context.Context, userID int64) (*store.GetUserDetailsRow, error) {
	details, err := us.queries.GetUserDetails(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoSuchUser
		}
		return nil, errors.Wrap(err, "Error getting user details")
	}
	return &details, nil
}
