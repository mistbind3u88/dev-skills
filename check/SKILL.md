---
name: check
description: mark で未チェックの項目（lint・build・test・review）を検出し、一通り実行する。
allowed-tools: Bash(make lint:*) Bash(make test:*) Bash(make build:*) Bash(make -n:*) Bash(npm run:*) Bash(npm test:*) Bash(yarn run:*) Bash(yarn test:*) Bash(pnpm run:*) Bash(pnpm test:*) Bash(cargo build:*) Bash(cargo clippy:*) Bash(cargo test:*) Bash(go build:*) Bash(go vet:*) Bash(go test:*) Read
---

# check スキル

mark タグが未設置の項目を検出し、実行する。

## 手順

### 1. チェックタグを確認する

`/mark --status` を実行して、現在の HEAD のチェック通過状況を確認する。

全項目が現在の HEAD にタグ設置済みなら、その旨を報告して終了する。

### 2. ビルドツールを検出する

以下の優先順位でコマンドを特定する。ドキュメントに記載があればそれを使い、自動検出には進まない。

#### 優先: リポジトリのドキュメント

リポジトリの `AGENTS.md`、`CLAUDE.md`、`README.md`、CI 定義、またはスキル定義（`.claude/skills/` 配下）に lint・test・build の実行コマンドが記載されていればそれを使う。ドキュメントで特定できた時点で検出を終了し、フォールバックには進まない。

#### フォールバック: プロジェクトルートのファイルから自動検出

ドキュメントに該当コマンドの記載がない項目についてのみ、プロジェクトルートのファイルから推定する。

| ファイル       | 検出方法                                                   |
| -------------- | ---------------------------------------------------------- |
| `Makefile`     | `make` のターゲット一覧から `lint`、`test`、`build` を探す |
| `package.json` | `scripts` フィールドから `lint`、`test`、`build` を探す    |
| `Cargo.toml`   | `cargo clippy`、`cargo test`、`cargo build` を使う         |
| `go.mod`       | `go vet`、`go test`、`go build` を使う                     |

該当するターゲットやスクリプトが存在しない項目はスキップする。

### 3. 未チェック項目を実行する

タグが現在の HEAD にない項目を順に実行する。

| チェック | タグなし                         | タグあり（現在の HEAD） |
| -------- | -------------------------------- | ----------------------- |
| lint     | 実行し、成功したら `/mark lint`  | スキップ                |
| build    | 実行し、成功したら `/mark build` | スキップ                |
| test     | 実行し、成功したら `/mark test`  | スキップ                |
| review   | `/codex-review` を実行           | スキップ                |

実行順序: lint → build → test → review。いずれかが失敗したら停止し、失敗内容をユーザーに報告する。

### 4. 結果サマリーを表示する

全チェック項目の結果を一覧で表示する。

```
チェック結果:
  lint:         OK
  build:        OK (スキップ: ターゲットなし)
  test:         OK
  codex review: OK
```

## 注意

- 検出できなかった項目（lint/build/test）は「スキップ」として扱い、ブロッカーにしない
- build/lint/test が成功したら `/mark <type>` でタグを設置する
- review は `/codex-review` が完了時に自動でタグを設置する
- `$ARGUMENTS` で `--skip-review` が指定された場合は review をスキップする
