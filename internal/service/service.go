package service

import (
	"context"
	"fmt"
	"url-shortener/internal/model"
)

type UrlShortener struct {
	urlGen <-chan string
	repo   UrlRepository
}

func NewUrlShortener(shortUrlGenerator <-chan string, repository UrlRepository) *UrlShortener {
	return &UrlShortener{urlGen: shortUrlGenerator, repo: repository}
}

func (us *UrlShortener) Generate(ctx context.Context, originUrl string) (string, error) {
	for shortUrl := range us.urlGen {
		switch err := us.repo.AddEntry(ctx, originUrl, shortUrl); err {
		case nil:
			return shortUrl, nil
		case model.ErrShortConflict:
		case model.ErrOriginConflict:
			return "", err
		default:
			return "", fmt.Errorf("us: %w", err)
		}
	}
	return "", model.ErrUrlGeneratorInternal
}

func (us UrlShortener) LookUp(ctx context.Context, shortUrl string) (string, error) {
	originUrl, err := us.repo.FindEntry(ctx, shortUrl)
	if err != nil {
		return "", err
	}
	return originUrl, nil
}
