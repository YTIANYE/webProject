package handler

import (
	"context"
	"fmt"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"math/rand"
	"time"
	"webProject/service/user/model"
	"webProject/service/user/utils"

	"github.com/micro/go-micro/util/log"

	user "webProject/service/user/proto/user"
)

type User struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) SendSms(ctx context.Context, req *user.Request, rsp *user.Response) error {
	log.Log("Received User.SendSms request")

	// 校验图片验证码 是否正确
	result := model.CheckImgCode(req.Uuid, req.ImgCode)
	if result {
		fmt.Println("图片验证码校验成功")

		// 发送短信
		client, _err := CreateClient(tea.String("LTAI5tNjSPYJNT4wyTbWzYLr"), tea.String("B3R0SpewZr83NCyDVuiPlPDQwk4hYj"))
		if _err != nil {
			fmt.Println("发送短信错误 _err", _err)
		}

		// 生成一个随机的六位数，作为验证码
		rand.Seed(time.Now().UnixNano()) //播种随机种子
		// 生成6位随机数
		smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
		smsCodeString := "{" + "\"code\":" + "\"" + smsCode + "\"" + "}"

		sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
			SignName:      tea.String("阿里云短信测试"),
			TemplateCode:  tea.String("SMS_154950909"),
			PhoneNumbers:  tea.String(req.Phone),
			TemplateParam: tea.String(smsCodeString),
		}

		// 打印 API 的返回值
		_, _err = client.SendSms(sendSmsRequest)
		if _err != nil {
			rsp.Errno = utils.RECODE_SMSERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_SMSERR)
			fmt.Println("发送短信错误 _err", _err)
		} else {
			// 发送短信成功
			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
			fmt.Println("发送短信成功")

			// 将 电话号：短信验证码 存入 Redis
			err := model.SaveSmsCode(req.Phone, smsCode)
			if err != nil {
				rsp.Errno = utils.RECODE_DBWRITERR
				rsp.Errmsg = utils.RecodeText(utils.RECODE_DBWRITERR)
				fmt.Println("存储短信验证码到redis失败：", err)
			} else {
				fmt.Println("存储短信验证码到redis成功")
			}
		}

	} else {
		// 校验失败 发送错误信息
		rsp.Errno = utils.RECODE_CHECKERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_CHECKERR)
		fmt.Println("校验失败， 发送短信错误")
	}

	return nil
}

func (e *User) Register(ctx context.Context, req *user.RegReq, rsp *user.Response) error {
	// 校验名短信验证码是否正确（短信验证码存储在redis中）
	err := model.CheckSmsCode(req.Mobile, req.SmsCode)
	if err != nil {
		fmt.Println("短信验证码错误")
		rsp.Errno = utils.RECODE_CHECKMSMERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_CHECKMSMERR)
		return err
	}

	// 校验正确，注册用户，将数据写入MySQL数据库
	err = model.RegisterUser(req.Mobile, req.Password)
	if err != nil {
		fmt.Println("写入数据库错误")
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		return err
	} else {
		fmt.Println("写入数据库成功")
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
	}
	return nil
}
