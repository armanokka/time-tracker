package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
)

// SetupRedis launches local Redis instance via testcontainers. Returned testcontainers.Container MUST be terminated
func SetupRedis(ctx context.Context) (testcontainers.Container, *redis.Client) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Could not start redis: %s", err)
	}
	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	return redisC, redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
}

func TestAuthRedisRepo_GetUser(t *testing.T) {
	ctx := context.Background()

	redisC, rdb := SetupRedis(ctx)
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	defer rdb.Close()

	repo := NewAuthRedisRepo(rdb)

	user := getTestUser()
	assert.Nil(t, repo.SetUser(ctx, user, 60))
	gotUser, err := repo.GetUser(ctx, user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, gotUser)

	gotUser, err = repo.GetUser(ctx, 134)
	assert.Nil(t, gotUser)
	assert.Equal(t, err, redis.Nil)
}

func TestAuthRedisRepo_SetUser(t *testing.T) {
	ctx := context.Background()

	redisC, rdb := SetupRedis(ctx)
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	defer rdb.Close()

	repo := NewAuthRedisRepo(rdb)

	user := getTestUser()
	assert.Nil(t, repo.SetUser(ctx, user, 3))
	gotUser, err := repo.GetUser(ctx, user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, gotUser)
}

func TestAuthRedisRepo_DeleteUser(t *testing.T) {
	ctx := context.Background()

	redisC, rdb := SetupRedis(ctx)
	defer func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	defer rdb.Close()

	repo := NewAuthRedisRepo(rdb)

	user := getTestUser()
	assert.Nil(t, repo.SetUser(ctx, user, 3))
	assert.Nil(t, repo.DeleteUser(ctx, user.ID))
	gotUser, err := repo.GetUser(ctx, user.ID)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, gotUser)
}
