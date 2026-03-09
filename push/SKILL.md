---
name: push
description: lint・test・build・codex review の通過を確認してから git push する。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git rev-parse:*) Bash(git push:*) Bash(git tag:*) Bash(make lint:*) Bash(make test:*) Bash(make build:*) Bash(make -n:*) Bash(npm run:*) Bash(npm test:*) Bash(yarn run:*) Bash(yarn test:*) Bash(pnpm run:*) Bash(pnpm test:*) Bash(cargo build:*) Bash(cargo clippy:*) Bash(cargo test:*) Bash(go build:*) Bash(go vet:*) Bash(go test:*) Read
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

### 3. チェックタグを確認する

現在の HEAD に `check/*` タグが設置されているかを確認する。

```bash
git tag --points-at HEAD | grep '^check/'
```

### 4. タグに基づいてチェックを実行する

各チェック項目について、タグが現在の HEAD にあれば通過済みとしてスキップする。タグがなければ実行する。

| チェック | タグ             | タグあり       | タグなし                     |
| -------- | ---------------- | -------------- | ---------------------------- |
| build    | `check/build`    | スキップ       | 実行し、成功したらタグ設置   |
| lint     | `check/lint`     | スキップ       | 実行し、成功したらタグ設置   |
| test     | `check/test`     | スキップ       | 実行し、成功したらタグ設置   |
| review   | `check/review`   | スキップ       | ブロック（ユーザーに確認）   |

実行順序: build → lint → test。いずれかが失敗したら **push せずに停止** し、失敗内容をユーザーに報告する。

成功したチェックにはタグを設置する:

```bash
git tag -f "check/<type>" HEAD
```

review はこのスキル内では実行しない。タグがなければ「codex review が未実施です。先にレビューを実行しますか？」とユーザーに確認する。

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
- codex review のみ、タグ未設置の場合はブロッカーとして扱う
- build/lint/test が成功したら自動的に `check/<type>` タグを設置する
- `$ARGUMENTS` で `--skip-review` が指定された場合は codex review の確認をスキップする
- `$ARGUMENTS` で `--force` が指定された場合は `git push --force-with-lease` を使う（ユーザーに確認後）
