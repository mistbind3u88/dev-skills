---
name: push
description: lint・test・build・codex review の通過を確認してから git push する。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git rev-parse:*) Bash(git push:*) Bash(make lint:*) Bash(make test:*) Bash(make build:*) Bash(make -n:*) Bash(npm run:*) Bash(npm test:*) Bash(yarn run:*) Bash(yarn test:*) Bash(pnpm run:*) Bash(pnpm test:*) Bash(cargo build:*) Bash(cargo clippy:*) Bash(cargo test:*) Bash(go build:*) Bash(go vet:*) Bash(go test:*) Bash(codex review:*) Read
---

# push スキル

lint・test・build が通り、codex review 済みであることを確認してから push する。

## 手順

### 1. 前提確認

```bash
git status -s
git log --oneline main..HEAD
git rev-parse --abbrev-ref HEAD
```

- 未コミットの変更がある場合は push せず、先にコミットするようユーザーに伝える
- main ブランチにいる場合は警告してユーザーに確認する

### 2. ビルドツールを検出する

プロジェクトルートのファイルから利用可能なコマンドを特定する。

| ファイル       | 検出方法                                                   |
| -------------- | ---------------------------------------------------------- |
| `Makefile`     | `make` のターゲット一覧から `lint`、`test`、`build` を探す |
| `package.json` | `scripts` フィールドから `lint`、`test`、`build` を探す    |
| `Cargo.toml`   | `cargo clippy`、`cargo test`、`cargo build` を使う         |
| `go.mod`       | `go vet`、`go test`、`go build` を使う                     |

複数のビルドツールが存在する場合は、`Makefile` > `package.json` > 言語固有ツールの優先順位で選択する。
該当するターゲットやスクリプトが存在しない項目はスキップする。

### 3. build・lint・test を実行する

検出したコマンドを以下の順で実行する。いずれかが失敗したら **push せずに停止** し、失敗内容をユーザーに報告する。

1. **build** (例: `make build`、`npm run build`、`cargo build`)
2. **lint** (例: `make lint`、`npm run lint`、`cargo clippy`)
3. **test** (例: `make test`、`npm test`、`cargo test`)

### 4. codex review の実行済みを確認する

`~/.codex/sessions/` 配下のセッションファイルから、現在のブランチ・リポジトリに対する review セッションが直近に実行されたかを確認する。

```bash
# 今日の日付ディレクトリからレビューセッションを探す
find ~/.codex/sessions/$(date +%Y/%m/%d) -name "rollout-*.jsonl" -newer $(git rev-parse --git-dir)/HEAD 2>/dev/null \
  | while read f; do head -1 "$f"; done \
  | grep -l '"originator":"codex_exec"'
```

セッションファイルのメタデータ (`session_meta`) で以下を確認する:

- `originator` が `codex_exec` (review 実行である)
- `git.branch` が現在のブランチと一致する
- `git.commit_hash` が現在の HEAD、または現在の HEAD の祖先である

レビュー済みが確認できない場合は、ユーザーに「codex review が未実施です。先にレビューを実行しますか？」と確認する。

### 5. 結果サマリーを表示する

全チェック項目の結果を一覧で表示する。

```
push 前チェック:
  lint:         OK
  test:         OK
  build:        OK (スキップ: ターゲットなし)
  codex review: OK (2026-03-06T10:30 に実施)
```

全て OK の場合のみ次のステップに進む。

### 6. push する

```bash
git push -u origin HEAD
```

push 後、結果を報告する。

## 注意

- 検出できなかった項目 (lint/test/build) は「スキップ」として扱い、ブロッカーにしない
- codex review のみ、未実施の場合はブロッカーとして扱う
- `$ARGUMENTS` で `--skip-review` が指定された場合は codex review の確認をスキップする
- `$ARGUMENTS` で `--force` が指定された場合は `git push --force-with-lease` を使う（ユーザーに確認後）
