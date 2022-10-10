package app

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"hamster-client/module/keystorage"
)

type KeyStorage struct {
	ctx    context.Context
	client keystorage.Client
}

func NewKeyStorageApp(service *keystorage.Client) KeyStorage {
	return KeyStorage{
		ctx:    context.Background(),
		client: *service,
	}
}

func (k *KeyStorage) Get(key string) string {
	v := k.client.Get(key)
	if k.client.Err() != nil {
		runtime.LogErrorf(k.ctx, "kv storage get error: %s", k.client.Err())
		return ""
	}
	return v
}

func (k *KeyStorage) Set(key, value string) {
	k.client.Set(key, value)
	if k.client.Err() != nil {
		runtime.LogErrorf(k.ctx, "kv storage set error: %s", k.client.Err())
	}
}

func (k *KeyStorage) Delete(key string) {
	k.client.Delete(key)
	if k.client.Err() != nil {
		runtime.LogErrorf(k.ctx, "kv storage delete error: %s", k.client.Err())
	}
}
