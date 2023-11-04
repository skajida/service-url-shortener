package entity

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

const (
	tokenLength  = 10
	base         = int64(len(characterSet))
	characterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
)

type TokenGenerator struct {
	gen *rand.Rand
}

func NewTokenGenerator(ctx context.Context, wg *sync.WaitGroup) <-chan string {
	tokenGen := TokenGenerator{gen: rand.New(rand.NewSource(time.Now().UnixNano()))}
	return tokenGen.makeGen(ctx, wg)
}

func int64ToShortUrl(seed int64) string {
	shortUrl := make([]byte, tokenLength)
	for i := range shortUrl {
		shortUrl[i] = characterSet[seed%base]
		seed /= base
	}
	return string(shortUrl)
}

func binPow(a, n int64) int64 {
	var power int64 = 1
	for n != 0 {
		if n%2 == 1 {
			power *= a
		}
		a *= a
		n >>= 1
	}
	return power
}

func (tg TokenGenerator) yield() string {
	return int64ToShortUrl(tg.gen.Int63n(binPow(base, tokenLength)))
}

func (tg TokenGenerator) makeGen(ctx context.Context, wg *sync.WaitGroup) <-chan string {
	channel := make(chan string, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(channel)
		for {
			select {
			case <-ctx.Done():
				return
			case channel <- tg.yield():
			}
		}
	}()
	return channel
}
