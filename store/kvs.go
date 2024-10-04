package store

import (
	"Bangseungjae/go-todo-app/config"
	"context"
	"github.com/go-redis/redis/v8"
)

func NewKVS(ctx context.Context, cfg *config.Config) (*KVS, error) {

}

type KVS struct {
	Cli *redis.Client
}
