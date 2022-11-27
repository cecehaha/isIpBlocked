package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/jordan-wright/email"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
)

type IPCheckResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Ping         bool   `json:"ping"`
		Tcp          bool   `json:"tcp"`
		Ip           string `json:"ip"`
		CountryClode string `json:"countryClode"`
	} `json:"data"`
}

var (
	HOST       string
	PORT       string
	foreignURL string
	homeURL    string
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func sendMail(title string, message string) (err error) {
	emailsTo := os.Getenv("EMAIL_TO")
	if emailsTo == "" {
		log.Printf("EMAIL is empty, skip sending email")
		return
	}
	var to []string
	for _, toAddr := range regexp.MustCompile(`\s+`).Split(emailsTo, -1) {
		to = append(to, toAddr)
	}
	// 发送邮件
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		log.Printf("EMAIL_FROM or EMAIL_PASSWORD or SMTP_HOST or SMTP_PORT is empty, skip sending email")
		return
	}
	SmtpAuth := os.Getenv("SMTP_AUTH")
	var auth smtp.Auth
	if SmtpAuth == "plain" {
		auth = smtp.PlainAuth("", from, password, smtpHost)
	} else if SmtpAuth == "login" {
		auth = LoginAuth(from, password)
	} else if SmtpAuth == "crammd5" {
		auth = smtp.CRAMMD5Auth(from, password)
	} else {
		return fmt.Errorf("SMTP_AUTH is not plain or login")
	}

	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = title
	e.Text = []byte(message)

	err = e.Send(smtpHost+":"+smtpPort, auth)
	if err != nil {
		err = fmt.Errorf("sendMail failed, %s", err)
		return
	}
	log.Println("sendMail success")
	return nil
}

func sendTg(title string, message string) (err error) {
	// 通知
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	tgChatID := os.Getenv("TG_CHAT_ID")
	if tgBotToken == "" || tgChatID == "" {
		log.Printf("TG_BOT_TOKEN or TG_CHAT_ID is empty, skip sending telegram message")
	} else {
		// POST 发送tg消息
		tgURL := fmt.Sprintf(
			"https://api.telegram.org/bot%s/sendMessage?chat_id=%s",
			tgBotToken,
			tgChatID,
		)
		data := fmt.Sprintf(`{"text":"%s\n\n%s"}`, title, message)
		var resp *http.Response
		resp, err = http.Post(tgURL, "application/json", io.NopCloser(strings.NewReader(data)))
		if err != nil {
			err = fmt.Errorf("sendTg failed, %s", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				log.Printf("sendTg failed, %s", err)
			}
		}(resp.Body)
	}
	log.Println("sendTg success")
	return
}

func notify(title string, message string) (err error) {
	// email 通知
	sendEmailErr := sendMail(title, message)
	// tg 通知
	sendTgErr := sendTg(title, message)
	// 检查错误
	errString := ""
	if sendEmailErr != nil {
		errString = sendEmailErr.Error()
	}
	if sendTgErr != nil {
		errString += sendTgErr.Error()
	}
	if errString != "" {
		err = fmt.Errorf(errString)
		return
	} else {
		return nil
	}
}

func isIP(host string) bool {
	// 正则检查HOST地址是否为IP
	pattern := `^((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)$`
	isReg := regexp.MustCompile(pattern).MatchString(host)
	return isReg
}

