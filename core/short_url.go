package core

import (
	"fmt"
	"math/big"
	"ohurlshortener/utils"
	"reflect"
	"time"
)

type ShortUrl struct {
	ID        int64     `db:"id" json:"id"`
	ShortUrl  string    `db:"short_url" json:"short_url"`
	DestUrl   string    `db:"dest_url" json:"desc_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Valid     bool      `db:"is_valid" json:"is_valid"`
}

func (url ShortUrl) IsEmpty() bool {
	return reflect.DeepEqual(url, ShortUrl{})
}

func GenerateShortLink(initialLink string) (string, error) {
	if utils.EemptyString(initialLink) {
		return "", fmt.Errorf("empty string")
	}
	urlHash, err := utils.Sha256Of(initialLink)
	if err != nil {
		return "", err
	}
	number := new(big.Int).SetBytes(urlHash).Uint64()
	str := utils.Base58Encode([]byte(fmt.Sprintf("%d", number)))
	return str[:8], nil
}
