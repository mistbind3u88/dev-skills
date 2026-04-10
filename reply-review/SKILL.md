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

### 2. 対応種別を判定する

指摘に対してコードを変更して対応したか、対応不要と判断したかを判定する。

- **対応済み**: コードを変更して指摘に対応した → ステップ 3A へ
- **対応不要**: 対応しないと判断した → ステップ 3B へ

### 3A. 対応済みのリプライを作成する

対応コミットを推定し、リプライ本文を作成する。

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

#### リプライ本文

```markdown
対応しました。

- <変更内容の要約 1>
- <変更内容の要約 2>

該当コミット: {sha}
```

#### 対応内容の示し方

コミットか差分リンクのどちらか一方を使う。両方は不要。

- **fixup コミットが push 済みの場合**: fixup コミットのハッシュを記載する
- **autosquash / rebase でコミットを潰した場合**: force push 前後の HEAD の compare リンクを添付する。`..`（2ドット）を使うこと（`...` 3ドットではない）

```markdown
差分: https://github.com/{owner}/{repo}/compare/{old_sha}..{new_sha}
```

派生対応がある場合は見出しで整理する。

```markdown
対応しました。以下の付随対応も行いました。

### 1. <直接対応の説明> ({sha})

- <詳細>

### 2. <派生対応の説明> ({sha})

- <詳細>
```

### 3B. 対応不要のリプライを作成する

対応しない理由を簡潔に記載する。

```markdown
<対応不要の理由>
```

### 4. ユーザーに確認する

**ユーザーに確認**: 本文を表示し、追加・修正・削除の指示を受け付ける。承認されたら投稿に進む。

### 4. リプライを投稿する

```bash
gh api "repos/{owner}/{repo}/pulls/{pr}/comments" \
  -f body="<本文>" \
  -F in_reply_to={comment_id}
```

投稿後、コメントの URL を報告する。

### 5. Bot コメントにリアクションを返す

リプライ先がエージェント系レビュアー（Copilot、codex-connector 等）のコメントで、エモートによるフィードバックを求められている場合（「Useful? React with 👍 / 👎」等）、リプライ投稿と合わせてリアクションも付ける。

- 指摘に対応した場合: 👍 (`+1`)
- 対応不要と判断した場合: 👎 (`-1`)

```bash
gh api repos/{owner}/{repo}/pulls/comments/{comment_id}/reactions -f content="+1"
```

リプライとリアクションは必ずセットで行う。片方だけにならないよう注意する。

## 注意

- **対応コミットはpush済みであること**: 「対応しました」とリプライする場合、該当コミットがリモートにpush済みであることを確認する。push前のローカルコミットを参照してリプライすると、レビュアーがコミットリンクから変更内容を確認できない
- **二重投稿の防止**: ステップ1で既存リプライを確認し、自分（`git config user.name` または GitHub ログインユーザー）のリプライが既にある場合は投稿をスキップする。force push によりコメントが outdated 状態になっても、API 上はリプライが残っている場合があるため、`in_reply_to_id` でのフィルタに加えて投稿者名でも照合する
- コミットハッシュは7桁の短縮形で記載し、コードブロック（`` ` ``）で囲まない。GitHub UI 上でコミットへのリンクとして自動認識させるため
- リプライの読み手はレビュアー（人間）。diff の羅列ではなく、何をどう変えたかが伝わる要約にする
- owner/repo は `gh api repos/{owner}/{repo}` または `git remote get-url origin` から取得する
