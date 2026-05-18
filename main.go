package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"discord-bot/db"
	"discord-bot/handler"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN が設定されていません")
	}

	// Supabase 初期化
	if err := db.Init(); err != nil {
		log.Fatalf("Supabase 初期化失敗: %v", err)
	}

	// Discord セッション作成
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Discord セッション作成失敗: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	h := handler.New()

	// イベントハンドラ登録
	dg.AddHandler(h.OnMessage)
	dg.AddHandler(h.OnInteraction)
	dg.AddHandler(h.OnReady)

	if err := dg.Open(); err != nil {
		log.Fatalf("Discord 接続失敗: %v", err)
	}
	defer dg.Close()

	log.Println("Bot が起動しました。Ctrl+C で終了します。")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	log.Println("Bot を終了します。")
}
