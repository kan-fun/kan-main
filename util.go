package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	mrand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	sign "github.com/kan-fun/kan-core"
	"github.com/kan-fun/kan-server-core/model"
)

func autoMigrate() {
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.ChannelEmail{})
}

type codeClaims struct {
	CodeHash  string `json:"code_hash"`
	ChannelID string `json:"channel_id"`
	jwt.StandardClaims
}

type idClaims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func generateKey() (string, error) {
	randomBytes := make([]byte, 32)

	if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
		return "", errors.New("Can't generate key")
	}

	key := base64.URLEncoding.EncodeToString(randomBytes)

	return key, nil
}

func getPrivateKey(test bool) (*rsa.PrivateKey, error) {
	if test {
		reader := rand.Reader
		bitSize := 512

		return rsa.GenerateKey(reader, bitSize)
	} else {
		url, ok := os.LookupEnv("KAN_PRIVATE_KEY_URL")
		if !ok {
			return nil, errors.New("KAN_PRIVATE_KEY_URL not set")
		}

		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
		if err != nil {
			return nil, err
		}

		return privateKey, nil
	}
}

func generateCode(channelID string) (raw string, token string, err error) {
	ints := make([]string, 6)

	for i := 0; i <= 5; i++ {
		v := mrand.Intn(10)
		ints[i] = strconv.Itoa(v)
	}

	raw = strings.Join(ints, "")
	hash := sign.HashString(raw, secretKeyGlobal)

	token, err = generateCodeToken(hash, channelID)
	if err != nil {
		return "", "", err
	}

	return
}

func generateIDToken(id string) (tokenString string, err error) {
	claims := idClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "kan-fun.com",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err = token.SignedString(privateKeyGlobal)
	if err != nil {
		return "", err
	}

	return
}

func generateCodeToken(codeHash string, channelID string) (tokenString string, err error) {
	claims := codeClaims{
		codeHash,
		channelID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "kan-fun.com",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err = token.SignedString(privateKeyGlobal)
	if err != nil {
		return "", err
	}

	return
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func checkSignature(c *gin.Context, specificParameter map[string]string) (*model.User, error) {
	accessKey := c.GetHeader("X-Ca-Key")
	if accessKey == "" {
		return nil, errors.New("No AccessKey")
	}

	signature := c.GetHeader("X-Ca-Signature")
	if signature == "" {
		return nil, errors.New("No Signature")
	}

	signatureNonce := c.GetHeader("X-Ca-Nonce")
	if signatureNonce == "" {
		return nil, errors.New("No SignatureNonce")
	}

	timestamp := c.GetHeader("X-Ca-Timestamp")
	if timestamp == "" {
		return nil, errors.New("No Timestamp")
	}

	commonParameter := sign.CommonParameter{
		AccessKey:      accessKey,
		SignatureNonce: signatureNonce,
		Timestamp:      timestamp,
	}

	var user model.User
	db.Select("id, secret_key").Where("access_key = ?", accessKey).First(&user)
	if user.ID == 0 {
		return nil, errors.New("User not Exist")
	}

	credential, err := sign.NewCredential(accessKey, user.SecretKey)
	if err != nil {
		return nil, err
	}

	s := credential.Sign(commonParameter, specificParameter)
	if s != signature {
		return nil, errors.New("Signature not Valid")
	}

	return &user, nil
}
