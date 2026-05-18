package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// GuildSetting はサーバーごとの設定を表す
type GuildSetting struct {
	GuildID  string
	Enabled  bool
	Prob     float64 // 確率 (%)
}

// Init はデータベース接続を初期化し、テーブルを作成する
func Init() error {
	dsn := os.Getenv("SUPABASE_DB_URL")
	if dsn == "" {
		return fmt.Errorf("SUPABASE_DB_URL が設定されていません")
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("DB接続失敗: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("DB ping 失敗: %w", err)
	}

	if err := migrate(); err != nil {
		return fmt.Errorf("マイグレーション失敗: %w", err)
	}

	log.Println("Supabase に接続しました")
	return nil
}

// migrate はテーブルが存在しない場合に作成する
func migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS guild_settings (
		guild_id  TEXT PRIMARY KEY,
		enabled   BOOLEAN NOT NULL DEFAULT TRUE,
		prob      DOUBLE PRECISION NOT NULL DEFAULT 10.0
	);`
	_, err := DB.Exec(query)
	return err
}

// GetSetting はサーバーの設定を取得する。存在しない場合はデフォルト設定を返す。
func GetSetting(guildID string) (*GuildSetting, error) {
	row := DB.QueryRow(
		`SELECT guild_id, enabled, prob FROM guild_settings WHERE guild_id = $1`,
		guildID,
	)

	s := &GuildSetting{}
	err := row.Scan(&s.GuildID, &s.Enabled, &s.Prob)
	if err == sql.ErrNoRows {
		// 初回アクセス時にデフォルト設定を挿入
		s = &GuildSetting{GuildID: guildID, Enabled: true, Prob: 10.0}
		if err := upsertSetting(s); err != nil {
			return nil, err
		}
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("設定取得失敗: %w", err)
	}
	return s, nil
}

// upsertSetting は設定を挿入または更新する
func upsertSetting(s *GuildSetting) error {
	_, err := DB.Exec(`
		INSERT INTO guild_settings (guild_id, enabled, prob)
		VALUES ($1, $2, $3)
		ON CONFLICT (guild_id) DO UPDATE
		SET enabled = EXCLUDED.enabled,
		    prob    = EXCLUDED.prob
	`, s.GuildID, s.Enabled, s.Prob)
	if err != nil {
		return fmt.Errorf("設定保存失敗: %w", err)
	}
	return nil
}

// SetEnabled は検知ON/OFFを切り替える
func SetEnabled(guildID string, enabled bool) error {
	s, err := GetSetting(guildID)
	if err != nil {
		return err
	}
	s.Enabled = enabled
	return upsertSetting(s)
}

// SetProb は確率を更新する
func SetProb(guildID string, prob float64) error {
	s, err := GetSetting(guildID)
	if err != nil {
		return err
	}
	s.Prob = prob
	return upsertSetting(s)
}
