---
name: gh-edit
description: GitHub の PR や Issue を作成・更新する。既存の内容を読み込んだ上で編集する。
allowed-tools: Bash(gh pr:*) Bash(gh issue:*) Bash(git log:*) Bash(git diff:*) Read
---

# gh-edit スキル

GitHub の PR や Issue を作成・更新する。

## 手順

### 1. `$ARGUMENTS` から対象と操作を特定する

- PR番号、Issue番号、またはURLが指定される
- 未指定の場合は現在のブランチに紐づくPRを対象とする
- 対象が存在しない場合は新規作成する

### 2-A. 新規作成

#### PR の作成

ブランチのコミット履歴からタイトルと概要を作成する。PR は常に draft で作成する。

```bash
# コミット履歴を確認
git log --oneline main..HEAD
git diff --stat main..HEAD

# draft で作成
gh pr create --draft --title "<タイトル>" --body "$(cat <<'EOF'
<TEMPLATE.md のスタイルに従って概要欄を作成>
EOF
)"
```

#### Issue の作成

```bash
gh issue create --title "<タイトル>" --body "$(cat <<'EOF'
<内容>
EOF
)"
```

### 2-B. 既存の更新

#### 概要欄を読み込む

```bash
# PRの場合
gh pr view <番号> --json body --jq .body

# Issueの場合
gh issue view <番号> --json body --jq .body
```

**重要**: 既存の内容を必ず確認してから編集すること。白紙から書き直さない。

#### ユーザーの指示に基づいて更新する

既存の内容をベースに、修正・追記・削除を行う。

```bash
# PRの場合
gh pr edit <番号> --body "$(cat <<'EOF'
<既存の内容をベースに更新>
EOF
)"

# Issueの場合
gh issue edit <番号> --body "$(cat <<'EOF'
<既存の内容をベースに更新>
EOF
)"
```

## 注意

- PR は常に draft で作成する
- 既存の概要欄が空でない場合、必ず既存の内容を読み込んでから更新する
- ユーザーが明示的に全面書き換えを指示しない限り、既存の構造を維持する
- PR 概要欄のスタイルは同ディレクトリの `TEMPLATE.md` に従う
