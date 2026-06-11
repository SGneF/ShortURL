package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"shortURL/dao"
	"shortURL/model"
	"shortURL/pkg/base62"
	"shortURL/pkg/filter"
	"strings"
)

var (
	ErrSensitiveShortCode = errors.New("generated short code hits sensitive word, please retry")
	ErrCircularShorten    = errors.New("circular shorten is not allowed")
)

type ShortenResult struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func Shorten(longURL string, domain string) (*ShortenResult, error) {
	longURL = strings.TrimSpace(longURL)
	if longURL == "" {
		return nil, errors.New("long_url is empty")
	}

	if strings.Contains(longURL, domain) {
		return nil, ErrCircularShorten
	}

	hash := md5.Sum([]byte(longURL))
	md5Str := hex.EncodeToString(hash[:])

	existing, err := dao.FindByMd5(md5Str)
	if err == nil && existing != nil {
		return &ShortenResult{
			ShortURL: domain + "/" + existing.Surl,
			LongURL:  longURL,
		}, nil
	}

	id, err := dao.NextSequenceID()
	if err != nil {
		return nil, err
	}

	shortCode := base62.Encode(id)

	if filter.IsSensitive(shortCode) {
		id, err = dao.NextSequenceID()
		if err != nil {
			return nil, err
		}
		shortCode = base62.Encode(id)
		if filter.IsSensitive(shortCode) {
			return nil, ErrSensitiveShortCode
		}
	}

	m := &model.ShortURLMap{
		Lurl: longURL,
		Md5:  md5Str,
		Surl: shortCode,
	}
	if err := dao.CreateShortURL(m); err != nil {
		return nil, err
	}

	return &ShortenResult{
		ShortURL: domain + "/" + shortCode,
		LongURL:  longURL,
	}, nil
}
