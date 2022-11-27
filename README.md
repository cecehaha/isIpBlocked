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
下载配置文件到 `isIpBlocked` 同目录：
```bash
wget https://raw.githubusercontent.com/cecehaha/isIpBlocked/main/.env.example && cp .env.example .env
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

# 运行
直接运行来测试一下：
```bash
./isIpBlocked
```

设置定时任务：
```bash
crontab -e
```

之后编辑文件加入（每五分钟执行一次）：
```bash
*/5 * * * * /root/isIpBlocked
```

注意：`/root/isIpBlocked` 需要是 `isIpBlocked` 的绝对路径

# 说明
被墙主要有：
1. 端口被墙：表现为某些端口（比如80，443）的国外节点 tcp 通但国内 tcp 不通，如果服务器没有禁ping的话国内外都能ping通
2. 整个IP被墙：表现为国内ping tcp都不通，国外都通

当然这个检测需要你本身没有禁ping和开放了tcp端口。
