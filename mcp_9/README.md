# mcp_9

 Version: 0.9.1

 date    : 2025/10/02 

 update :

***

GoLang  MCP Server, postgres

* go version go1.25.3 linux/amd64

***
### setup
* config/config.go
* postgres set

```
const (
	Host     = "localhost"
	Port     = 5432
	User     = "user1"
	Password = "pass"
	Dbname   = "postgres"
)
```

***
* TABLE: table.sql
***
* start

```
go mod init example.com/go-mcp-server-9
go mod tidy

go get github.com/joho/godotenv
go get github.com/lib/pq

go build
```
***
### Test

* test-code: test_create.js

```
import { spawn } from "child_process";

class RpcClient {
  constructor(command) {
    this.proc = spawn(command);
    this.idCounter = 1;
    this.pending = new Map();

    this.proc.stdout.setEncoding("utf8");
    this.proc.stdout.on("data", (data) => this._handleData(data));
    this.proc.stderr.on("data", (err) => console.error("Rust stderr:", err.toString()));
    this.proc.on("exit", (code) => console.log(`Rust server exited (${code})`));
  }

  _handleData(data) {
    // 複数行対応
    data.split("\n").forEach((line) => {
      console.log("line=", line);
      if (!line.trim()) return;
      try {
        const msg = JSON.parse(line);
        if (msg.id && this.pending.has(msg.id)) {
          const { resolve } = this.pending.get(msg.id);
          this.pending.delete(msg.id);
          resolve(msg.result);
        }
      } catch (e) {
        //console.error("JSON parse error:", e, line);
      }
    });
  }

  call(method, params = {}) {
    const id = this.idCounter++;
    const payload = {
      jsonrpc: "2.0",
      id,
      method,
      params,
    };

    return new Promise((resolve, reject) => {
      this.pending.set(id, { resolve, reject });
      this.proc.stdin.write(JSON.stringify(payload) + "\n");
    });
  }

  close() {
    this.proc.kill();
  }
}

// -----------------------------
// 実行例
// -----------------------------
async function main() {
  const client = new RpcClient("/home/naka/work/go/mcp/mcp_9/go-mcp-server-9");

  const result1 = await client.call(
    "tools/call", 
    { 
      name: "test_create", 
      arguments: {title: "tit2", content: "cc2"}, 
    },
  );
  console.log("add結果:", result1);

  client.close();
}

main().catch(console.error);

```
***
### blog

***