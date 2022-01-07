// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"github.com/alibabacloud-go/tea/tea"
	"os"
)

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func _main(args []*string) (_err error) {
	client, _err := CreateClient(tea.String("LTAI5tNjSPYJNT4wyTbWzYLr"), tea.String("B3R0SpewZr83NCyDVuiPlPDQwk4hYj"))
	if _err != nil {
		return _err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("阿里云短信测试"),
		TemplateCode:  tea.String("SMS_154950909"),
		PhoneNumbers:  tea.String("13011961739"),
		TemplateParam: tea.String("{\"code\":\"521520\"}"),
	}

	// 复制代码运行请自行打印 API 的返回值
	_, _err = client.SendSms(sendSmsRequest)
	if _err != nil {
		return _err
	}
	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
}
