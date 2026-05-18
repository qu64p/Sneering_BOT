package handler

import (
	"fmt"
	"log"

	"discord-bot/db"

	"github.com/bwmarrin/discordgo"
)

// Handler はイベントハンドラをまとめる構造体
type Handler struct{}

// New は Handler を生成する
func New() *Handler {
	return &Handler{}
}

// スラッシュコマンドの定義
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "prob",
		Description: "冷笑が発生する確率を設定します (0.1〜100.0 %)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "value",
				Description: "確率 (%) を入力 (0.1〜100.0)",
				Required:    true,
				MinValue:    floatPtr(0.1),
				MaxValue:    100.0,
			},
		},
	},
	{
		Name:        "switch",
		Description: "冷笑の検知をON/OFFで切り替えます",
	},
	{
		Name:        "nowsetting",
		Description: "現在の設定（ON/OFF・確率）を確認します",
	},
}

func floatPtr(v float64) *float64 { return &v }

// OnReady はBot起動時にスラッシュコマンドを登録する
func (h *Handler) OnReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Bot 起動完了: %s#%s", r.User.Username, r.User.Discriminator)

	// グローバルコマンドとして登録
	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			log.Printf("コマンド登録失敗 (%s): %v", cmd.Name, err)
		} else {
			log.Printf("コマンド登録成功: /%s", cmd.Name)
		}
	}
}

// OnInteraction はスラッシュコマンドのインタラクションを処理する
func (h *Handler) OnInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// DM からの実行を防ぐ
	if i.GuildID == "" {
		respond(s, i, "このコマンドはサーバー内でのみ使用できます。")
		return
	}

	switch i.ApplicationCommandData().Name {
	case "prob":
		handleProb(s, i)
	case "switch":
		handleSwitch(s, i)
	case "nowsetting":
		handleNowSetting(s, i)
	}
}

// handleProb は確率を設定する
func handleProb(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	prob := options[0].FloatValue()

	if err := db.SetProb(i.GuildID, prob); err != nil {
		log.Printf("確率設定エラー: %v", err)
		respond(s, i, "❌ 確率の設定に失敗しました。")
		return
	}

	respond(s, i, fmt.Sprintf("✅ 冷笑確率を **%.1f%%** に設定しました。", prob))
}

// handleSwitch はON/OFFを切り替える
func handleSwitch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	setting, err := db.GetSetting(i.GuildID)
	if err != nil {
		log.Printf("設定取得エラー: %v", err)
		respond(s, i, "❌ 設定の取得に失敗しました。")
		return
	}

	newEnabled := !setting.Enabled
	if err := db.SetEnabled(i.GuildID, newEnabled); err != nil {
		log.Printf("ON/OFF切替エラー: %v", err)
		respond(s, i, "❌ 設定の切り替えに失敗しました。")
		return
	}

	status := "🟢 ON"
	if !newEnabled {
		status = "🔴 OFF"
	}
	respond(s, i, fmt.Sprintf("✅ 冷笑検知を **%s** に切り替えました。", status))
}

// handleNowSetting は現在の設定を表示する
func handleNowSetting(s *discordgo.Session, i *discordgo.InteractionCreate) {
	setting, err := db.GetSetting(i.GuildID)
	if err != nil {
		log.Printf("設定取得エラー: %v", err)
		respond(s, i, "❌ 設定の取得に失敗しました。")
		return
	}

	status := "🟢 ON"
	if !setting.Enabled {
		status = "🔴 OFF"
	}

	msg := fmt.Sprintf(
		"📊 **現在の設定**\n検知: %s\n冷笑確率: **%.1f%%**",
		status,
		setting.Prob,
	)
	respond(s, i, msg)
}

// respond はエフェメラル（自分にだけ見える）なレスポンスを返す
func respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("レスポンス送信エラー: %v", err)
	}
}
