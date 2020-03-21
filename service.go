package main

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type Service interface {
	email(address string, subject string, body string) error
	sms(number string, code string) error
}

type RealService struct {
}

type MockService struct {
}

func (s RealService) email(address string, subject string, body string) error {
	request := requests.NewCommonRequest()
	request.Domain = "dm.aliyuncs.com"
	request.Version = "2015-11-23"
	request.ApiName = "SingleSendMail"

	request.QueryParams["AccountName"] = "no-reply@mail.progress.cool"
	request.QueryParams["AddressType"] = "1"
	request.QueryParams["ReplyToAddress"] = "false"
	request.QueryParams["ToAddress"] = address
	request.QueryParams["Subject"] = subject
	request.QueryParams["HtmlBody"] = body

	_, err := client_global.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	return nil
}

func (s MockService) email(address string, subject string, body string) error {
	return nil
}

func (s RealService) sms(number string, code string) error {
	request := requests.NewCommonRequest()
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"

	request.QueryParams["PhoneNumbers"] = number
	request.QueryParams["SignName"] = "Progress"
	request.QueryParams["TemplateCode"] = "SMS_185811363"
	request.QueryParams["TemplateParam"] = fmt.Sprintf("{\"code\":\"%s\"}", code)

	_, err := client_global.ProcessCommonRequest(request)
	if err != nil {
		return err
	}

	return nil
}

func (s MockService) sms(number string, code string) error {
	return nil
}
