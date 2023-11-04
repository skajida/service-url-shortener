package repository

import (
	"context"
	"sync"
	"url-shortener/internal/model"
)

type InMemoryDatabase struct {
	shortToOrigin map[string]string
	originToShort map[string]string
	mapMutex      sync.RWMutex
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{shortToOrigin: make(map[string]string), originToShort: make(map[string]string)}
}

func (im *InMemoryDatabase) checkIfExists(ctx context.Context, originUrl, shortUrl string) error {
	im.mapMutex.RLock()
	defer im.mapMutex.RUnlock()
	if _, exists := im.originToShort[originUrl]; exists {
		return model.ErrOriginConflict
	}
	if _, exists := im.shortToOrigin[shortUrl]; exists {
		return model.ErrShortConflict
	}
	return nil
}

func (im *InMemoryDatabase) AddEntry(ctx context.Context, originUrl, shortUrl string) error {
	if err := im.checkIfExists(ctx, originUrl, shortUrl); err != nil {
		return err
	}
	im.mapMutex.Lock()
	defer im.mapMutex.Unlock()
	im.shortToOrigin[shortUrl] = originUrl
	im.originToShort[originUrl] = shortUrl
	return nil
}

func (im *InMemoryDatabase) FindEntry(ctx context.Context, shortUrl string) (string, error) {
	im.mapMutex.RLock()
	defer im.mapMutex.RUnlock()
	if originUrl, exists := im.shortToOrigin[shortUrl]; exists {
		return originUrl, nil
	}
	return "", model.ErrShortBadParam
}
