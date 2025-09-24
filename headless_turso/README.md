# headless_turso

 Version: 0.9.1

 date    : 2025/09/22 

 update :

***

GoLang Turso SDK , Headless CMS

***
### API document

https://github.com/kuc-arc-f/golang_2ex/blob/main/headless_turso/document/api.md

***
### setup
* .env
* API_KEY: API auth key

```
API_KEY=123
TURSO_DATABASE_URL=""
TURSO_AUTH_TOKEN=
```
***
* TABLE: schema.sql

```
CREATE TABLE IF NOT EXISTS todo (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  data TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
***
* start

```
go mod init example.com/headlessturso
go mod tidy

go get github.com/joho/godotenv
go get github.com/tursodatabase/libsql-client-go/libsql

go run .
```

***
### blog

***