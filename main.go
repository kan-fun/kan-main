package main

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var aliyunRegionID string
var aliyunAccessKey string
var aliyunSecretKey string
var secretKeyStr string

var db *gorm.DB

var clientGlobal *sdk.Client
var ossClientGlobal *oss.Client
var tableStoreClientGlobal *tablestore.TableStoreClient
var serviceGlobal service

var privateKeyGlobal *rsa.PrivateKey
var secretKeyGlobal []byte

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
		aliyunRegionIDLocal, ok := os.LookupEnv("KAN_ALIYUN_REGION_ID")
		if !ok {
			panic("KAN_ALIYUN_REGION_ID not set")
		}
		aliyunRegionID = aliyunRegionIDLocal

		aliyunAccessKeyLocal, ok := os.LookupEnv("KAN_ALIYUN_ACCESS_KEY")
		if !ok {
			panic("KAN_ALIYUN_ACCESS_KEY not set")
		}
		aliyunAccessKey = aliyunAccessKeyLocal

		aliyunSecretKeyLocal, ok := os.LookupEnv("KAN_ALIYUN_SECRET_KEY")
		if !ok {
			panic("KAN_ALIYUN_SECRET_KEY not set")
		}
		aliyunSecretKey = aliyunSecretKeyLocal

		secretKeyStrLocal, ok := os.LookupEnv("KAN_SECRET_KEY_STR")
		if !ok {
			panic("KAN_SECRET_KEY_STR not set")
		}
		secretKeyStr = secretKeyStrLocal
	}

	rand.Seed(time.Now().UnixNano())

	// Set Private Key
	privateKeyLocal, err := getPrivateKey(test)
	if err != nil {
		panic(err)
	}

	privateKeyGlobal = privateKeyLocal

	secretKeyLocal, err := base64.URLEncoding.DecodeString(secretKeyStr)
	if err != nil {
		panic(err)
	}

	secretKeyGlobal = secretKeyLocal

	// Connect DB
	dbLocal, err := connectDB(test)
	if err != nil {
		panic(err)
	}

	db = dbLocal
	autoMigrate()

	// Init Aliyun Client
	clientLocal, err := sdk.NewClientWithAccessKey(aliyunRegionID, aliyunAccessKey, aliyunSecretKey)
	if err != nil {
		panic(err)
	}

	clientGlobal = clientLocal

	// Init OSS Client
	ossClientLocal, err := oss.New("oss-cn-beijing.aliyuncs.com", aliyunAccessKey, aliyunSecretKey)
	if err != nil {
		panic(err)
	}

	ossClientGlobal = ossClientLocal

	// Init TableStore Client
	tableStoreClientLocal := tablestore.NewClient("https://kan.cn-beijing.ots.aliyuncs.com", "kan", aliyunAccessKey, aliyunSecretKey)
	tableStoreClientGlobal = tableStoreClientLocal
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
	r.GET("/log/pub", logPub)
	r.GET("/log/sub", logSub)
	r.GET("/bin", bin)

	return r
}

func main() {
	setup(false)
	serviceGlobal = realService{}
	r := setupRouter()
	r.Run(":8080")
}
