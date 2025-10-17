# mcp_6

 Version: 0.9.1

 date    : 2025/10/16

***

GoLang remoto MCP Server , TURSO database

***
### setup

* .env

```
TURSO_DATABASE_URL = 
TURSO_AUTH_TOKEN = 
```

***
* build

```
go mod init example.com/go-remote-mcp-server6
go mod tidy

go get github.com/joho/godotenv
go get github.com/tursodatabase/libsql-client-go/libsql


go build
go run .

```
***
* setting.json , GEMINI-CLI

```
    "myRemoteServer": {
      "httpUrl": "http://localhost:8080/mcp", 
      "headers": {
        "Authorization": "123" 
      },
      "timeout": 5000 
    }        

```

***
* prompt

```
お茶 , 130　円を購入。品名、価格 の値をAPIに送信して欲しい。
```
***

