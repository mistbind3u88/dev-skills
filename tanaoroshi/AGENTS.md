# tanaoroshi — Issue/PR 棚卸しツール

## 前提ツール

- [gh](https://cli.github.com/) — 認証済みであること
- [go](https://go.dev/) — `go run` で直接実行する

## 既知の制約

### gh コマンド失敗時の挙動

`gh issue list` / `gh pr list` が失敗した場合（認証エラー、レート制限、リポジトリへのアクセス不可など）、該当リポジトリの結果は空配列になりプログラムは正常終了する。stderr に警告は出るが、exit code には反映されない。棚卸し結果が不完全になる可能性がある点を踏まえ、警告が出ていないか必ず確認する。

### collect の取得上限

`collect` は `--limit 100` 固定で GitHub CLI を呼び出す。100 件を超える Open Issue/PR があるリポジトリではデータが切り詰められる。超過時は stderr に警告が出るため、該当リポジトリは別途個別に確認する。
