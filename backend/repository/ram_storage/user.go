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

type UserRepositoryDB struct {
    DB         *sql.DB
    RedisClient *redis.Client
}

func NewUserRepository(db *sql.DB, redisClient *redis.Client) domain.UserService {
    return &UserRepositoryDB{
        DB:         db,
        RedisClient: redisClient,
    }
}

func (r *UserRepositoryDB) Get(login string) (*domain.User, error) {
    ctx := context.Background()

    // Check Redis cache first
    userJSON, err := r.RedisClient.Get(ctx, "user:login:"+login).Result()
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
    err = row.Scan(&user.Login, &user.Email, &user.PhoneNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    // Cache the result in Redis
    userJSON, err = json.Marshal(user)
    if err != nil {
        return &user, nil // Proceed without caching
    }
    err = r.RedisClient.Set(ctx, "user:login:"+login, userJSON, 5*time.Minute).Err()
    if err != nil {
        // Handle error but proceed
    }

    return &user, nil
}

func (r *UserRepositoryDB) GetByEmail(email string) (*domain.User, error) {
    ctx := context.Background()

    // Query database directly
    row := r.DB.QueryRow("SELECT login, email, phone_number FROM users WHERE email = $1", email)
    var user domain.User
    err := row.Scan(&user.Login, &user.Email, &user.PhoneNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    // Cache the user by login in Redis
    userJSON, err := json.Marshal(user)
    if err != nil {
        return &user, nil // Proceed without caching
    }
    err = r.RedisClient.Set(ctx, "user:login:"+user.Login, userJSON, 5*time.Minute).Err()
    if err != nil {
        // Handle error but proceed
    }

    return &user, nil
}

func (r *UserRepositoryDB) GetByPhoneNumber(phoneNumber string) (*domain.User, error) {
    ctx := context.Background()

    // Query database directly
    row := r.DB.QueryRow("SELECT login, email, phone_number FROM users WHERE phone_number = $1", phoneNumber)
    var user domain.User
    err := row.Scan(&user.Login, &user.Email, &user.PhoneNumber)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    // Cache the user by login in Redis
    userJSON, err := json.Marshal(user)
    if err != nil {
        return &user, nil // Proceed without caching
    }
    err = r.RedisClient.Set(ctx, "user:login:"+user.Login, userJSON, 5*time.Minute).Err()
    if err != nil {
        // Handle error but proceed
    }

    return &user, nil
}
