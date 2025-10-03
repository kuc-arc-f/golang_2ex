# mcp_3

 Version: 0.9.1

 date    : 2025/10/02 

 update :

***

goLang + Turuso database , MCP Server

* GEMINI-CLI use
***
* TABLE: table.sql
***
* start

```
go mod init example.com/go-mcp-server-3
go mod tidy

go build
```
***
* settings.json , GEMINI-CLI
```
  "mcpServers": {
    "go-mcp-server-3": {
      "command": "/path/mcp_3/go-mcp-server-3.exe",
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