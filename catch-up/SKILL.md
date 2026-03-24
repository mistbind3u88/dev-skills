---
name: catch-up
description: main の最新をrebaseで取り込み、コンフリクトを解消する。
allowed-tools: Bash(git fetch:*) Bash(git rebase:*) Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git rev-parse:*) Bash(git add:*) Bash(git checkout:*) Read
---

# catch-up スキル

現在のブランチに main の最新を rebase で取り込む。

## 手順

### 1. 前提確認

```bash
git status -s
git rev-parse --abbrev-ref HEAD
git log --oneline main..HEAD
```

- 未コミットの変更がある場合は先にコミットするようユーザーに伝える
- main ブランチにいる場合は `git pull --rebase` で済むので、そちらを案内する

### 2. リモートの最新を取得する

```bash
git fetch origin
```

### 3. バックアップブランチを作成する

rebase 前に `/backup-branch` で現在の状態を保存する。

### 4. rebase を開始する

```bash
git rebase origin/main
```

#### コンフリクトが発生した場合

1. コンフリクトファイルを確認する

```bash
git status -s
git diff --name-only --diff-filter=U
```

2. 各コンフリクトファイルを Read で開き、コンフリクトマーカーを確認する
3. ユーザーに解消方針を確認する（ユーザーが方針を示している場合はそれに従う）
4. Edit で解消する
5. 解消したファイルをステージして続行する

```bash
git add <解消したファイル>
git rebase --continue
```

コンフリクトが連続する場合はコミットごとに繰り返す。

### 5. rebase 完了後の確認

```bash
git log --oneline origin/main..HEAD
git status -s
```

rebase 前後でコミット数が変わっていないことを確認する。

## 注意

- rebase 中に `git rebase --abort` が必要な場合はユーザーに確認する
- コンフリクト解消時、どちらを優先するか不明な場合は必ずユーザーに確認する。勝手に判断しない
- rebase 後に force push が必要になるが、push はユーザーの指示があるまで行わない
- バックアップブランチは rebase が問題なく完了したことを確認するまで残す
