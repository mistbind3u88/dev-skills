---
name: reply-review
description: PR レビューコメントへの対応内容をスレッドにリプライする。対応コミットを自動推定し、ユーザー確認後に投稿する。
allowed-tools: Bash(gh api:*) Bash(git log:*) Bash(git diff:*) Bash(git show:*) Bash(git rev-parse:*)
---

# reply-review スキル

PR レビューコメントへの対応内容を、該当スレッドにリプライとして投稿する。

## 手順

### 1. レビューコメントを取得する

`$ARGUMENTS` から PR レビューコメントの URL を受け取り、comment ID を抽出する。

```
# URL 形式
https://github.com/{owner}/{repo}/pull/{pr}#discussion_r{comment_id}
```

```bash
# コメント本文・対象ファイル・既存リプライを取得
gh api "repos/{owner}/{repo}/pulls/comments/{comment_id}" --jq '{body, path, line, created_at}'

# 既存リプライを確認（二重投稿を防ぐ）
gh api "repos/{owner}/{repo}/pulls/{pr}/comments" --jq '[.[] | select(.in_reply_to_id == {comment_id}) | {body, created_at}]'
```

### 2. 対応コミットを推定する

指摘内容と対象ファイルを手がかりに、PR 内のコミットから対応したものを特定する。

```bash
# PR 内の全コミットを取得
git log --oneline main..HEAD

# 対象ファイルを変更したコミットを絞り込む
git log --oneline main..HEAD -- {path}

# 各コミットの diff を確認し、指摘内容との関連を判定
git show --stat {sha}
git diff {sha}~1 {sha} -- {path}
```

#### 派生対応の検出

指摘への直接対応だけでなく、そこから派生した変更も検出する。

- 直接対応コミットで変更された関数・型が、他のコミットでも変更されていないか確認する
- 同じファイルや関連ファイル（DTO、テスト等）への変更を追跡する

### 3. リプライ本文を作成し、ユーザーに確認する

以下の形式でリプライ本文を作成し、投稿前にユーザーに提示する。

```markdown
対応しました。

- <変更内容の要約 1>
- <変更内容の要約 2>

該当コミット: {sha1}
```

派生対応がある場合は見出しで整理する。

```markdown
対応しました。以下の付随対応も行いました。

### 1. <直接対応の説明> ({sha})
- <詳細>

### 2. <派生対応の説明> ({sha})
- <詳細>
```

**ユーザーに確認**: 本文を表示し、追加・修正・削除の指示を受け付ける。承認されたら投稿に進む。

### 4. リプライを投稿する

```bash
gh api "repos/{owner}/{repo}/pulls/{pr}/comments" \
  -f body="<本文>" \
  -F in_reply_to={comment_id}
```

投稿後、コメントの URL を報告する。

## 注意

- 既存リプライがある場合は二重投稿にならないよう確認する
- コミットハッシュは短縮形（7-9文字）で記載し、コードブロック（`` ` ``）で囲まない。GitHub UI 上でコミットへのリンクとして自動認識させるため
- リプライの読み手はレビュアー（人間）。diff の羅列ではなく、何をどう変えたかが伝わる要約にする
- owner/repo は `gh api repos/{owner}/{repo}` または `git remote get-url origin` から取得する
