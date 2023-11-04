//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package controller

import "context"

type urlShortener interface {
	Generate(ctx context.Context, originUrl string) (string, error)
	LookUp(ctx context.Context, shortUrl string) (string, error)
}
