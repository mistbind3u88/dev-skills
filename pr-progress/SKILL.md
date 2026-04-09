---
name: pr-progress
description: PR のトップレベルに進捗コメントを投稿する。force push 後の差分報告やコンフリクト解消報告など、PR の経過を記録する。
allowed-tools: Bash(gh pr comment:*) Bash(gh pr view:*) Bash(git rev-parse:*) Bash(gh repo view:*) Read
---

# pr-progress スキル

PR のトップレベルに進捗コメントを投稿する。

## 対象場面

| 場面                         | トリガー元                       |
| ---------------------------- | -------------------------------- |
| コンフリクト解消後の差分報告 | `/push`（force push 時）         |
| autosquash 後の差分報告      | `/ship`（autosquash のみの場合） |
| catch-up 後の差分報告        | `/catch-up`                      |

## 手順

### 1. PR の存在を確認する

```bash
gh pr view --json number,isDraft --jq '{number, isDraft}'
```

- PR が存在しない場合はスキップして終了する
- draft の PR にはコメントしない

### 2. コメント本文を作成する

呼び出し元から以下の情報を受け取る:

- **旧 HEAD のハッシュ**（バックアップブランチまたは push 前に記録した値）
- **新 HEAD のハッシュ**（現在の HEAD）
- **操作の種類**（コンフリクト解消 / autosquash / catch-up）
- **補足情報**（解消したファイル一覧、解消方針など。操作の種類による）

compare リンクは2ドット（`..`）を使う:

```
https://github.com/<owner/repo>/compare/<旧HEAD>..<新HEAD>
```

リポジトリ名・PR 番号・ハッシュは会話コンテキストから直接埋め込む。

#### コンフリクト解消の場合

```markdown
main を rebase で取り込み、コンフリクトを解消しました。

- 解消したファイル: <ファイル一覧>
- 解消方針: <簡潔な説明>
- 差分: https://github.com/<owner/repo>/compare/<旧HEAD>..<新HEAD>
```

#### autosquash の場合

```markdown
fixup コミットを autosquash で整理しました。

- 差分: https://github.com/<owner/repo>/compare/<旧HEAD>..<新HEAD>
```

#### catch-up の場合

```markdown
main の最新を rebase で取り込みました。

- 差分: https://github.com/<owner/repo>/compare/<旧HEAD>..<新HEAD>
```

### 3. ユーザーに確認する

コメント本文を提示し、承認を得てから投稿する。

### 4. コメントを投稿する

```bash
gh pr comment <PR番号> --body "<本文>"
```

投稿後、コメントの URL を報告する。

## 注意

- draft PR にはコメントしない
- コミットハッシュは GitHub UI 上でコミットへのリンクとして自動認識させる。そのためコードブロック（`` ` ``）で囲まず、両側を空白または改行にする（括弧内に入れない）
- fixup の修正理由をコメントする場合は、対象の fixup コミットハッシュを必ず含める
- 読み手はレビュアー（人間）。操作の経緯が伝わる簡潔な記述にする
