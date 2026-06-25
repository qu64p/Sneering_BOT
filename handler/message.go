package handler

import (
	"log"
	"math/rand"

	"discord-bot/db"

	"github.com/bwmarrin/discordgo"
)


var cynicalReplies = []string{
  "うおw😅",
  "どわー笑😅",
  "ったく…w😅",
  "かっこよw😅",
  "えぐー！笑😅",
  "きちーww😅🤚",
  "おもろいなあ！ww😅👍",
  "ど、どした？笑😅",
  "ちょw一回落ち着けw落ち着けw😅",
  "ええてw😅🤚",
  "ええてwええてww😅🤚",
  "ちょwお前必死やんww😅",
  "すごいなあキミww😅👍",
  "お、おうww😅",
  "あぁ、そういうノリ…w😅",
  "流石にそれはうおwww😅",
  "かっけぇww😅",
  "満足した？笑😅",
  "ちょっと何言ってるか分かんないっすわww😅",
  "きっつwww😅",
  "流石に草ww😅",
  "はいはいww😅",
  "ちょw一旦深呼吸しろって笑😅",
  "あーねw😅",
  "楽しそうで何よりww😅",
  "もうええてww😅🤚",
  "わかったわかったww😅",
  "あっ、そういう感じなん？ww😅",
  "ほーん…w😅",
  "ちょwえぐいってww😅",

}


func (h *Handler) OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}
  if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply {
    return
  }
	if m.GuildID == "" {
		return
	}

	setting, err := db.GetSetting(m.GuildID)
	if err != nil {
		log.Printf("設定取得エラー (guild: %s): %v", m.GuildID, err)
		return
	}

	if !setting.Enabled {
		return
	}

	if rand.Float64()*100 >= setting.Prob {
		return
	}

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
