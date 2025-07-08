package main		// mainパッケージであることを宣言

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	// --- 送信者と受信者の情報を設定 ---
	// 送信者のGmailアドレス
	from := "your_email@gmail.com"
	// 上記で取得したアプリパスワード
	password := "your_app_password" // 16桁のアプリパスワード

	// 受信者のメールアドレス（複数指定可能）
	to := []string{
		"recipient1@example.com",
	}

	// --- SMTPサーバーの情報を設定 ---
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// --- メールの内容を作成 ---
	// メールの件名と本文
	subject := "Subject: Goからのテストメール\n"
	body := "これはGo言語から送信されたテストメールです。"
	// メッセージ全体を組み立てる
	message := []byte(subject + body)

	// --- 認証情報を作成 ---
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// --- メールを送信 ---
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatal(err) // エラーが発生した場合はログに出力して終了
	}

	fmt.Println("メールが正常に送信されました！")
}