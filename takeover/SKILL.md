---
name: takeover
description: 前セッションのコンテキスト（Issue、プラン、タスクドキュメント、コミット履歴）を収集し、作業の引き継ぎを受ける。
allowed-tools: Bash(git log:*) Bash(git status:*) Bash(git rev-parse:*) Bash(git diff:*) Bash(git merge-base:*) Bash(grep:*) Read
---

# takeover スキル

前セッションの作業コンテキストを収集し、引き継ぎを受ける。worktree での作業再開を主な用途とする。

## 引数

```
$ARGUMENTS: <補足情報（任意）>
```

引数は任意。Issue 番号や追加の指示があれば受け取るが、なくてもブランチ名から自動推定する。

## 手順

### 1. 現在の状態確認

```bash
git rev-parse --abbrev-ref HEAD
git status -s
git log --oneline -10
```

- 現在のブランチ名を取得する
- 未コミットの変更があれば記録する
- 直近のコミット履歴を取得する

### 2. Issue/PR の特定と取得

ブランチ名末尾の数字を Issue 番号として推定する（例: `f/dare10-blacklist-feed-33423` → `#33423`）。引数で明示的に指定されていればそちらを優先する。

Issue 番号が特定できた場合、スキル `/gh-read` で情報を取得する。

### 3. プランファイルの検索

`~/.claude/plans/` 配下のファイルから、ブランチ名・Issue 番号をキーワードにして関連するプランを検索する。

```bash
grep -rl "<キーワード>" ~/.claude/plans/
```

見つかったプランファイルの内容を読み込む。

### 4. タスクドキュメントの確認

リポジトリの元ディレクトリ（`~/Workspace/<repo>`）の `.claude/docs` 配下に関連するタスクドキュメントがあるか確認する。

### 5. ブランチの差分概要

main/master からの差分を把握する。

```bash
git merge-base HEAD main
git diff --stat <merge-base>..HEAD
```

### 6. 作業状況の提示

収集した情報をまとめて提示する:

- **ブランチ**: ブランチ名、worktree パス
- **Issue/PR**: 概要とステータス
- **実装計画**: プランファイルの内容
- **進捗**: main からのコミット履歴と差分概要
- **未コミットの変更**: あれば内容を提示
- **タスクドキュメント**: あれば内容を提示
- **次のアクション**: プランの進捗状況から次に着手すべき作業を提案
