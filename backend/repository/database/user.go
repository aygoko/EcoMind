package repository

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

type UserRepository struct {
	db        *gorm.DB
	redisPool *redis.Pool
}

func NewUserRepository(db *gorm.DB, redisPool *redis.Pool) repository.UserService {
	return &UserRepository{
		db:        db,
		redisPool: redisPool,
	}
}

func (r *UserRepository) Get(login string) (*repository.User, error) {
	return r.getCachedUser("user:username:"+login, func() (*repository.User, error) {
		return r.db.Where("username = ?", login).First(&repository.User{}).Take()
	})
}

func (r *UserRepository) GetByEmail(email string) (*repository.User, error) {
	return r.getCachedUser("user:email:"+email, func() (*repository.User, error) {
		return r.db.Where("email = ?", email).First(&repository.User{}).Take()
	})
}

func (r *UserRepository) GetByPhoneNumber(phoneNumber string) (*repository.User, error) {
	return r.getCachedUser("user:phone:"+phoneNumber, func() (*repository.User, error) {
		return r.db.Where("phone_number = ?", phoneNumber).First(&repository.User{}).Take()
	})
}

func (r *UserRepository) GetByID(id uint) (*repository.User, error) {
	key := "user:id:" + strconv.FormatUint(uint64(id), 10)
	return r.getCachedUser(key, func() (*repository.User, error) {
		return r.db.Where("id = ?", id).First(&repository.User{}).Take()
	})
}

func (r *UserRepository) Create(user *repository.User) error {
	if exists, _ := r.checkUserExistence(user.Username, user.Email, user.PhoneNumber); exists {
		return errors.New("user already exists")
	}

	if err := r.db.Create(user).Error; err != nil {
		return errors.New("database error")
	}

	r.clearUserCache(user.Username, user.Email, user.PhoneNumber, user.ID)
	return nil
}

func (r *UserRepository) Update(user *repository.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errors.New("database error")
	}

	r.clearUserCache(user.Username, user.Email, user.PhoneNumber, user.ID)
	return nil
}

func (r *UserRepository) UpdateAuthToken(id uint, token string) error {
	if err := r.db.Model(&repository.User{}).Where("id = ?", id).Update("auth_token", token).Error; err != nil {
		return errors.New("database error")
	}

	conn := r.redisPool.Get()
	defer conn.Close()
	conn.Do("DEL", "user:id:"+strconv.FormatUint(uint64(id), 10))
	return nil
}

// Helper function to handle Redis caching
func (r *UserRepository) getCachedUser(key string, queryFunc func() (*repository.User, error)) (*repository.User, error) {
	conn := r.redisPool.Get()
	defer conn.Close()

	userJSON, err := redis.String(conn.Do("GET", key))
	if err == nil {
		var user repository.User
		if err := json.Unmarshal([]byte(userJSON), &user); err == nil {
			return &user, nil
		}
	}

	user, err := queryFunc()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	userJSON, _ = json.Marshal(user)
	conn.Do("SET", key, userJSON, "EX", 3600)
	return user, nil
}

// Optimized check for existing user details
func (r *UserRepository) checkUserExistence(username, email, phoneNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&repository.User{}).Where(
		"username = ? OR email = ? OR phone_number = ?",
		username, email, phoneNumber,
	).Count(&count).Error
	if err != nil {
		return false, errors.New("database error")
	}
	return count > 0, nil
}

// Clears all relevant Redis keys for a user
func (r *UserRepository) clearUserCache(username, email, phoneNumber string, id uint) {
	conn := r.redisPool.Get()
	defer conn.Close()

	keys := []string{
		"user:username:" + username,
		"user:email:" + email,
		"user:phone:" + phoneNumber,
		"user:id:" + strconv.FormatUint(uint64(id), 10),
	}
	_, _ = conn.Do("DEL", keys...)
}
