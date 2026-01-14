package mail

import (
	"errors"
	"fmt"
	"mime"

	"gopkg.in/gomail.v2"
)

func Message(options ...Option) (*gomail.Message, error) {
	e := loadOptions(options...)

	if len(e.mailTo) <= 0 {
		return nil, errors.New("收件人为空")
	}

	//for _, v := range e.mailTo {
	//	if !e.VerifyEmailFormat(v.address) {
	//		return nil, fmt.Errorf("收件人中有错误邮件地址 `%s`", v.address)
	//	}
	//}

	m := gomail.NewMessage(
		gomail.SetEncoding(gomail.Base64),
	)

	//这种方式可以添加别名，即“XX官方”
	if e.from.address != "" {
		m.SetHeader("From", m.FormatAddress(e.from.address, e.from.name))
	}
	// 说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])

	var toArray []string
	for _, v := range e.mailTo {
		toArray = append(toArray, m.FormatAddress(v.address, v.name))
	}
	m.SetHeader("To", toArray...) // 发送给多个用户

	var ccArray []string
	for _, v := range e.mailCc {
		ccArray = append(ccArray, m.FormatAddress(v.address, v.name))
	}
	if len(ccArray) > 0 {
		m.SetHeader("Cc", ccArray...) // 抄送
	}

	var bccArray []string
	for _, v := range e.mailBcc {
		bccArray = append(bccArray, m.FormatAddress(v.address, v.name))
	}
	if len(ccArray) > 0 {
		m.SetHeader("Bcc", bccArray...) // 暗送
	}

	m.SetHeader("Subject", e.subject) //设置邮件主题
	m.SetBody("text/html", e.body)    //设置邮件正文
	if len(e.mailAttach) > 0 {
		for _, attach := range e.mailAttach {
			if attach.name == "" || attach.path == "" {
				continue
			}
			m.Attach(attach.path,
				gomail.Rename(attach.name),
				gomail.SetHeader(map[string][]string{
					"Content-Disposition": {
						fmt.Sprintf(`attachment; filename="%s"`, mime.BEncoding.Encode("UTF-8", attach.name)),
					},
				}),
			)
		}
	}
	return m, nil
}
