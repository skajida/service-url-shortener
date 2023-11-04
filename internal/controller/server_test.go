package controller_test

import (
	"context"
	"testing"
	"url-shortener/internal/controller"
	"url-shortener/internal/model"
	"url-shortener/internal/pb"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServer_ReduceUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		originUrl string
		prepare   func(*MockurlShortener)
		expected  func(assert.TestingT, *pb.ShortUrl, error)
	}{
		{
			name:      "successful reducing",
			originUrl: "https://www.ozon.ru/",
			prepare: func(mock *MockurlShortener) {
				mock.EXPECT().Generate(gomock.Any(), "https://www.ozon.ru/").Return("_MnpjKe_mn", nil)
			},
			expected: func(t assert.TestingT, shortUrl *pb.ShortUrl, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "_MnpjKe_mn", shortUrl.GetShortUrl())
			},
		},
		{
			name:      "origin url already exists",
			originUrl: "https://www.avito.ru/",
			prepare: func(mock *MockurlShortener) {
				mock.EXPECT().Generate(gomock.Any(), "https://www.avito.ru/").Return("", model.ErrOriginConflict)
			},
			expected: func(t assert.TestingT, shortUrl *pb.ShortUrl, err error) {
				assert.ErrorIs(t, err, status.Error(codes.AlreadyExists, "origin url entry already exists"))
			},
		},
		{
			name:      "short url generator closed channel",
			originUrl: "https://mail.ru/",
			prepare: func(mock *MockurlShortener) {
				mock.EXPECT().Generate(gomock.Any(), "https://mail.ru/").Return("", model.ErrUrlGeneratorInternal)
			},
			expected: func(t assert.TestingT, shortUrl *pb.ShortUrl, err error) {
				assert.ErrorIs(t, err, status.Error(codes.Internal, "internal service error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			urlShortener := NewMockurlShortener(ctrl)
			tc.prepare(urlShortener)

			controllerServer := controller.NewServer(zap.NewNop(), urlShortener)
			shortUrl, err := controllerServer.ReduceUrl(context.Background(), &pb.OriginUrl{OriginUrl: tc.originUrl})
			tc.expected(t, shortUrl, err)
		})
	}
}

func TestServer_GetOriginUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		shortUrl string
		prepare  func(*MockurlShortener)
		expected func(assert.TestingT, *pb.OriginUrl, error)
	}{
		{
			name:     "origin url found",
			shortUrl: "_MnpjKe_mn",
			prepare: func(mock *MockurlShortener) {
				mock.EXPECT().LookUp(gomock.Any(), "_MnpjKe_mn").Return("https://www.ozon.ru/", nil)
			},
			expected: func(t assert.TestingT, originUrl *pb.OriginUrl, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "https://www.ozon.ru/", originUrl.GetOriginUrl())
			},
		},
		{
			name:     "no such short url entry",
			shortUrl: "On0jz2Mqd5",
			prepare: func(ms *MockurlShortener) {
				ms.EXPECT().LookUp(gomock.Any(), "On0jz2Mqd5").Return("", model.ErrShortBadParam)
			},
			expected: func(t assert.TestingT, ou *pb.OriginUrl, err error) {
				assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "shorted url entry doesn't exist"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			urlShortener := NewMockurlShortener(ctrl)
			tc.prepare(urlShortener)

			controllerServer := controller.NewServer(zap.NewNop(), urlShortener)
			originUrl, err := controllerServer.GetOriginUrl(context.Background(), &pb.ShortUrl{ShortUrl: tc.shortUrl})
			tc.expected(t, originUrl, err)
		})
	}
}
