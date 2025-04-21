package repository

import (
    "context"
    "encoding/json"
    "errors"
    "log"
    "time"

    "github.com/aygoko/EcoMInd/backend/domain"
    "github.com/go-redis/redis/v8"
    "database/sql"
)

const (
    redisUserKeyPrefix = "user:login:"
    cacheTTL           = 5 * time.Minute
)

// Logger interface for structured logging
type Logger interface {
    Errorf(format string, args ...interface{})
    Infof(format string, args ...interface{})
}

// Default logger (replace with a structured logging library like zap or logrus in production)
var logger Logger = &defaultLogger{}

type defaultLogger struct{}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
    log.Printf("[ERROR] "+format, args...)
}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
    log.Printf("[INFO] "+format, args...)
}

// Helper function to scan a database row into a User struct
func scanUserRow(row *sql.Row, user *domain.User) error {
    return row.Scan(&user.Login, &user.Email, &user.PhoneNumber)
}

// UserRepositoryDB implements the UserService interface
type UserRepositoryDB struct {
    DB          *sql.DB
    RedisClient *redis.Client
    Logger      Logger
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB, redisClient *redis.Client, logger Logger) domain.UserService {
    if logger == nil {
        logger = &defaultLogger{}
    }
    return &UserRepositoryDB{
        DB:          db,
        RedisClient: redisClient,
        Logger:      logger,
    }
}

// CacheUser stores the user in Redis cache
func (r *UserRepositoryDB) cacheUser(user *domain.User) error {
    key := redisUserKeyPrefix + user.Login
    userJSON, err := json.Marshal(user)
    if err != nil {
        r.Logger.Errorf("failed to marshal user for caching: %v", err)
        return err
    }
    err = r.RedisClient.Set(context.Background(), key, userJSON, cacheTTL).Err()
    if err != nil {
        r.Logger.Errorf("failed to cache user in Redis: %v", err)
        return err
    }
    r.Logger.Infof("cached user with login: %s", user.Login)
    return nil
}

// InvalidateCache removes a user from the Redis cache
func (r *UserRepositoryDB) invalidateCache(login string) error {
    key := redisUserKeyPrefix + login
    err := r.RedisClient.Del(context.Background(), key).Err()
    if err != nil {
        r.Logger.Errorf("failed to invalidate cache for user: %s, error: %v", login, err)
        return err
    }
    r.Logger.Infof("invalidated cache for user with login: %s", login)
    return nil
}

// Get retrieves a user by login with cache check
func (r *UserRepositoryDB) Get(login string) (*domain.User, error) {
    key := redisUserKeyPrefix + login

    // Check Redis cache first
    userJSON, err := r.RedisClient.Get(context.Background(), key).Result()
    if err == nil {
        var user domain.User
        if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
            r.Logger.Errorf("failed to unmarshal cached user data: %v", err)
            return nil, err
        }
        r.Logger.Infof("retrieved user from cache with login: %s", login)
        return &user, nil
    } else if !errors.Is(err, redis.Nil) {
        r.Logger.Errorf("Redis error while fetching user: %v", err)
        return nil, err
    }

    // Query database if not found in cache
    row := r.DB.QueryRow("SELECT login, email, phone_number FROM users WHERE login = $1", login)
    var user domain.User
    if err := scanUserRow(row, &user); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            r.Logger.Infof("user not found in database with login: %s", login)
            return nil, domain.ErrUserNotFound
        }
        r.Logger.Errorf("database error while fetching user: %v", err)
        return nil, err
    }

    // Cache the result in Redis
    if err := r.cacheUser(&user); err != nil {
        r.Logger.Errorf("failed to cache user after database fetch: %v", err)
    }

    r.Logger.Infof("retrieved user from database with login: %s", login)
    return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepositoryDB) GetByEmail(email string) (*domain.User, error) {
    row := r.DB.QueryRow("SELECT login, email, phone_number FROM users WHERE email = $1", email)
    var user domain.User
    if err := scanUserRow(row, &user); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            r.Logger.Infof("user not found in database with email: %s", email)
            return nil, domain.ErrUserNotFound
        }
        r.Logger.Errorf("database error while fetching user by email: %v", err)
        return nil, err
    }

    // Cache the user in Redis
    if err := r.cacheUser(&user); err != nil {
        r.Logger.Errorf("failed to cache user after database fetch: %v", err)
    }

    r.Logger.Infof("retrieved user from database with email: %s", email)
    return &user, nil
}

// GetByPhoneNumber retrieves a user by phone number
func (r *UserRepositoryDB) GetByPhoneNumber(phoneNumber string) (*domain.User, error) {
    row := r.DB.QueryRow("SELECT login, email, phone_number FROM users WHERE phone_number = $1", phoneNumber)
    var user domain.User
    if err := scanUserRow(row, &user); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            r.Logger.Infof("user not found in database with phone number: %s", phoneNumber)
            return nil, domain.ErrUserNotFound
        }
        r.Logger.Errorf("database error while fetching user by phone number: %v", err)
        return nil, err
    }

    // Cache the user in Redis
    if err := r.cacheUser(&user); err != nil {
        r.Logger.Errorf("failed to cache user after database fetch: %v", err)
    }

    r.Logger.Infof("retrieved user from database with phone number: %s", phoneNumber)
    return &user, nil
}

// UpdateUser updates user data and invalidates the cache
func (r *UserRepositoryDB) UpdateUser(user *domain.User) error {
    _, err := r.DB.Exec("UPDATE users SET email = $1, phone_number = $2 WHERE login = $3", user.Email, user.PhoneNumber, user.Login)
    if err != nil {
        r.Logger.Errorf("failed to update user in database: %v", err)
        return err
    }

    // Invalidate cache
    if err := r.invalidateCache(user.Login); err != nil {
        r.Logger.Errorf("failed to invalidate cache after user update: %v", err)
    }

    r.Logger.Infof("updated user with login: %s", user.Login)
    return nil
}