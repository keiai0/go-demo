# PostgreSQL EXPLAIN デモプログラム

このプログラムは、PostgreSQLの`EXPLAIN`コマンドを使った実行計画の分析を実践的に学ぶためのデモです。

## 📚 ドキュメント

このプロジェクトには3つのレベルの学習教材があります:

1. **[クイックスタート](../../docs/explain_quickstart.md)** - 5分で始める基本
2. **[実践演習SQL](./examples.sql)** - 手を動かして学ぶ
3. **[完全ガイド](../../docs/postgresql_explain_guide.md)** - 詳細な解説

## 前提条件

- Docker がインストールされていること
- PostgreSQLコンテナが起動していること

## セットアップ

### 1. データベースの起動

```bash
# プロジェクトルートから実行
cd ../dev/private/go/go-demo

# Dockerコンテナ起動
docker compose up -d db

# スキーマとダミーデータを投入
make refresh-schema
```

### 2. 依存パッケージのインストール

```bash
cd demo/explain_demo
go mod tidy
```

## 🚀 最も簡単な実行方法

```bash
cd demo/explain_demo

# クイックデモを実行（基本的な例を自動実行）
./quick_demo.sh

# またはインタラクティブな実行スクリプト
./run_examples.sh
```

## 実行方法（詳細）

### 方法1: PostgreSQLに直接接続（推奨・最も簡単）

最もシンプルで効果的な方法です。

```bash
# PostgreSQLに接続
docker exec -it go-demo-db psql -U postgres -d go_demo

# 接続後、EXPLAINコマンドを実行
EXPLAIN ANALYZE SELECT * FROM users;

# examples.sqlの内容をコピー＆ペーストして試すこともできます
```

### 方法2: Goプログラムをコンテナ内で実行

PostgreSQLコンテナと同じネットワーク内で実行します。

```bash
cd demo/explain_demo

# Docker環境で実行
docker run --rm \
  --network go-demo_go-demo-network \
  -v $(pwd):/app \
  -w /app \
  -e DB_HOST=db \
  golang:1.23 \
  sh -c "go mod download && go run main.go"
```

または、シンプルにGoがインストールされている場合:

```bash
cd demo/explain_demo
DB_HOST=db go run main.go
```

### 方法2: ローカルから実行

ローカルのPostgreSQLクライアントライブラリを使用します。

```bash
cd demo/explain_demo

# デフォルト（localhost:5432）
go run main.go

# カスタム設定
DB_HOST=localhost DB_PORT=5432 go run main.go
```

**注意:** ローカルから実行する場合、PostgreSQLコンテナのポートが5432でホストに公開されている必要があります。

### 方法3: PostgreSQLで直接SQLを実行（最も簡単）

```bash
# PostgreSQLに接続
docker exec -it go-demo-db psql -U postgres -d go_demo

# examples.sqlの内容をコピー＆ペースト
# または直接EXPLAINコマンドを実行
EXPLAIN ANALYZE SELECT * FROM users;
```

このプログラムは以下の例を実行します:

1. **基本的なSELECT** - Seq Scanの確認
2. **WHERE句でのフィルタリング** - フィルタ条件の影響
3. **インデックスの効果** - インデックスあり/なしの比較
4. **JOIN操作** - 結合方法の確認
5. **集約クエリ** - GROUP BYとORDER BYの実行計画
6. **大量データ** - データ量による違い
7. **複合インデックス** - 複数カラムのインデックス活用

## 手動でEXPLAINを試す

### PostgreSQLに接続

```bash
# コンテナ内のpsqlに接続
docker exec -it go-demo-db psql -U postgres -d go_demo
```

### 基本的なクエリ

```sql
-- 実行計画のみ表示
EXPLAIN SELECT * FROM users;

-- 実際に実行して実測値も表示
EXPLAIN ANALYZE SELECT * FROM users;

-- 詳細情報を表示
EXPLAIN (ANALYZE, BUFFERS, VERBOSE) SELECT * FROM users;
```

### インデックスの作成と効果確認

```sql
-- インデックスを作成
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- インデックスありでの検索
EXPLAIN ANALYZE SELECT * FROM todos WHERE user_id = 1;

-- インデックスの一覧確認
SELECT tablename, indexname, indexdef
FROM pg_indexes
WHERE tablename = 'todos';
```

### 大量データでの実験

```sql
-- 10万件のデータを追加
INSERT INTO todos (user_id, title, note, done_flag)
SELECT
    (random() * 3 + 1)::int,  -- user_id 1-4
    'TODO ' || generate_series,
    'Note ' || generate_series,
    random() > 0.5
FROM generate_series(1, 100000);

-- インデックスなしでの検索時間を確認
EXPLAIN ANALYZE SELECT * FROM todos WHERE user_id = 2;

-- インデックス作成
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- インデックスありでの検索時間を確認
EXPLAIN ANALYZE SELECT * FROM todos WHERE user_id = 2;
```

### 複合インデックスの実験

```sql
-- 複合インデックス作成
CREATE INDEX idx_todos_user_done ON todos(user_id, done_flag);

-- 複合条件での検索
EXPLAIN ANALYZE
SELECT * FROM todos
WHERE user_id = 1 AND done_flag = false;

-- Index Only Scanの確認
EXPLAIN ANALYZE
SELECT user_id, done_flag
FROM todos
WHERE user_id = 1;
```

## よく使うコマンド

### 統計情報の更新

```sql
-- 全テーブルの統計情報を更新
ANALYZE;

-- 特定テーブルの統計情報を更新
ANALYZE todos;
```

### インデックスの管理

```sql
-- インデックス一覧
\di

-- インデックスのサイズ確認
SELECT
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) as size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;

-- 使われていないインデックスを発見
SELECT
    schemaname || '.' || tablename AS table,
    indexname AS index,
    idx_scan AS scans
FROM pg_stat_user_indexes
WHERE idx_scan = 0
    AND indexrelname NOT LIKE '%pkey'
ORDER BY pg_relation_size(indexrelid) DESC;

-- インデックス削除
DROP INDEX idx_todos_user_id;
```

### テーブル情報

```sql
-- テーブルサイズ
SELECT
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- レコード数
SELECT COUNT(*) FROM todos;

-- テーブル構造
\d todos
```

## トラブルシューティング

### データベースに接続できない

```bash
# コンテナが起動しているか確認
docker ps | grep go-demo-db

# コンテナを再起動
docker compose restart db
```

### テーブルが存在しない

```bash
# スキーマを再適用
make refresh-schema
```

### ポート5432が使えない

`.env`ファイルで`POSTGRES_PORT`を変更してください。
また、`demo/explain_demo/main.go`の`port`定数も同じ値に変更してください。

## 参考資料

- [PostgreSQL EXPLAIN完全ガイド](../../docs/postgresql_explain_guide.md)
- [PostgreSQL公式ドキュメント - EXPLAIN](https://www.postgresql.org/docs/current/sql-explain.html)

## 学習のヒント

1. まずはデモプログラムを実行して、全体像を把握する
2. PostgreSQLに直接接続して、自分でクエリを書いてみる
3. インデックスを追加/削除して、実行計画の変化を観察する
4. 大量データを投入して、パフォーマンスの違いを体感する
5. 実行計画の各要素（コスト、行数、実行時間）の意味を理解する

実際に手を動かすことが、EXPLAINを理解する一番の近道です！
