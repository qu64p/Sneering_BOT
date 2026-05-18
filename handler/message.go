package handler

import (
	"log"
	"math/rand"

	"discord-bot/db"

	"github.com/bwmarrin/discordgo"
)

// cynicalReplies は冷笑リプライの一覧
var cynicalReplies = []string{
	"うおw😅",
	"どわー笑😅",
	"ったく…w😅",
	"かっこよw😅",
	"えぐー！笑😅",
	"きちーww😅🤚",
	"おもろいなあ！ww😅👍",
	"ど、どした？笑😅",
	"ちょw1回落ち着けw落ち着けw😅",
	"ええてw😅🤚",
	"ええてw ええてww😅🤚",
	"ちょwお前必死やんww😅",
	"すごいなあキミww😅👍",
	"お、おうww😅",
	"あぁ、そういうノリ…w😅",
	"流石にそれはうおwww😅",
}

// OnMessage はメッセージイベントを処理する
func (h *Handler) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot 自身のメッセージは無視
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}
	// DM は無視（GuildID が空）
	if m.GuildID == "" {
		return
	}

	setting, err := db.GetSetting(m.GuildID)
	if err != nil {
		log.Printf("設定取得エラー (guild: %s): %v", m.GuildID, err)
		return
	}

	// 検知がOFFなら何もしない
	if !setting.Enabled {
		return
	}

	// 確率判定 (0.0〜100.0 %)
	if rand.Float64()*100 >= setting.Prob {
		return
	}

	// ランダムにリプライを選択
	reply := cynicalReplies[rand.Intn(len(cynicalReplies))]

	_, err = s.ChannelMessageSendReply(
		m.ChannelID,
		reply,
		m.Reference(),
	)
	if err != nil {
		log.Printf("リプライ送信エラー: %v", err)
	}
}
