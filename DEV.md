# 開発ガイド

## 前提

- macOS / Linux
- Go は [mise](https://mise.jdx.dev/) で管理 (`mise.toml` でバージョン固定)
- 動作確認には `shpool` 本体が PATH にあると便利 (`mise` 経由で `cargo:shpool` を入れてもよい)

## 環境セットアップ

```sh
mise install   # mise.toml の Go バージョンを入れる
```

## ビルド/実行せずに試す

```sh
go run ./cmd/shp --print-name   # cwd 由来のセッション名を表示 (shpool 不要)
go run ./cmd/shp --help
go run ./cmd/shp                # TUI: cwd 名 (デフォルト選択) + 既存セッション
go run ./cmd/shp my-name        # TUIを介さず指定名で attach (--dir .)
go run ./cmd/shp -f             # cwd 名で force-attach (--dir .)
```

## テスト / ビルド

主な操作は `mise run` のタスクから実行できる。`mise tasks` でタスク一覧。

```sh
mise run build      # go build -o ./bin/shp ./cmd/shp
mise run install    # ~/.local/bin/shp にインストール
mise run check      # vet + test
mise run test       # go test ./...
mise run vet        # go vet ./...
```

素の go コマンドでも当然OK。

```sh
go test ./...
go vet ./...
go build ./cmd/shp
```

## ソースから入れる

リリースバイナリではなくソースから直接入れたい場合。

```sh
# go install で $GOBIN (デフォルト $HOME/go/bin) に入れる
go install github.com/uzulla/shpool-launch/cmd/shp@latest

# あるいはこのリポジトリをクローンして任意の場所へ
git clone https://github.com/uzulla/shpool-launch.git
cd shpool-launch
go build -o ~/.local/bin/shp ./cmd/shp
```

クロスビルド例:

```sh
GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o dist/shp-linux-amd64 ./cmd/shp
GOOS=linux  GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o dist/shp-linux-arm64 ./cmd/shp
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o dist/shp-darwin-amd64 ./cmd/shp
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o dist/shp-darwin-arm64 ./cmd/shp
```

## パッケージ構成

```text
cmd/shp/main.go              CLI エントリポイント (flag パースと分岐)
internal/session/name.go     cwd → セッション名の変換
internal/session/name_test.go
internal/shpool/attach.go    syscall.Exec で `shpool attach` に置換
internal/shpool/list.go      `shpool list` 実行と出力パース
internal/shpool/list_test.go
internal/tui/select.go       Bubble Tea の peco 風絞り込みTUI
internal/tui/select_test.go
```

依存ライブラリは Bubble Tea / Lipgloss のみ。CLI引数パースは標準 `flag`。

## セッション名生成ルール

1. `os.Getwd()` で現在ディレクトリを取得
2. 絶対パス化 (`filepath.Abs`)
3. `$HOME` 配下なら `$HOME/` を取り除く
4. 先頭・末尾の `/` を取り除く
5. `/` を `.` に置換
6. 英数字、`.`、`-`、`_` を維持する
7. それ以外の文字 (空白を含む) はすべて `_` に置換
8. 置換が発生した場合、またはディレクトリ名自体に `.` が含まれる場合は、衝突緩和のため相対パス由来の短い hash suffix を付ける

例:

| 入力 | 出力 |
| --- | --- |
| `/home/uzulla/work/company-a/api` | `work.company-a.api` |
| `/Users/uzulla/src/foo bar/api`   | `src.foo_bar.api-a8163889` |
| `/tmp/test/api`                    | `tmp.test.api` |
| `/Users/uzulla/プロジェクト/api`  | `______.api-f10f1a6e` (非ASCIIは1文字あたり `_` + suffix) |
| `/srv/app-v1.2.3/api`              | `srv.app-v1.2.3.api-597dbea5` |

MVP ではこのルールは固定。将来は `~/work` や `~/src` を root として相対化するなどの設定を追加できる設計にしてある (`session.FromPath(path, home)` で home を引数化済み)。

## attach の実行方式

`shpool attach` 実行時は、子プロセスを起動するのではなく `syscall.Exec` で `shp` プロセスを `shpool` に置き換える (`internal/shpool/attach.go`)。
これによりシグナル伝播や端末制御が `shpool` に直接渡る。

attach には常に `--dir .` を渡す。既存セッションへの attach では実質無視され、新規作成時は現在のディレクトリから shell が始まる。

## `shpool list` パース方針

一覧取得は `shpool list` で実行する (`internal/shpool/list.go`)。
一覧取得ができない場合は、既存セッション候補なしとして cwd 由来の候補だけを表示する。

- 空行を除外
- `NAME STATUS` のヘッダ行を除外
- 各行の先頭カラムをセッション名として扱う
- 候補が0件なら標準エラーに `No shpool sessions found.` を表示して終了

## TUI

`internal/tui/select.go`。Bubble Tea v1 + Lipgloss。

- alt screen で起動
- 上部に `QUERY> ...` を表示し、下にフィルタ済み候補
- フィルタ: case-insensitive substring AND (`Filter(items, query)` として独立、テスト対象)
- キー操作: 文字入力で絞り込み、Backspace、↑ / Ctrl-P、↓ / Ctrl-N、Enter、Esc / Ctrl-C
- 絞り込み結果が0件の Enter は新規セッション名として query を返す。query が `-` で始まり cwd 既定項目がある場合は `defaultItem + query` を返す。
- Backspace は rune 単位で削除する
- 候補が端末高に収まらない場合はカーソル位置に合わせてスクロールする
- `SelectWithDefault(items, defaultItem)` で「先頭に置く既定項目」を渡せる。既定項目は `(cwd)` ラベル付きで表示され、初期カーソル位置になる。既存リスト内に重複があれば dedupe。

`shp` (引数なし) はこの `SelectWithDefault` を使って `cwd 由来のセッション名` + `shpool list の既存セッション` を一つのリストで提示する。

MVP では fuzzy / negative / 正規表現 / 複数選択 / preview pane / sort 切替は持たない。

## エラー処理

| 状況 | 出力 | 終了コード |
| --- | --- | --- |
| `shpool` が PATH に無い | `Error: shpool command not found in PATH. Please install shpool first` | 1 |
| `shpool list` で一覧取得不可 | cwd 候補だけを表示。cwd 名も空なら `No shpool sessions found.` | 0 |
| `shpool list` 失敗 (一覧取得不可として扱わないもの) | `Error: shpool list failed: <stderr>` | 1 |
| セッション0件 | `No shpool sessions found.` (stderr) | 0 |
| TUI でキャンセル | (何も出さない) | 0 |
| 不正な引数 | `Error: too many arguments` 等 + usage | 1 |
