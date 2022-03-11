package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

func ExitOnError(message string, err error) {
	if err != nil {
		log.Printf("[%s] - %s", message, err)
		os.Exit(-1)
	}
}

func PrintOnError(message string, err error) {
	if err != nil {
		log.Printf("[%s] - %s", message, err)
	}
}

func RaiseError(message string) error {
	if !EemptyString(message) {
		return fmt.Errorf(message)
	}
	return nil
}

func EemptyString(str string) bool {
	str = strings.TrimSpace(str)
	return strings.EqualFold(str, "")
}

func UserAgentIpHash(useragent string, ip string) string {
	input := fmt.Sprintf("%s-%s-%s-%d", useragent, ip, time.Now().String(), rand.Int())
	data, _ := Sha256Of(input)
	str := Base58Encode(data)
	return str[:10]
}

func Sha256Of(input string) ([]byte, error) {
	algorithm := sha256.New()
	_, err := algorithm.Write([]byte(strings.TrimSpace(input)))
	if err != nil {
		return nil, err
	}
	return algorithm.Sum(nil), nil
}

func Base58Encode(data []byte) string {
	return base58.Encode(data)
}
