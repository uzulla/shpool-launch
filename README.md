# shp - shpool session selector

`shpool attach` を便利にするツール。

- 引数なしで `shp` を実行すると、peco 風の TUI が開く。カレントディレクトリ由来のセッション名が先頭・初期選択になっていて、Enter でそのまま attach できる(無ければ現在のディレクトリで新規作成)。
- 既存セッションも同じリストに並ぶので、↓ や絞り込みで切り替えて attach できる。
- 既存セッション一覧が取れない場合でも、cwd 由来の候補だけで起動できる。

`shpool` (https://github.com/shell-pool/shpool) 本体が PATH に存在することが前提です。

## インストール

[Releases](https://github.com/uzulla/shpool-launch/releases) から OS / アーキテクチャに合うものをダウンロードして、PATH の通った場所に置く。

あるいは mise (go backend) で入れる。

```sh
mise use -g go@latest # Goがはいっているなら不要
mise use -g "go:github.com/uzulla/shpool-launch/cmd/shp@latest"
```

ローカルからソースビルドする手順は [DEV.md](./DEV.md) を参照。

## 使い方

### `shp`

TUI ピッカーを開く。先頭は cwd 由来のセッション名 (`/`を`.`に変換したものがデフォルト選択)、その後ろに `shpool list` の既存セッションが並ぶ。Enter で attach。
既存セッション一覧が取れない場合は、cwd 候補だけを表示する。

```sh
$ pwd
/Users/uzulla/work/company-a/api

$ shp

QUERY>

> work.company-a.api (cwd)
  work.company-b.api
  sandbox.api-test
```

- Enter でカーソル位置に attach (該当セッションが無ければ `shpool attach --dir .` で現在のディレクトリに新規作成)
- 絞り込み結果が0件の状態で Enter すると、入力した文字列で新規セッションを作成する。空白を含む入力は作成しない。cwd 候補がある状態で `-another` のように `-` から入力した場合は、cwd 候補にサフィックスとして足して `work.company-a.api-another` のように作成する。
- 文字を入力すれば case-insensitive substring AND で絞り込める
- 絞り込みは case-insensitive の substring AND マッチ。`company api` のように空白で AND 検索できる。
- 候補が画面に収まらない場合はカーソル移動に合わせてスクロールする。

### `shp <session-name>`

TUI を介さず、指定された名前で直接 attach する。

```sh
shp my-session
# → shpool attach --dir . my-session
```

### `shp -f` / `shp -f <session-name>`

force attach する(既存セッションを奪い取る)。こちらは TUI を介さず直接 attach。

```sh
shp -f             # cwd 由来の名前で force-attach
shp -f my-session  # 指定名で force-attach
```

いずれも新規作成時は `--dir .` を渡すため、作成される shell は現在のディレクトリから始まる。

## セッション名

基本は `$HOME` からの相対パスを使い、`/` を `.` に置き換える。
英数字、`.`、`-`、`_` 以外は `_` に置き換えるが、置換やディレクトリ名中の `.` による衝突を避けるため、必要な場合だけ短い hash suffix を付ける。

例:

| cwd | セッション名 |
| --- | --- |
| `/Users/uzulla/work/company-a/api` | `work.company-a.api` |
| `/Users/uzulla/src/foo bar/api` | `src.foo_bar.api-a8163889` |
| `/Users/uzulla/プロジェクト/api` | `______.api-f10f1a6e` |

## 開発

開発環境のセットアップやテスト方法、内部仕様は [DEV.md](./DEV.md) を参照。

## ライセンス

MIT
