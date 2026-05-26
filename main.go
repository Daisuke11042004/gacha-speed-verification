package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// DSN（接続情報）
	dsn := "root:Pi01290502@tcp(127.0.0.1:3306)/gacha_system"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("設定エラー:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DBに接続できません:", err)
	}
	fmt.Println("MySQLへの接続に成功しました！実験（検索モード）を開始します。")

	r := gin.Default()

	// ガチャ履歴「検索」API
	r.GET("/gacha", func(c *gin.Context) {
		// 【要素分解：起点】ストップウォッチ開始
		start := time.Now()

		// 検索対象のユーザーID（101番のログを10万件から全捜索する）
		searchUserID := 101

		// SQLの検索命令（発注書：user_id が 一致するものをすべて取得）
		query := "SELECT id, user_id, item_id, rarity FROM gacha_logs WHERE user_id = ?"

		rows, err := db.Query(query, searchUserID)
		if err != nil {
			log.Println("SQL実行エラー:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
			return
		}
		defer rows.Close()

		// 取得したデータを数えるカウンター
		count := 0
		for rows.Next() {
			count++
		}

		// 【要素分解：差分】経過時間の計算
		duration := time.Since(start)

		// ターミナルにナノ秒・ミリ秒単位で速度をリアルタイム出力
		fmt.Printf("\n====================================\n")
		fmt.Printf("[実験結果] インデックスなし状態での検索\n")
		fmt.Printf("ヒット件数: %d 件\n", count)
		fmt.Printf("検索にかかった時間: %v\n", duration)
		fmt.Printf("====================================\n\n")

		c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"hit_count":    count,
			"process_time": duration.String(),
			"message":      "これがインデックスなし（目次なし）の限界速度です",
		})
	})

	r.Run(":8080")
}
