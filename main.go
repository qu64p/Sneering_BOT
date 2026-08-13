package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"discord-bot/db"
	"discord-bot/handler"

	"github.com/bwmarrin/discordgo"
)


func startHealthServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	go func() {
		log.Printf("ヘルスチェックサーバー起動: :%s/health", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("ヘルスチェックサーバー起動失敗: %v", err)
		}
	}()
}

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

	dg.LogLevel = discordgo.LogError

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


	startHealthServer()

	log.Println("Botが起動しました。")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	log.Println("Botを終了します。")
}