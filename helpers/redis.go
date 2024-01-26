package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"ngc-grpc/model"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetUserCache(c *redis.Client, user *model.UserCache, ctx context.Context) error {
	p, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := "user" + user.ID

	err = c.Set(ctx, key, p, 1*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUserCache(c *redis.Client, id string, ctx context.Context) (model.UserCache, error) {
	var user model.UserCache
	key := "user" + user.ID

	str, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return model.UserCache{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return model.UserCache{}, err
	}

	err = json.Unmarshal([]byte(str), &user)
	if err != nil {
		return model.UserCache{}, err
	}

	return user, nil
}
