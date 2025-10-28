# mcp_7

 Version: 0.9.1

 date    : 2025/10/02 

 update :

***

GoLang  MCP Server , Excel output example

* Turuso database use
* go version go1.24.4 

***
### setup
* config/config.go
* TURSO_DATABASE_URL, TURSO_AUTH_TOKEN set

```
const TURSO_DATABASE_URL = ""
const TURSO_AUTH_TOKEN = ""
```

***
* TABLE: table.sql
***
* start

```
go mod init example.com/go-mcp-server-7
go mod tidy

go get github.com/xuri/excelize/v2
go get github.com/joho/godotenv
go get github.com/tursodatabase/libsql-client-go/libsql


go build
```
***
### Test

* test-code: test_list.js

***
### blog

https://zenn.dev/link/comments/eab8c506730d88

***