# OIDC アーキテクチャ解説：フロー・OAuth 2.0との違い・登場人物の関係

## 目次

1. [OIDCのフロー図（認可コードフロー）](#1-oidcのフロー図認可コードフロー)
2. [OAuth 2.0との違い](#2-oauth-20との違い)
3. [RP / OP / Resource Server / Session の関係](#3-rp--op--resource-server--session-の関係)

---

## 1. OIDCのフロー図（認可コードフロー）

### 全体フロー

最も一般的な**認可コードフロー（Authorization Code Flow）**の全体像です。

```mermaid
sequenceDiagram
    participant User as ユーザー（ブラウザ）
    participant RP as RP（リライングパーティ）
    participant OP as OP（OpenIDプロバイダ）
    participant RS as Resource Server（リソースサーバー）

    Note over RP: state, nonce, code_verifier を生成<br/>セッションに保存<br/>code_challenge を計算

    User->>RP: 1. ログインボタンをクリック
    RP->>User: 2. OPへリダイレクト<br/>response_type=code<br/>scope=openid profile email<br/>state=xxx&nonce=yyy<br/>code_challenge=zzz

    User->>OP: 3. 認証画面にアクセス
    OP->>User: 4. ID/PW入力画面を表示
    User->>OP: 5. 認証情報を送信
    OP->>OP: 6. 認証成功<br/>OP認証セッションを作成

    OP->>User: 7. RPのコールバックURLへリダイレクト<br/>code=AUTH_CODE&state=xxx
    User->>RP: 8. コールバック受信

    Note over RP: 【検証A】state の検証

    RP->>OP: 9. トークンリクエスト（サーバー間通信）<br/>code=AUTH_CODE<br/>code_verifier<br/>client_id + client_secret
    OP->>OP: 10. PKCE検証<br/>code_verifier → code_challenge一致？

    OP->>RP: 11. トークンレスポンス<br/>IDトークン + アクセストークン + リフレッシュトークン

    Note over RP: 【検証B】IDトークンの検証<br/>署名・iss・aud・exp・nonce

    RP->>RP: 12. RPセッションを作成
    RP->>User: 13. セッションCookie発行

    User->>RP: 14. リソース要求（セッションCookie）
    RP->>RS: 15. APIリクエスト（アクセストークン）
    RS->>RP: 16. リソース返却
    RP->>User: 17. レスポンス返却
```

### 認可コードフロー + PKCE（SPA / モバイルアプリ向け）

`client_secret`を安全に保持できないクライアント向けのフローです。

```mermaid
sequenceDiagram
    participant User as ユーザー（ブラウザ / アプリ）
    participant RP as RP（SPA / モバイルアプリ）
    participant OP as OP（OpenIDプロバイダ）

    Note over RP: code_verifier = ランダム値<br/>code_challenge = SHA256(code_verifier)<br/>state, nonce を生成

    User->>RP: 1. ログイン開始
    RP->>OP: 2. 認証リクエスト<br/>response_type=code<br/>code_challenge=xxx<br/>code_challenge_method=S256

    OP->>User: 3. 認証画面
    User->>OP: 4. 認証
    OP->>RP: 5. 認可コード返却

    RP->>OP: 6. トークンリクエスト<br/>code + code_verifier<br/>（client_secret は不要）

    Note over OP: SHA256(code_verifier)<br/>== code_challenge を検証

    OP->>RP: 7. IDトークン + アクセストークン
```

### Implicit フロー（非推奨）

```mermaid
sequenceDiagram
    participant User as ユーザー（ブラウザ）
    participant RP as RP
    participant OP as OP

    User->>RP: 1. ログイン開始
    RP->>OP: 2. 認証リクエスト<br/>response_type=id_token token

    OP->>User: 3. 認証
    OP->>User: 4. リダイレクト<br/>フラグメント（#）にトークンを含む<br/>#id_token=xxx&access_token=yyy

    Note over User: ⚠️ トークンがURLに露出<br/>ブラウザ履歴・Refererヘッダーから漏洩リスク

    User->>RP: 5. フラグメントからトークン取得（JS）
```

> Implicit フローはセキュリティ上の問題から**非推奨**です。PKCEを使った認可コードフローを使用してください。

---

## 2. OAuth 2.0との違い

### 根本的な違い：「認可」vs「認証 + 認可」

- **OAuth 2.0** = 認可（Authorization）：リソースへのアクセス許可。「何ができるか」
- **OIDC** = 認証（Authentication）+ 認可：ユーザーの身元確認 + アクセス許可。「誰であるか」+「何ができるか」

OIDCはOAuth 2.0を**包含する上位プロトコル**です。OAuth 2.0の仕組みをそのまま使いながら、認証のための標準仕様を追加しています。

### 比較表

| 観点 | OAuth 2.0 | OIDC |
|------|-----------|------|
| **目的** | 認可（リソースアクセスの許可） | 認証（身元確認）+ 認可 |
| **答える問い** | 「このアプリはリソースにアクセスしてよいか？」 | 「このユーザーは誰か？」 |
| **発行されるトークン** | アクセストークン（+ リフレッシュトークン） | アクセストークン + **IDトークン**（+ リフレッシュトークン） |
| **ユーザー情報** | 標準化されていない（各API独自） | **標準化**（IDトークン + UserInfoエンドポイント） |
| **scopeの例** | `read:photos`, `write:calendar` | `openid`, `profile`, `email` |
| **仕様の性質** | フレームワーク（実装の自由度が高い） | プロトコル（厳密な仕様） |
| **メタデータの発見** | なし | **ディスカバリ**（`.well-known/openid-configuration`） |
| **トークン形式** | 未規定（Opaque Token可） | IDトークンは**必ずJWT** |

### OIDCが追加したもの

- **IDトークン（JWT）** — ユーザー認証情報の標準化
- **UserInfoエンドポイント** — 追加のユーザー属性を取得
- **ディスカバリ** — `.well-known/openid-configuration` でOP情報を自動取得
- **scope: openid** — OIDCフローの起動トリガー
- **nonce パラメータ** — リプレイ攻撃対策
- **JWKSエンドポイント** — 署名検証用公開鍵の提供

### OAuth 2.0だけで「認証」すると何が危険か

OAuth 2.0は認可プロトコルであり、認証のための仕組みを持ちません。それを無理に認証に使うと**トークン置換攻撃**が成立します。

```mermaid
sequenceDiagram
    participant Attacker as 攻撃者
    participant EvilApp as 悪意あるアプリ
    participant LegitApp as 正規アプリ
    participant Provider as プロバイダ（Google等）

    Attacker->>EvilApp: 1. ログイン（OAuth認証）
    EvilApp->>Provider: 2. 認可リクエスト
    Provider->>EvilApp: 3. アクセストークン発行

    Note over EvilApp: 攻撃者のアクセストークンを窃取

    EvilApp->>LegitApp: 4. 攻撃者のアクセストークンを<br/>正規アプリに送りつける

    LegitApp->>Provider: 5. /userinfo でユーザー情報取得
    Provider->>LegitApp: 6. 攻撃者のユーザー情報を返却

    Note over LegitApp: ⚠️ アクセストークンの「発行先」を<br/>検証する仕組みがないため<br/>攻撃者を正規ユーザーとして受け入れてしまう
```

**OIDCのIDトークンが解決する理由：**

| 検証項目 | OAuth 2.0のアクセストークン | OIDCのIDトークン |
|---------|--------------------------|-----------------|
| 発行元（iss） | 検証手段なし | `iss` クレームで検証 |
| 発行先（aud） | **検証手段なし** | `aud` クレームで自分のclient_idか検証 |
| 改ざん検知 | 不可（Opaque Tokenの場合） | JWT署名で検証 |
| リプレイ対策 | なし | `nonce` クレームで検証 |

特に**`aud`（audience）の検証**がポイントです。IDトークンには「どのアプリ向けに発行されたか」が明記されているため、他のアプリ向けのトークンを拒否できます。

---

## 3. RP / OP / Resource Server / Session の関係

### 各コンポーネントの責務

```
ユーザー（ブラウザ）
  ├──[セッションCookie]──→ RP（リライングパーティ）= あなたのアプリ
  │                          ├── RPセッション
  │                          └── ビジネスロジック
  │
  └──[認証画面 / OPセッションCookie]──→ OP（OpenIDプロバイダ）= Google / Entra ID 等
                                          ├── OPセッション
                                          ├── 認証エンジン
                                          └── トークン発行

OP ──[IDトークン / アクセストークン]──→ RP ──[アクセストークン]──→ Resource Server（API）
                                                                      ├── 保護されたリソース
                                                                      └── トークン検証
```

### 各コンポーネントの詳細

#### OP（OpenID Provider）

| 項目 | 内容 |
|------|------|
| **何を持つか** | ユーザーの認証情報（ID/PW、MFA設定等）、クライアント登録情報（client_id, client_secret, redirect_uri）、署名鍵（秘密鍵/公開鍵）|
| **何を保証するか** | ユーザーの**身元（Identity）**。「このユーザーは確かに本人である」ことをIDトークンの署名で保証する |
| **セッション** | **OPセッション**を持つ。ユーザーがOPに対して認証済みであることを記憶する。これによりSSO（シングルサインオン）が実現する |

#### RP（Relying Party）

| 項目 | 内容 |
|------|------|
| **何を持つか** | client_id / client_secret、IDトークン（検証後は破棄可）、アクセストークン / リフレッシュトークン、ユーザー情報（IDトークンから抽出） |
| **何を保証するか** | ユーザーとの**アプリケーションセッション**の維持。IDトークンの検証結果に基づいて「このセッションは認証済みユーザーのものである」ことを保証する |
| **セッション** | **RPセッション**を持つ。サーバーサイドセッション + セッションCookieでユーザーを追跡する |

#### Resource Server（リソースサーバー）

| 項目 | 内容 |
|------|------|
| **何を持つか** | 保護されたリソース（APIデータ）、アクセストークンの検証手段（OPの公開鍵 or トークンイントロスペクションエンドポイント） |
| **何を保証するか** | **リソースへのアクセス制御**。アクセストークンのスコープに基づいて、許可された操作のみを実行する |
| **セッション** | **セッションを持たない**（ステートレス）。リクエストごとにアクセストークンを検証する |

### セッションの全体像

認証フロー全体には**3つの独立したセッション**が存在します。

**ユーザーのブラウザが持つもの：**
- RPセッションCookie（例: `session_id=abc123`） → 毎リクエストでRPサーバーに送信
- OPセッションCookie（例: `SSID=xyz789`） → OPアクセス時に送信

**RPサーバーのセッションストア：**
- user_id, user_name, email
- access_token, refresh_token
- ログイン時刻

**OPサーバーのセッションストア：**
- 認証済みユーザーID
- 認証時刻
- 認証方法（MFA等）

#### セッション比較表

| セッション | 管理者 | 存在場所 | ライフサイクル | 目的 |
|-----------|--------|---------|---------------|------|
| **OPセッション** | OP | OPサーバー + ブラウザ（Cookie） | OPでログアウトするまで | ユーザーが認証済みであることを記憶。SSO実現 |
| **RPセッション** | RP | RPサーバー + ブラウザ（Cookie） | RPアプリのセッション有効期限まで | アプリにログイン済みであることを追跡 |
| **（Resource Serverはセッションなし）** | - | - | - | リクエスト単位でアクセストークンを検証 |

### セッションとトークンのライフサイクル

```mermaid
sequenceDiagram
    participant User as ユーザー
    participant RP as RP
    participant OP as OP
    participant RS as Resource Server

    Note over User,OP: ===== フェーズ1: 認証 =====

    User->>RP: ログイン開始
    RP->>OP: 認証リクエスト

    Note over OP: OPセッションが存在しない場合<br/>→ ログイン画面を表示

    OP->>User: 認証画面
    User->>OP: 認証情報入力

    Note over OP: ✅ OPセッション作成<br/>（ブラウザにOPセッションCookie発行）

    OP->>RP: 認可コード
    RP->>OP: トークンリクエスト
    OP->>RP: IDトークン + アクセストークン

    Note over RP: IDトークン検証<br/>→ ユーザー情報を抽出
    Note over RP: ✅ RPセッション作成<br/>（ブラウザにRPセッションCookie発行）<br/>アクセストークンをセッションに保存

    Note over User,RS: ===== フェーズ2: リソースアクセス =====

    User->>RP: APIリクエスト（RPセッションCookie）
    RP->>RS: アクセストークン付きリクエスト
    Note over RS: アクセストークンを検証<br/>（セッション不要・ステートレス）
    RS->>RP: リソース返却
    RP->>User: レスポンス

    Note over User,RS: ===== フェーズ3: SSO（別のRPにログイン） =====

    User->>RP: 別のRPアプリにアクセス
    RP->>OP: 認証リクエスト

    Note over OP: OPセッションが既に存在する<br/>→ ログイン画面をスキップ

    OP->>RP: 認可コード（認証画面なし）
    RP->>OP: トークンリクエスト
    OP->>RP: IDトークン + アクセストークン

    Note over RP: ✅ 別のRPセッションを作成<br/>（ユーザーはID/PW入力なしでログイン完了）
```

### 各コンポーネントが保持するトークンと情報の整理

| コンポーネント | 保持するもの | 行うこと |
|--------------|------------|---------|
| **OP（OpenIDプロバイダ）** | ユーザー認証情報、署名用秘密鍵、クライアント登録情報 | IDトークン（JWT）・アクセストークン・リフレッシュトークンの発行 |
| **RP（リライングパーティ）** | client_id / client_secret、アクセストークン、リフレッシュトークン、ユーザー情報（セッション内） | IDトークンの検証、セッション管理、RSへのAPIコール |
| **Resource Server** | 保護リソース、OPの公開鍵（JWKS） | アクセストークン検証、スコープに基づくアクセス制御 |

**トークンの流れ：** OP →（IDトークン + アクセストークン）→ RP →（アクセストークン）→ Resource Server

### よくある誤解と注意点

#### 誤解1: 「IDトークンをResource Serverに送る」

```
❌ RP → Resource Server: Authorization: Bearer <IDトークン>
✅ RP → Resource Server: Authorization: Bearer <アクセストークン>
```

IDトークンは**RPがユーザーの身元を確認するため**のトークンです。Resource ServerへのAPIリクエストには**アクセストークン**を使用します。

#### 誤解2: 「IDトークンをそのままセッションとして使い続ける」

```
❌ 毎リクエストでIDトークンの有効期限を確認してセッション管理
✅ IDトークンは認証時に一度検証し、以降はRPセッションで管理
```

IDトークンは短命（数分〜数十分）に設定されるべきです。ログイン後のセッション維持はRPのサーバーサイドセッションが担います。

#### 誤解3: 「OPセッションが切れたらRPセッションも切れる」

OPセッションとRPセッションは**独立**しています。OPでログアウトしてもRPセッションは自動的には切れません（OIDCには[RP-Initiated Logout](https://openid.net/specs/openid-connect-rpinitiated-1_0.html)や[Back-Channel Logout](https://openid.net/specs/openid-connect-backchannel-1_0.html)といったログアウト仕様がありますが、実装は任意です）。

### 関係性のまとめ

```
OP（認証の権威）
├── ユーザーの身元を保証する
├── IDトークンを発行して「誰であるか」を伝える
├── アクセストークンを発行して「何ができるか」を伝える
├── OPセッションを管理してSSOを実現する
└── 署名鍵を管理してトークンの真正性を担保する

RP（アプリケーション）
├── IDトークンを検証して「誰がログインしたか」を確認する
├── RPセッションを作成してアプリ内でのログイン状態を維持する
├── アクセストークンを使ってResource Serverからデータを取得する
└── セッションCookieでユーザーのブラウザを追跡する

Resource Server（APIサーバー）
├── アクセストークンを検証してリクエストの正当性を確認する
├── スコープに基づいてアクセス制御を行う
├── セッションは持たない（ステートレス）
└── ユーザーの認証には関与しない（RPとOPの責務）
```
