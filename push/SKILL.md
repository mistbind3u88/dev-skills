---
name: push
description: lint・build・test・codex review の通過を確認してから git push する。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git rev-parse:*) Bash(git push:*)
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

## 注意

- `$ARGUMENTS` で `--force` が指定された場合は `git push --force-with-lease` を使う（ユーザーに確認後）
