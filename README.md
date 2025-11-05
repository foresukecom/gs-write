# gs-write

標準入力を新しいGoogleスプレッドシートに書き出すシンプルなCLIツールです。
UNIX哲学に基づき、他のコマンドとパイプ(`|`)で連携することを前提に設計されています。

![Go](https://img.shields.io/badge/Go-1.24-blue.svg)

## 概要

`gs-write` は、`ls -l` や `cat report.csv`、`ps aux` といったコマンドの実行結果を、直接新しいGoogleスプレッドシートに出力するためのコマンドです。
実行のたびに常に新しいスプレッドシートが作成され、完了後にはそのシートのURLが標準出力に返されます。これにより、シェルスクリプト内でのデータ連携や、結果のURLをクリップボードにコピーするなどの操作が容易になります。

## 特徴

- 標準入力からCSVデータを受け取り、Googleスプレッドシートに書き込み
- 常に新しいスプレッドシートを作成（既存シートへの追記機能はありません）
- パイプ(`|`)で他のコマンドとスムーズに連携可能
- `--title`オプションでスプレッドシートのタイトルを自由に指定可能
  - タイトルが指定されない場合は、実行日時から自動で命名 (`YYYYMMDDHHMMSS+gs`)
- 成功時に、作成されたスプレッドシートのURLを標準出力に返す

## インストール

Go言語の環境がセットアップされている場合、以下のコマンドでインストールできます。

```bash
go install github.com/your-username/gs-write@latest
```

または、このリポジトリをクローンしてビルド：

```bash
git clone https://github.com/your-username/gs-write.git
cd gs-write
go build -o gs-write .
```

## セットアップ

### 1. Google Cloud Consoleでプロジェクトを作成

1. [Google Cloud Console](https://console.cloud.google.com/)にアクセス
2. 新しいプロジェクトを作成
3. Google Sheets APIを有効化
4. OAuth 2.0クライアントIDを作成（アプリケーションの種類：デスクトップアプリ）
5. `credentials.json`をダウンロード

### 2. 認証

初回使用時、または認証が必要な場合は以下のコマンドを実行します：

#### 方法1: インタラクティブに認証情報を入力

```bash
gs-write auth
```

プロンプトに従って、`credentials.json`の内容を貼り付けてください。

#### 方法2: ファイルから認証情報を読み込み

```bash
gs-write auth --credentials ./credentials.json
```

認証が成功すると、認証情報とトークンが `~/.config/gs-write/config.toml` に保存されます。

## 使い方

### 基本的な使い方

```bash
# CSVファイルをスプレッドシートに変換
cat data.csv | gs-write

# コマンドの出力をスプレッドシートに保存
ls -l | gs-write

# タイトルを指定してスプレッドシートを作成
ps aux | gs-write --title "Process List $(date +%Y%m%d)"

# URLをクリップボードにコピー（macOS）
cat report.csv | gs-write | pbcopy

# URLをクリップボードにコピー（Linux with xclip）
cat report.csv | gs-write | xclip -selection clipboard
```

### オプション

- `--title <タイトル>`: スプレッドシートのタイトルを指定します。指定しない場合は、タイムスタンプから自動生成されます。

### サブコマンド

#### `gs-write auth`

Google Sheets APIとの認証を行います。

```bash
# インタラクティブに認証
gs-write auth

# ファイルから認証情報を読み込み
gs-write auth --credentials ./credentials.json
```

#### `gs-write version`

バージョン情報を表示します。

```bash
gs-write version
```

## データフォーマット

`gs-write`はCSV形式のデータを想定しています。標準入力から読み込んだデータはカンマ区切り(`,`)として解析されます。

### 例

```bash
echo "Name,Age,City
Alice,30,Tokyo
Bob,25,Osaka" | gs-write --title "User List"
```

## プロジェクト構成

```
.
├── README.md           # このファイル
├── cmd/                # Cobraコマンド定義
│   ├── auth.go         # 認証コマンド
│   ├── root.go         # ルートコマンド（メイン機能）
│   └── version.go      # バージョンコマンド
├── pkg/                # 内部パッケージ
│   ├── auth/           # 認証処理
│   │   └── auth.go
│   └── sheets/         # Google Sheets API クライアント
│       └── sheets.go
├── go.mod              # Go Modules
├── go.sum              # Go Modules チェックサム
└── main.go             # エントリーポイント
```

## トラブルシューティング

### 認証エラーが発生する

認証トークンの有効期限が切れている可能性があります。再度認証を実行してください：

```bash
gs-write auth
```

### APIクォータエラー

Google Sheets APIには利用制限があります。大量のリクエストを送信している場合は、少し時間をおいてから再試行してください。

## ライセンス

MIT License

## 貢献

Issue や Pull Request を歓迎します！

## 開発

### 開発環境

このプロジェクトはVS Code Dev Containersを使用して開発できます。

1. VS Codeでプロジェクトを開く
2. "Reopen in Container"を選択
3. コンテナ内で開発を開始

### ビルド

```bash
go build -o gs-write .
```

### テスト

```bash
go test ./...
```
