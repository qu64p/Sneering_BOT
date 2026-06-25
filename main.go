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
		log.Fatal("DISCORD_TOKENが設定されていません")
	}

	if err := db.Init(); err != nil {
		log.Fatalf("Supabase初期化失敗: %v", err)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Discordセッション作成失敗: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuilds

	h := handler.New()

	dg.AddHandler(h.OnMessage)
	dg.AddHandler(h.OnInteraction)
	dg.AddHandler(h.OnReady)
	dg.AddHandler(h.OnGuildDelete)

	if err := dg.Open(); err != nil {
		log.Fatalf("Discord接続失敗: %v", err)
	}
	defer dg.Close()

	log.Println("Botが起動しました。Ctrl+Cで終了します。")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	log.Println("Botを終了します。")
}