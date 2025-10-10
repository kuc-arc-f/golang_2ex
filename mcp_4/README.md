# mcp_4

 Version: 0.9.1

 date    : 2025/10/10 

 update :

***

goLang + Turuso database , MCP Server

* GEMINI-CLI use
* Turso SDK GoLang
***
* TABLE: table.sql
***
* start

```
go mod init example.com/go-mcp-server-4
go mod tidy

go get github.com/tursodatabase/libsql-client-go/libsql
gi build
```
***
* settings.json , GEMINI-CLI
```
  "mcpServers": {
    "go-mcp-server-4": {
      "command": "/work/go/mcp/mcp_4/go-mcp-server-4.exe",
      "env": {
        "HOGE": ""
      }
    }
  },
```
***
### Prompt

```
コーヒー , 170　円を購入。品名、価格 の値をAPIに送信して欲しい。
```
***
### blog


***