func isDomain(host string) bool {
	// 正则检查HOST地址是否为域名
	pattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`
	isReg := regexp.MustCompile(pattern).MatchString(host)
	return isReg
}

func isPort(port string) bool {
	// 正则检查PORT是否为范围在 1-65535 的整数
	pattern := `^([1-9]|[1-9]\d{1,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`
	isReg := regexp.MustCompile(pattern).MatchString(port)
	return isReg
}

func initURL() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("load .env failed, %s", err)
	}
	HOST = os.Getenv("HOST")
	PORT = os.Getenv("PORT")
	// 检查IP和PORT格式
	if HOST == "" {
		panic("IP is empty")
	}
	if PORT == "" {
		PORT = "22"
	}

	// 正则检查HOST地址和PORT合法性
	if (isIP(HOST) || isDomain(HOST)) && isPort(PORT) {
		// 格式化 URL
		homeURL = fmt.Sprintf("https://api.24kplus.com/ipcheck?host=%s&port=%s", HOST, PORT)
		foreignURL = fmt.Sprintf("https://api.idcoffer.com/ipcheck?host=%s&port=%s", HOST, PORT)
		return
	}
	panic("IP or PORT is invalid")
}

func IPCheck(apiURL string) (result *IPCheckResult, err error) {
	// 发送请求
	resp, err := http.Get(apiURL)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	// 解析返回值
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return
	}
	if result.Code == 1 {
		return
	}
	return nil, fmt.Errorf("IPCheck failed, code: %d, message: %s", result.Code, result.Message)
}

func main() {
	// 检查IP是否被封禁
	initURL()
	homeResult, err := IPCheck(homeURL)
	if err != nil {
		if err = notify("IPCheck failed", err.Error()); err != nil {
			panic(err)
		}
		return
	}
	homePing, homeTcp := homeResult.Data.Ping, homeResult.Data.Tcp

	var globalResult *IPCheckResult
	globalResult, err = IPCheck(foreignURL)
	if err != nil {
		if err = notify("IPCheck failed", err.Error()); err != nil {
			panic(err)
		}
		return
	}
	globalPing, globalTcp := globalResult.Data.Ping, globalResult.Data.Tcp

	var title string
	message := fmt.Sprintf(
		"IP: %s\nPORT: %s\n%s 节点结果：Ping: %t tcp: %t\n%s 节点结果：Ping: %t tcp: %t",
		HOST, PORT,
		homeResult.Data.CountryClode, homePing, homeTcp,
		globalResult.Data.CountryClode, globalPing, globalTcp,
	)

	if !homePing && !homeTcp && !globalPing && !globalTcp {
		title = fmt.Sprintf("Host %s is Down", HOST)
		message += "\n无法判断是否被阻断"
		if err = notify(title, message); err != nil {
			panic(err)
		}
		return
	}
	if !homePing && !homeTcp {
		title = fmt.Sprintf("Host %s is blocked", HOST)
		if globalPing && globalTcp {
			message += "\n国外节点可访问，IP被完全阻断"
		} else if globalPing && !globalTcp {
			message += "\n国外节点可ping，IP被icmp阻断，无法判断是否被tcp阻断"
		} else if !globalPing && globalTcp {
			message += "\n国外节点可tcping，IP被tcp阻断，可能是被封端口，无法判断是否被icmp阻断"
		}
		if err = notify(title, message); err != nil {
			panic(err)
		}
	}
	if !homePing && homeTcp {
		title = fmt.Sprintf("Host %s is icmp blocked", HOST)
		if globalPing {
			message += "\n国外节点可ping，IP被icmp阻断，无法判断是否被tcp阻断"
		} else {
			message += "\n国外节点不可ping，无法判断是否被阻断，端口未被tcp阻断"
			fmt.Println(message)
			return
		}
		if err = notify(title, message); err != nil {
			panic(err)
		}
	}
	if homePing && !homeTcp {
		title = fmt.Sprintf("Host %s is tcp blocked", HOST)
		if globalTcp {
			message += "\n国外节点可tcping，IP被tcp阻断，可能是被封端口，无法判断是否被icmp阻断"
		} else {
			message += "\n国外节点不可tcping，无法判断是否被阻断"
			fmt.Println(message)
			return
		}
		if err = notify(title, message); err != nil {
			panic(err)
		}
	}
	if homePing && homeTcp {
		log.Printf("Host %s is not blocked", HOST)
	}
	fmt.Println(message)

	// Output:
	// 国内ping结果，国内tcp结果，国外ping结果，国外tcp结果
	// 可能被阻断的情况下，发送邮件/tg通知
}
