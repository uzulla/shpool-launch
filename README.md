# shp - shpool tui launcher

`shpool attach` を便利にするツール。

- 引数なしで `shp` を実行すると、peco 風の TUI が開く。カレントディレクトリ由来のセッション名が先頭・初期選択になっていて、Enter でそのまま attach できる(無ければ新規作成)。
- 既存セッションも同じリストに並ぶので、↓ や絞り込みで切り替えて attach できる。

`shpool` 本体が PATH に存在することが前提です。

## インストール

### mise (go backend) で入れる

[mise](https://mise.jdx.dev/) を使うのが手軽です。`cargo:shpool` などと同じ感覚で go の wrapper backend が使えます。

グローバルに入れる場合:

```sh
mise use -g "go:github.com/uzulla/shpool-launch/cmd/shp@latest"
```

プロジェクトの `mise.toml` に固定したい場合は:

```toml
[tools]
"go:github.com/uzulla/shpool-launch/cmd/shp" = "latest"
```

### `go install`

```sh
go install github.com/uzulla/shpool-launch/cmd/shp@latest
```

### ローカルビルド

```sh
go build -o ~/.local/bin/shp ./cmd/shp
```

## 使い方

### `shp`

TUI ピッカーを開く。先頭は cwd 由来のセッション名 (デフォルト選択)、その後ろに `shpool list` の既存セッションが並ぶ。Enter で attach。

```sh
$ pwd
/Users/uzulla/work/company-a/api

$ shp
```

```text
QUERY>

> work.company-a.api (cwd)
  work.company-b.api
  sandbox.api-test
```

- Enter でカーソル位置に attach (該当セッションが無ければ shpool が新規作成)
- 文字を入力すれば case-insensitive substring AND で絞り込める

### `shp <session-name>`

TUI を介さず、指定された名前で直接 attach する。

```sh
shp my-session
# → shpool attach my-session
```

### `shp -f` / `shp -f <session-name>`

force attach する(既存セッションを奪い取る)。こちらは TUI を介さず直接 attach。

```sh
shp -f             # cwd 由来の名前で force-attach
shp -f my-session  # 指定名で force-attach
```

操作:

| キー | 動作 |
| --- | --- |
| 文字入力 | 絞り込み |
| Backspace | 1文字削除 |
| ↑ / Ctrl-P | 上へ移動 |
| ↓ / Ctrl-N | 下へ移動 |
| Enter | 選択して attach |
| Esc / Ctrl-C | キャンセル |

絞り込みは case-insensitive の substring AND マッチ。`company api` のように空白で AND 検索できる。

## セッション名の付き方

カレントディレクトリのパスから決まる。たとえば:

| 入力 | 出力 |
| --- | --- |
| `/home/uzulla/work/company-a/api` | `work.company-a.api` |
| `/Users/uzulla/src/foo bar/api`   | `src.foo_bar.api` |
| `/tmp/test/api`                    | `tmp.test.api` |

詳しい変換ルールは [DEV.md](./DEV.md) を参照。

## 開発

開発環境のセットアップやテスト方法、内部仕様は [DEV.md](./DEV.md) を参照。

## ライセンス

MIT
