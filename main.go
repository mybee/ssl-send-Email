package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

func main() {
	// 目前支持ssl的smtp成为主流,gmail和qqmail不再支持非ssl的端口(但163的还可用),
	// 示例是为了发送程序通用以支持更多的email提供商而开发,已测试163,
	// qq的邮箱很快而且都成功,gmail较慢有时成功有时会超时.
	host :="smtp.qq.com"
	port := 465
	email :="1144894155@qq.com"
	password :="wwvzvpyXXXXXXXxhfebh"
	toEmail :="mf_deer@163.com"

	header := make(map[string]string)
	header["From"] ="test"+"<"+ email +">"
	header["To"] = toEmail
	header["Subject"] ="邮件标题"
	header["Content-Type"] ="text/html; charset=UTF-8"

	body :="我是一封电子邮件!golang发出."

	message :=""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s", k, v)
	}
	message +=""+ body

	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)

	err := SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		[]byte(message),
	)

	if err != nil {
		panic(err)
	}
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		fmt.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
to []string, msg []byte) (err error) {

	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		fmt.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				fmt.Println("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}