package repository

import (
    "context"
    "encoding/json"
    "errors"
    "time"

    "github.com/aygoko/EcoMInd/backend/domain"
    "github.com/go-redis/redis/v8"
    "github.com/lib/pq"
    "database/sql"
)

const (
    redisUserKeyPrefix = "user:login:"
    cacheTTL           = 5 * time.Minute
)

// Helper function to scan a database row into a User struct
func scanUserRow(row *sql.Row, user *domain.User) error {
    return row.Scan(&user.Login, &user.Email, &user.PhoneNumber)
}

type UserRepositoryDB struct {
    DB         *sql.DB
    RedisClient *redis.Client
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB, redisClient *redis.Client) domain.UserService {
    return &UserRepositoryDB{
        DB:         db,
        RedisClient: redisClient,
    }
}

// cacheUser stores the user in Redis cache
func (r *UserRepositoryDB) cacheUser(user *domain.User) error {
    key := redisUserKeyPrefix + user.Login
    userJSON, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return r.RedisClient.Set(context.Background(), key, userJSON, cacheTTL).Err()
}

// Get retrieves a user by login with cache check
func (r *UserRepositoryDB) Get(login string) (*domain.User, error) {
    key := redisUserKeyPrefix + login

    // Check Redis cache first
    userJSON, err := r.RedisClient.Get(context.Background(), key).Result()
    if err == nil {
        var user domain.User
        if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
            return nil, err
        }
        return &user, nil
    } else if !errors.Is(err, redis.Nil) {
        return nil, err
    }

    // Query database if not found in cache
    row := r.DB.QueryRow("SELECT login, email, phone_number FROM users WHERE login = $1", login)
    var user domain.User
    if err := scanUserRow(row, &user); err != nil {
        return nil, err
    }

    // Cache the result in Redis
    if err := r.cacheUser(&user); err != nil {
        // Log error but proceed
    }

    return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepositoryDB) GetByEmail(email string) (*domain.User, error) {
    row := r.DB.QueryRow("SELECT login, email, phone_number FROM users WHERE email = $1", email)
    var user domain.User
    if err := scanUserRow(row, &user); err != nil {
        return nil, err
    }

    // Cache the user in Redis
    if err := r.cacheUser(&user); err != nil {
        // Log error but proceed
    }

    return &user, nil
}

// GetByPhoneNumber retrieves a user by phone number
func (r *UserRepositoryDB) GetByPhoneNumber(phoneNumber string) (*domain.User, error) {
    row := r.DB.QueryRow("SELECT login, email, phone_number FROM users WHERE phone_number = $1", phoneNumber)
    var user domain.User
    if err := scanUserRow(row, &user); err != nil {
        return nil, err
    }

    // Cache the user in Redis
    if err := r.cacheUser(&user); err != nil {
        // Log error but proceed
    }

    return &user, nil
}