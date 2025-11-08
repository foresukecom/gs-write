# gs-write

[English](README_EN.md) | 日本語

標準入力を新しいGoogleスプレッドシートに書き出すシンプルなCLIツールです。
パイプ(`|`)で他のコマンドと簡単に連携できます。

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
- `--freeze-rows`と`--freeze-cols`オプションで行と列の固定表示が可能
- `--filter-header-row`オプションで基本フィルタの設定が可能
- 成功時に、作成されたスプレッドシートのURLを標準出力に返す

## インストール

### バイナリをダウンロード（推奨）

[Releases ページ](https://github.com/foresukecom/gs-write/releases)から、お使いの環境に応じたファイルをダウンロードしてください。

#### macOS

**Intel Mac の場合:**
```bash
# ダウンロード
curl -LO https://github.com/foresukecom/gs-write/releases/latest/download/gs-write_Darwin_x86_64.tar.gz
# 解凍
tar xzf gs-write_Darwin_x86_64.tar.gz
# 実行可能な場所に移動
sudo mv gs-write /usr/local/bin/
```

**Apple Silicon (M1/M2/M3) の場合:**
```bash
# ダウンロード
curl -LO https://github.com/foresukecom/gs-write/releases/latest/download/gs-write_Darwin_arm64.tar.gz
# 解凍
tar xzf gs-write_Darwin_arm64.tar.gz
# 実行可能な場所に移動
sudo mv gs-write /usr/local/bin/
```

#### Linux

```bash
# ダウンロード
curl -LO https://github.com/foresukecom/gs-write/releases/latest/download/gs-write_Linux_x86_64.tar.gz
# 解凍
tar xzf gs-write_Linux_x86_64.tar.gz
# 実行可能な場所に移動
sudo mv gs-write /usr/local/bin/
```

#### Windows

1. [gs-write_Windows_x86_64.zip](https://github.com/foresukecom/gs-write/releases/latest/download/gs-write_Windows_x86_64.zip) をダウンロード
2. ZIPファイルを解凍
3. `gs-write.exe` を任意のフォルダに配置
4. 必要に応じてPATHを通す

### Go言語環境からインストール

Go言語の環境がセットアップされている場合、以下のコマンドでインストールできます。

```bash
go install github.com/foresukecom/gs-write@latest
```

### ソースからビルド

このリポジトリをクローンしてビルド：

```bash
git clone https://github.com/foresukecom/gs-write.git
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

#### 認証フロー

コマンドを実行すると、以下の手順で認証を行います：

1. **認証URLが表示されます**
   ```
   Please visit the following URL to authorize this application:
   https://accounts.google.com/o/oauth2/auth?...
   ```

2. **ブラウザでURLにアクセス**
   - 表示されたURLをブラウザで開きます
   - Googleアカウントでログインします

3. **アプリケーションを承認**
   - 「[アプリ名] が Google アカウントへのアクセスを求めています」という画面が表示されます
   - 「続行」をクリックします

4. **認証コードを取得**
   - 承認後、`http://localhost` へリダイレクトされ、通常は `ERR_CONNECTION_REFUSED` エラーが表示されます
   - **これは正常な動作です**
   - ブラウザのアドレスバーに表示されているURLから `code=` 以降の文字列をコピーします
   - 例: `http://localhost/?code=4/0AbCD...XYZ&scope=...` の場合、`4/0AbCD...XYZ` の部分がコードです

5. **認証コードを入力**
   ```
   Enter the authorization code:
   ```
   - ターミナルに戻り、コピーした認証コードを貼り付けてEnterキーを押します

6. **認証完了**
   ```
   Authentication successful!
   Authentication saved to: ~/.config/gs-write/auth.json
   ```

認証が成功すると、認証情報とトークンが `~/.config/gs-write/auth.json` に保存されます。

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

### 固定表示（Freeze Panes）

ヘッダー行や特定の列を固定表示することができます：

```bash
# 1行目（ヘッダー行）を固定
cat data.csv | gs-write --freeze-rows 1

# 最初の2列を固定
cat data.csv | gs-write --freeze-cols 2

# 1行目と1列目の両方を固定
cat data.csv | gs-write --freeze-rows 1 --freeze-cols 1

# タイトルと固定表示を組み合わせ
cat employee.csv | gs-write --title "社員リスト" --freeze-rows 1
```

### 基本フィルタ（Basic Filter）

データに対して基本フィルタを設定することができます：

```bash
# 1行目をヘッダーとしてフィルタを設定
cat data.csv | gs-write --filter-header-row 1

# 2行目をヘッダーとしてフィルタを設定
cat data.csv | gs-write --filter-header-row 2

# 固定表示とフィルタを組み合わせ
cat data.csv | gs-write --freeze-rows 1 --filter-header-row 1

# すべてのオプションを組み合わせ
cat employee.csv | gs-write --title "社員リスト" --freeze-rows 1 --filter-header-row 1
```

### 文字エンコーディング

Shift_JIS（SJIS）やEUC-JPなど、UTF-8以外のエンコーディングのCSVファイルを扱うことができます：

```bash
# Shift_JIS（SJIS）エンコーディングのCSVファイルを読み込む
cat data_sjis.csv | gs-write --encoding sjis

# EUC-JPエンコーディングのCSVファイルを読み込む
cat data_eucjp.csv | gs-write --encoding euc-jp

# UTF-8（デフォルト）の場合は指定不要
cat data_utf8.csv | gs-write

# エンコーディングとその他のオプションを組み合わせ
cat data_sjis.csv | gs-write --encoding sjis --title "社員リスト" --freeze-rows 1
```

サポートされているエンコーディング：
- `utf-8`（デフォルト）：UTF-8エンコーディング
- `sjis`：Shift_JIS（Windows標準の日本語エンコーディング）
- `euc-jp`：EUC-JP（Unix系の日本語エンコーディング）

### オプション

- `--title <タイトル>`: スプレッドシートのタイトルを指定します。指定しない場合は、タイムスタンプから自動生成されます。
- `--freeze-rows <行数>`: 上から指定した行数を固定表示します。設定ファイルの値を上書きします。
- `--freeze-cols <列数>`: 左から指定した列数を固定表示します。設定ファイルの値を上書きします。
- `--filter-header-row <行番号>`: 指定した行をヘッダーとして基本フィルタを設定します。設定ファイルの値を上書きします。
- `--encoding <エンコーディング>`: 入力CSVの文字エンコーディングを指定します（`utf-8`, `sjis`, `euc-jp`）。デフォルトは`utf-8`です。

### 設定ファイル

頻繁に使用する設定値（固定行数・固定列数・フィルタヘッダー行など）は、設定ファイルに保存しておくことができます。

設定ファイルは `~/.config/gs-write/config.toml` に保存され、`gs-write config` コマンドで管理できます。

コマンドラインオプションが指定された場合、設定ファイルの値は上書きされます（優先順位: CLI > 設定ファイル > デフォルト値）。

```bash
# 現在の設定を表示
gs-write config list

# 特定の設定値を取得
gs-write config get freeze.rows
gs-write config get filter.header_row

# 設定値を変更
gs-write config set freeze.rows 1
gs-write config set freeze.cols 2
gs-write config set filter.header_row 1

# 設定値を削除（デフォルト値に戻す）
gs-write config unset freeze.rows
gs-write config unset filter.header_row
```

#### 利用可能な設定項目

- `freeze.rows`: 固定する行数（デフォルト: 0）
- `freeze.cols`: 固定する列数（デフォルト: 0）
- `filter.header_row`: フィルタのヘッダー行番号（デフォルト: 0 = フィルタなし）

### サブコマンド

#### `gs-write auth`

Google Sheets APIとの認証を行います。

```bash
# インタラクティブに認証
gs-write auth

# ファイルから認証情報を読み込み
gs-write auth --credentials ./credentials.json
```

#### `gs-write config`

設定ファイルを管理します。

```bash
# 設定一覧を表示
gs-write config list

# 特定の設定値を取得
gs-write config get freeze.rows

# 設定値を変更
gs-write config set freeze.rows 1

# 設定値を削除
gs-write config unset freeze.rows
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
├── README.md           # このファイル（日本語）
├── README_EN.md        # このファイル（英語）
├── cmd/                # Cobraコマンド定義
│   ├── auth.go         # 認証コマンド
│   ├── config.go       # 設定コマンド
│   ├── root.go         # ルートコマンド（メイン機能）
│   └── version.go      # バージョンコマンド
├── pkg/                # 内部パッケージ
│   ├── auth/           # 認証処理
│   │   └── auth.go
│   ├── config/         # 設定管理
│   │   └── config.go
│   └── sheets/         # Google Sheets API クライアント
│       └── sheets.go
├── go.mod              # Go Modules
├── go.sum              # Go Modules チェックサム
└── main.go             # エントリーポイント
```

## 設定ファイルの場所

gs-writeは以下のファイルを使用します：

- `~/.config/gs-write/auth.json` - OAuth 2.0認証情報とトークン
- `~/.config/gs-write/config.toml` - ユーザー設定（freeze.rows、freeze.colsなど）

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
