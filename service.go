package main

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"

	core "github.com/kan-fun/kan-core"
)

type service interface {
	email(address string, subject string, body string) error
	sms(number string, code string) error
	bin(platform string) ([]string, error)
	logPut(reversedID string, content string) error
	logGetToEnd(reversedID string, fromHead bool, lastAutoID int64) ([]string, int64, error)
}

type realService struct {
}

type mockService struct {
}

func (s realService) email(address string, subject string, body string) error {
	request := requests.NewCommonRequest()
	request.Domain = "dm.aliyuncs.com"
	request.Version = "2015-11-23"
	request.ApiName = "SingleSendMail"

	request.QueryParams["AccountName"] = "no-reply@kan-fun.com"
	request.QueryParams["AddressType"] = "1"
	request.QueryParams["ReplyToAddress"] = "false"
	request.QueryParams["ToAddress"] = address
	request.QueryParams["Subject"] = subject

	if core.IsAllWhiteChar(body) {
		request.QueryParams["HtmlBody"] = "<html></html>"
	} else {
		request.QueryParams["TextBody"] = body
	}

	_, err := clientGlobal.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	return nil
}

func (s mockService) email(address string, subject string, body string) error {
	return nil
}

func (s realService) sms(number string, code string) error {
	request := requests.NewCommonRequest()
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"

	request.QueryParams["PhoneNumbers"] = number
	request.QueryParams["SignName"] = "Progress"
	request.QueryParams["TemplateCode"] = "SMS_185811363"
	request.QueryParams["TemplateParam"] = fmt.Sprintf("{\"code\":\"%s\"}", code)

	_, err := clientGlobal.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	return nil
}

func (s mockService) sms(number string, code string) error {
	return nil
}

func (s realService) bin(platform string) (result []string, err error) {
	bucketName := "kan-bin"

	bucket, err := ossClientGlobal.Bucket(bucketName)
	if err != nil {
		return
	}

	path := fmt.Sprintf("%s/", platform)

	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker), oss.Prefix(path))
		if err != nil {
			return nil, err
		}

		for _, object := range lsRes.Objects {
			if object.Key != path {
				result = append(result, object.Key)
			}
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}

	return
}

func (s mockService) bin(platform string) (result []string, err error) {
	return
}

func (s realService) logPut(reversedID string, content string) (err error) {
	putRowRequest := new(tablestore.PutRowRequest)
	putRowChange := new(tablestore.PutRowChange)
	putRowChange.TableName = "log"

	putPk := new(tablestore.PrimaryKey)
	putPk.AddPrimaryKeyColumn("reversed_id", reversedID)
	putPk.AddPrimaryKeyColumnWithAutoIncrement("auto_id")
	putRowChange.PrimaryKey = putPk

	putRowChange.AddColumn("content", content)

	putRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
	putRowRequest.PutRowChange = putRowChange

	_, err = tableStoreClientGlobal.PutRow(putRowRequest)

	return
}

func (s mockService) logPut(reversedID string, content string) (err error) {
	return
}

func (s realService) logGetToEnd(reversedID string, fromHead bool, lastAutoID int64) (result []string, newLastAutoID int64, err error) {
	newLastAutoID = lastAutoID

	getRangeRequest := &tablestore.GetRangeRequest{}
	rangeRowQueryCriteria := &tablestore.RangeRowQueryCriteria{}
	rangeRowQueryCriteria.TableName = "log"

	startPK := new(tablestore.PrimaryKey)
	startPK.AddPrimaryKeyColumn("reversed_id", reversedID)
	if fromHead {
		startPK.AddPrimaryKeyColumnWithMinValue("auto_id")
	} else {
		startPK.AddPrimaryKeyColumn("auto_id", lastAutoID+1)
	}

	endPK := new(tablestore.PrimaryKey)
	endPK.AddPrimaryKeyColumn("reversed_id", reversedID)
	endPK.AddPrimaryKeyColumnWithMaxValue("auto_id")

	rangeRowQueryCriteria.StartPrimaryKey = startPK
	rangeRowQueryCriteria.EndPrimaryKey = endPK
	rangeRowQueryCriteria.Direction = tablestore.FORWARD
	rangeRowQueryCriteria.MaxVersion = 1
	rangeRowQueryCriteria.Limit = 10
	getRangeRequest.RangeRowQueryCriteria = rangeRowQueryCriteria

	getRangeResp, err := tableStoreClientGlobal.GetRange(getRangeRequest)

	for {
		if err != nil {
			return
		}

		for _, row := range getRangeResp.Rows {
			result = append(result, row.Columns[0].Value.(string))
		}

		if len(getRangeResp.Rows) != 0 {
			newLastAutoID = getRangeResp.Rows[len(getRangeResp.Rows)-1].PrimaryKey.PrimaryKeys[1].Value.(int64)
		}

		if getRangeResp.NextStartPrimaryKey == nil {
			break
		} else {
			getRangeRequest.RangeRowQueryCriteria.StartPrimaryKey = getRangeResp.NextStartPrimaryKey
			getRangeResp, err = tableStoreClientGlobal.GetRange(getRangeRequest)
		}
	}

	return
}

func (s mockService) logGetToEnd(reversedID string, fromHead bool, lastAutoID int64) (result []string, newLastAutoID int64, err error) {
	return
}
