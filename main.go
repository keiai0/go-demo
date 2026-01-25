package main

import (
	"context"
	"fmt"

	"github.com/keiai0/go-demo/books"
	"github.com/keiai0/go-demo/handler"
	"github.com/keiai0/go-demo/handler/bank"
	"github.com/keiai0/go-demo/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
	"net/http"
	"time"
)

func main() {
	// DB接続
	dsn := "host=db user=postgres password=postgres dbname=go_demo port=5432 sslmode=disable TimeZone=Asia/Tokyo"
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Redis接続
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// Redisの接続確認
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	// CacheHandlerの初期化
	cacheHandler := &handler.CacheHandler{
		DB:    dbConn,
		Redis: rdb,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!!!")
	})

	// 時関
	http.HandleFunc("/demo/time", handler.TimeDemoHandler)

	// アルゴリズム
	http.HandleFunc("/demo/algorithm", handler.AlgorithmDemoHandler)

	// 書籍
	http.HandleFunc("/demo/books", books.BooksDemoHandler)

	// キャッシュデモ
	// テストデータ(10,000件)を作成
	http.HandleFunc("/demo/cache/init", cacheHandler.InitTestDataHandler)
	// Redisキャッシュを使用してログ取得
	http.HandleFunc("/demo/cache/with", cacheHandler.CacheWithHandler)
	// キャッシュなしでDBから直接ログ取得
	http.HandleFunc("/demo/cache/without", cacheHandler.CacheWithoutHandler)
	// Redisキャッシュをクリア
	http.HandleFunc("/demo/cache/clear", cacheHandler.ClearCacheHandler)

	// レート制限付きAPI (1分間に10回まで)
	rateLimiter := middleware.NewRateLimiter(10, time.Minute)
	http.Handle("/api/limited", rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"message": "Success! This endpoint is rate-limited to 10 requests per minute."}`)); err != nil {
			return
		}
	})))

	// セキュリティデモAPI
	securityHandler := &handler.SecurityDemoHandler{DB: dbConn}
	// 情報エンドポイント
	http.HandleFunc("/api/security", securityHandler.SecurityInfoHandler)

	// CORS デモ
	http.HandleFunc("/api/security/cors/vulnerable", securityHandler.CORSVulnerableHandler)
	http.HandleFunc("/api/security/cors/secure", securityHandler.CORSSecureHandler)

	// CSRF デモ
	http.HandleFunc("/api/security/csrf/token", securityHandler.CSRFTokenHandler)
	http.HandleFunc("/api/security/csrf/vulnerable", securityHandler.CSRFVulnerableHandler)
	http.HandleFunc("/api/security/csrf/secure", securityHandler.CSRFSecureHandler)

	// XSS デモ
	http.HandleFunc("/api/security/xss/vulnerable", securityHandler.XSSVulnerableHandler)
	http.HandleFunc("/api/security/xss/secure", securityHandler.XSSSecureHandler)

	// SQL Injection デモ
	http.HandleFunc("/api/security/sql-injection/vulnerable", securityHandler.SQLInjectionVulnerableHandler)
	http.HandleFunc("/api/security/sql-injection/secure", securityHandler.SQLInjectionSecureHandler)

	// 並列URL取得API
	parallelFetchHandler := &handler.ParallelFetchHandler{}
	// 最速レスポンスを返す
	http.HandleFunc("/api/fetch/fastest", parallelFetchHandler.FetchFastestHandler)
	// 全結果を返す
	http.HandleFunc("/api/fetch/all", parallelFetchHandler.FetchAllHandler)

	// 複数API並列実行・集約API
	aggregateHandler := &handler.AggregateAPIHandler{}
	// カスタムAPIを並列実行して集約
	http.HandleFunc("/api/aggregate", aggregateHandler.AggregateHandler)
	// プリセットAPIを並列実行して集約（デモ用）
	http.HandleFunc("/api/aggregate/preset", aggregateHandler.PresetAggregateHandler)

	// 銀行振込API
	bankTransferHandler := &bank.Handler{DB: dbConn}
	// テスト用口座を初期化
	http.HandleFunc("/api/bank/init", bankTransferHandler.InitAccountsHandler)
	// 通常の振込処理
	http.HandleFunc("/api/bank/transfer", bankTransferHandler.NormalTransferHandler)
	// 口座情報を取得
	http.HandleFunc("/api/bank/account", bankTransferHandler.GetAccountHandler)
	// 全口座一覧を取得
	http.HandleFunc("/api/bank/accounts", bankTransferHandler.ListAccountsHandler)
	// Dirty Readデモ
	http.HandleFunc("/api/bank/dirty-read", bankTransferHandler.DirtyReadDemoHandler)
	// Phantom Readデモ
	http.HandleFunc("/api/bank/phantom-read", bankTransferHandler.PhantomReadDemoHandler)
	// デッドロックデモ
	http.HandleFunc("/api/bank/deadlock", bankTransferHandler.DeadlockDemoHandler)
	// デッドロック回避策1: ロック順序の統一
	http.HandleFunc("/api/bank/transfer-safe", bankTransferHandler.DeadlockAvoidanceHandler)
	// デッドロック回避策2: タイムアウト設定
	http.HandleFunc("/api/bank/transfer-timeout", bankTransferHandler.DeadlockTimeoutHandler)

	fmt.Println("localhost:8080 server runnig ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Todo struct {
	ID        int `gorm:"primaryKey"`
	UserID    int
	User      User
	Title     string
	DoneFlag  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type User struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Category struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type TodoCategory struct {
	ID         int `gorm:"primaryKey"`
	TodoID     int
	CategoryID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

type Todos []Todo
