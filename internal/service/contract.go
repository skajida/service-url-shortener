//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package service

import "context"

type UrlRepository interface {
	AddEntry(ctx context.Context, url, shortUrl string) error
	FindEntry(ctx context.Context, shortUrl string) (url string, err error)
}
