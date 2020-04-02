package main

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type service interface {
	email(address string, subject string, body string) error
	sms(number string, code string) error
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

	if strings.TrimSpace(body) == "" {
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
