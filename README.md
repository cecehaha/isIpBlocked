# Linux 安装
对于 `amd64`
```bash
wget -O isIpBlocked https://github.com/cecehaha/isIpBlocked/releases/latest/download/isIpBlocked_linux_amd64 && chmod +x isIpBlocked
```

对于 `arm64`
```bash
wget -O isIpBlocked https://github.com/cecehaha/isIpBlocked/releases/latest/download/isIpBlocked_linux_arm64 && chmod +x isIpBlocked
```

> 其他架构系统看 [release](https://github.com/cecehaha/isIpBlocked/releases/latest) 中有没有，没有可自行编译

# 配置
下载配置文件：
```bash
wget -O .env https://raw.githubusercontent.com/cecehaha/isIpBlocked/main/.env.example
```

编辑配置文件：
```bash
nano .env
```

说明：
```bash
# HOST 可以用域名也可也用 IP
HOST="1.1.1.1"
# PORT 可以改成你需要的端口，或者默认80，主要是测试端口tcp可用性
PORT=80

# 不需要email或者tg通知的话，留空即可
# EMAIL_TO 填写收通知的邮箱
EMAIL_TO="to@example.com"
# EMAIL_FROM 填写发通知的邮箱
EMAIL_FROM="from@hotmail.com"
EMAIL_PASSWORD="password"
SMTP_HOST="smtp-mail.outlook.com"
SMTP_PORT=587
# auth 方法：plain, login, cram-md5
# TODO 暂时只测试了 login 方法
SMTP_AUTH=login

# telegram 通知
TG_BOT_TOKEN=""
TG_CHAT_ID=""
```

- TG_BOT_TOKEN: 可以去 [@botFather](https://t.me/botFather) 处新建机器人来获取机器人的token
- TG_CHAT_ID: 可以去 [@userinfobot](https://t.me/userinfobot) 处来获取自己的账号ID，并填入
