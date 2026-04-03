---
name: start-dev
description: 新しい作業ブランチを作成し、Issue/PR の情報取得やタスクドキュメント確認を含む作業開始の準備を行う。
allowed-tools: Bash(git switch:*) Bash(git branch:*) Bash(git status:*) Bash(git rev-parse:*) Bash(git fetch:*) Bash(git pull:*) Read
---

# start-dev スキル

新しい作業を開始する準備を行う。ブランチ作成、Issue/PR の情報取得、タスクドキュメント確認をまとめて実施する。

## 引数

```
$ARGUMENTS: [Issue/PR 番号] [ブランチ名]
```

- Issue/PR 番号が指定された場合、`/gh-read` で情報を取得する
- ブランチ名が直接指定された場合はそのまま使用する

## 手順

### 1. 前提確認

```bash
git status -s
git rev-parse --abbrev-ref HEAD
```

- 未コミットの変更がある場合はコミットを促して停止する
- 現在のブランチを確認する

### 2. Issue/PR の情報取得

`$ARGUMENTS` に Issue/PR 番号が含まれる場合、`/gh-read` に委譲して情報を取得する。

取得した情報は以降のブランチ名生成と作業概要の提示に使用する。

### 3. ブランチ名の決定

- **Issue/PR がある場合**: タイトルからブランチ名を自動生成する
  - 形式: `f/<slug>-<number>`（例: `f/add-user-auth-123`）
  - slug はタイトルを小文字化し、英数字とハイフンのみに変換、30文字以内に切り詰める
- **Issue/PR がない場合**: 口頭指示の内容からブランチ名を提案し、ユーザーに確認する
- **ブランチ名が直接指定された場合**: そのまま使用する

### 4. main/master を最新化してブランチ作成

特別な指示がない限り、main/master を最新化してからブランチを切る。

```bash
git switch main
git pull --ff-only origin main
git switch -c <ブランチ名>
```

main への pull が fast-forward できない場合はユーザーに報告して停止する。

### 5. タスクドキュメント確認

worktree の場合は元ディレクトリの `.claude/docs` 配下に関連するタスクドキュメントがあるか確認する。

関連するドキュメントがあれば内容を提示する。

### 6. 作業概要の提示

以下を表示して作業開始の準備完了を報告する:

- 作成したブランチ名
- Issue/PR の概要（取得した場合）
- 確認すべきタスクドキュメント（存在する場合）

## 注意

- このスキルはブランチ作成と情報取得のみを行う。コードの変更は行わない
- Issue/PR の作成・更新は `/gh-edit` に委ねる
