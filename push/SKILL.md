---
name: push
description: lint・build・test・codex review の通過を確認してから git push する。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git rev-parse:*) Bash(git push:*) Bash(gh pr comment:*) Bash(gh pr view:*)
---

# push スキル

lint・build・test が通り、codex review 済みであることを確認してから push する。

## 手順

### 1. 前提確認

```bash
git status -s
git log --oneline main..HEAD
git rev-parse --abbrev-ref HEAD
```

- 未コミットの変更がある場合は push せず、先にコミットするようユーザーに伝える
- main ブランチにいる場合は警告してユーザーに確認する

### 2. チェックを実行する

`/check` を実行して、未チェック項目を一通り実行する。

`$ARGUMENTS` に `--skip-review` がある場合は `/check --skip-review` で渡す。

全チェックが OK でない場合は push せずに停止する。

### 3. push する

```bash
git push -u origin HEAD
```

push 後、結果を報告する。

### 4. コンフリクト解消時の PR コメント

rebase で main を取り込んでコンフリクトを解消した場合（`--force-with-lease` で push した場合）、PR にコメントを残す。

1. 現在のブランチに紐づく PR があるか確認する

```bash
gh pr view --json number --jq '.number'
```

2. force push 前後のコミットハッシュから compare リンクを作成する

```bash
# push 前に旧 HEAD を記録しておく
OLD_HEAD=$(git rev-parse HEAD)
# push 後の新 HEAD
NEW_HEAD=$(git rev-parse HEAD)
# リポジトリの owner/repo を取得
REPO=$(gh repo view --json nameWithOwner --jq '.nameWithOwner')
# compare リンク: https://github.com/<owner/repo>/compare/<old>...<new>
```

3. PR が存在すれば、コンフリクト解消の旨をコメントする

```bash
gh pr comment <PR番号> --body "mainをrebaseで取り込み、コンフリクトを解消しました。

- 解消したファイル: <コンフリクトがあったファイル一覧>
- 解消方針: <どのように解消したかの簡潔な説明>
- 差分: https://github.com/<owner/repo>/compare/<旧HEAD>...<新HEAD>"
```

PR が存在しない場合はスキップする。

## 注意

- `$ARGUMENTS` で `--force` が指定された場合は `git push --force-with-lease` を使う（ユーザーに確認後）
