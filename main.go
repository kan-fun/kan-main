package main

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	mrand "math/rand"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var aliyunRegionId string
var aliyunAccessKey string
var aliyunSecretKey string
var secretKeyStr string

var db *gorm.DB

var client_global *sdk.Client
var service_global Service

var privateKey_global *rsa.PrivateKey
var secretKey_global []byte

func connectDB(test bool) (db *gorm.DB, err error) {
	if test {
		db, err = gorm.Open("sqlite3", "DB.db")
	} else {
		user, ok := os.LookupEnv("WP_RDS_ACCOUNT_NAME")
		if !ok {
			return nil, errors.New("WP_RDS_ACCOUNT_NAME not set")
		}

		password, ok := os.LookupEnv("WP_RDS_ACCOUNT_PASSWORD")
		if !ok {
			return nil, errors.New("WP_RDS_ACCOUNT_PASSWORD not set")
		}

		address, ok := os.LookupEnv("WP_RDS_CONNECTION_ADDRESS")
		if !ok {
			return nil, errors.New("WP_RDS_CONNECTION_ADDRESS not set")
		}

		dbname, ok := os.LookupEnv("WP_RDS_DATABASE")
		if !ok {
			return nil, errors.New("WP_RDS_DATABASE not set")
		}

		dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, address, dbname)
		db, err = gorm.Open("mysql", dsn)
	}

	return
}

func setup(test bool) {
	if !test {
		aliyunRegionId_local, ok := os.LookupEnv("KAN_ALIYUN_REGION_ID")
		if !ok {
			panic("KAN_ALIYUN_REGION_ID not set")
		}
		aliyunRegionId = aliyunRegionId_local

		aliyunAccessKey_local, ok := os.LookupEnv("KAN_ALIYUN_ACCESS_KEY")
		if !ok {
			panic("KAN_ALIYUN_ACCESS_KEY not set")
		}
		aliyunAccessKey = aliyunAccessKey_local

		aliyunSecretKey_local, ok := os.LookupEnv("KAN_ALIYUN_SECRET_KEY")
		if !ok {
			panic("KAN_ALIYUN_SECRET_KEY not set")
		}
		aliyunSecretKey = aliyunSecretKey_local

		secretKeyStr_local, ok := os.LookupEnv("KAN_SECRET_KEY_STR")
		if !ok {
			panic("KAN_SECRET_KEY_STR not set")
		}
		secretKeyStr = secretKeyStr_local
	}

	mrand.Seed(time.Now().UnixNano())

	// Set Private Key
	privateKey_local, err := getPrivateKey(test)
	if err != nil {
		panic(err)
	}

	privateKey_global = privateKey_local

	secretKey_local, err := base64.URLEncoding.DecodeString(secretKeyStr)
	if err != nil {
		panic(err)
	}

	secretKey_global = secretKey_local

	// Connect DB
	db_local, err := connectDB(test)
	if err != nil {
		panic(err)
	}

	db = db_local
	autoMigrate()

	// Init Aliyun Client
	client_local, err := sdk.NewClientWithAccessKey(aliyunRegionId, aliyunAccessKey, aliyunSecretKey)
	if err != nil {
		panic(err)
	}

	client_global = client_local
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Set Route
	r.POST("/signup", signup)
	r.POST("/login", login)
	r.POST("/view-key", viewKey)
	r.POST("/send-email", sendEmail)
	r.POST("/send-email-code", sendEmailCode)
	r.POST("/send-sms-code", sendSMSCode)

	return r
}

func main() {
	setup(false)
	service_global = RealService{}
	r := setupRouter()
	r.Run()
}
