package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"shortURL/dao"
	"shortURL/pkg/bloom"

	"golang.org/x/sync/singleflight"
)

var (
	bloomFilter *bloom.Filter
	sfGroup     singleflight.Group
)

func InitBloom(capacity uint, fpRate float64) {
	bloomFilter = bloom.New(capacity, fpRate)
}

func LoadBloomFromDB() error {
	surls, err := dao.LoadAllSurls()
	if err != nil {
		return fmt.Errorf("load surls from db: %w", err)
	}
	bloomFilter.AddBatch(surls) // 把所有存在的短链加入布隆过滤器
	log.Printf("bloom filter loaded with %d short urls", len(surls))
	return nil
}

func BloomAdd(key string) {
	bloomFilter.Add(key)
}

var (
	ErrNotFound = errors.New("short url not found")
	ErrDeleted  = errors.New("short url has been deleted")
)

func GetLongURL(ctx context.Context, shortCode string) (string, error) {
	if !bloomFilter.Contains(shortCode) {
		return "", ErrNotFound
	}

	cacheKey := dao.ShortURLKeyPrefix + shortCode
	lurl, err := dao.CacheGet(ctx, cacheKey)
	if err == nil {
		return lurl, nil
	}

	v, err, _ := sfGroup.Do(shortCode, func() (any, error) {
		record, err := dao.FindBySurl(shortCode)
		if err != nil {
			return nil, ErrNotFound
		}

		go func() { // 查完 DB 后，把结果异步写回缓存，不阻塞当前请求的返回。
			dao.CacheSet(context.Background(), cacheKey, record.Lurl, 7*24*time.Hour)
		}()

		return record.Lurl, nil
	})
	if err != nil {
		return "", err
	}

	return v.(string), nil
}
