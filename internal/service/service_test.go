package service_test

import (
	"context"
	"sync"
	"testing"
	"url-shortener/internal/entity"
	"url-shortener/internal/model"
	"url-shortener/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUrlShortener_GenerateShortUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		originUrl string
		prepare   func(*MockUrlRepository)
		expected  func(assert.TestingT, error)
	}{
		{
			name:      "no error",
			originUrl: "https://www.ozon.ru/",
			prepare: func(mock *MockUrlRepository) {
				mock.EXPECT().AddEntry(gomock.Any(), "https://www.ozon.ru/", gomock.Len(10)).Return(nil)
			},
			expected: func(t assert.TestingT, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:      "origin url already exists",
			originUrl: "https://www.avito.ru/",
			prepare: func(mock *MockUrlRepository) {
				mock.EXPECT().AddEntry(gomock.Any(), "https://www.avito.ru/", gomock.Len(10)).Return(model.ErrOriginConflict)
			},
			expected: func(t assert.TestingT, err error) {
				assert.ErrorIs(t, err, model.ErrOriginConflict)
			},
		},
		{
			name:      "short url already exists once",
			originUrl: "https://github.com/",
			prepare: func(mock *MockUrlRepository) {
				mock.EXPECT().AddEntry(gomock.Any(), "https://github.com/", gomock.Len(10)).Return(model.ErrShortConflict).Times(1)
				mock.EXPECT().AddEntry(gomock.Any(), "https://github.com/", gomock.Len(10)).Return(nil).Times(1)
			},
			expected: func(t assert.TestingT, err error) {
				assert.NoError(t, err)
			},
		},
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	ctx, turnOffGen := context.WithCancel(context.Background())
	defer turnOffGen()
	gen := entity.NewTokenGenerator(ctx, &wg)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			urlRepo := NewMockUrlRepository(ctrl)
			tc.prepare(urlRepo)

			urlShortener := service.NewUrlShortener(gen, urlRepo)
			_, err := urlShortener.Generate(context.Background(), tc.originUrl)
			tc.expected(t, err)
		})
	}
}

func TestUrlShortener_GetOriginUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		shortUrl string
		prepare  func(*MockUrlRepository)
		expected func(assert.TestingT, string, error)
	}{
		{
			name:     "correct respond",
			shortUrl: "https://ya.ru/",
			prepare: func(mock *MockUrlRepository) {
				mock.EXPECT().FindEntry(gomock.Any(), "https://ya.ru/").Return("https://yandex.ru/", nil)
			},
			expected: func(t assert.TestingT, originUrl string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://yandex.ru/", originUrl)
			},
		},
		{
			name:     "short url doesn't exist",
			shortUrl: "https://ya.cc/m/nnw8B41",
			prepare: func(mock *MockUrlRepository) {
				mock.EXPECT().FindEntry(gomock.Any(), "https://ya.cc/m/nnw8B41").Return("", model.ErrShortBadParam)
			},
			expected: func(t assert.TestingT, originUrl string, err error) {
				assert.ErrorIs(t, err, model.ErrShortBadParam)
			},
		},
	}

	var wg sync.WaitGroup
	defer wg.Wait()
	ctx, turnOffGen := context.WithCancel(context.Background())
	defer turnOffGen()
	gen := entity.NewTokenGenerator(ctx, &wg)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			urlRepo := NewMockUrlRepository(ctrl)
			tc.prepare(urlRepo)

			urlShortener := service.NewUrlShortener(gen, urlRepo)
			originUrl, err := urlShortener.LookUp(context.Background(), tc.shortUrl)
			tc.expected(t, originUrl, err)
		})
	}
}
