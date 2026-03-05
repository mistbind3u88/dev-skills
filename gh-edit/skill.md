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

`$SKILL_DIR/AGENTS.md` と `$SKILL_DIR/TEMPLATE.md` を Read ツールで読み込み、スタイル規則と概要欄の構造を確認する。
ブランチのコミット履歴からタイトルと概要を作成する。PR は常に draft で作成する。

```bash
# コミット履歴を確認
git log --oneline main..HEAD
git diff --stat main..HEAD

# draft で作成
gh pr create --draft --title "<タイトル>" --body "$(cat <<'EOF'
<AGENTS.md と TEMPLATE.md に従って概要欄を作成>
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

#### AGENTS.md と TEMPLATE.md を読み込む

概要欄の編集前に、必ず `$SKILL_DIR/AGENTS.md` と `$SKILL_DIR/TEMPLATE.md` を Read ツールで読み込む。

#### タイトルと概要欄を読み込む

```bash
# PRの場合
gh pr view <番号> --json title,body --jq '.title,.body'

# Issueの場合
gh issue view <番号> --json title,body --jq '.title,.body'
```

**重要**: 既存の内容を必ず確認してから編集すること。白紙から書き直さない。

#### コミット履歴を確認する

```bash
git log --oneline main..HEAD
git diff --stat main..HEAD
```

#### タイトルと概要欄を更新する

既存の内容をベースに、修正・追記・削除を行う。
タイトルはコミット履歴と概要欄の内容を踏まえて、PR/Issue の現在のスコープを正確に反映しているか見直す。

```bash
# PRの場合
gh pr edit <番号> --title "<タイトル>" --body "$(cat <<'EOF'
<既存の内容をベースに更新>
EOF
)"

# Issueの場合
gh issue edit <番号> --title "<タイトル>" --body "$(cat <<'EOF'
<既存の内容をベースに更新>
EOF
)"
```

## 注意

- PR は常に draft で作成する
- 既存の概要欄が空でない場合、必ず既存の内容を読み込んでから更新する
- ユーザーが明示的に全面書き換えを指示しない限り、既存の構造を維持する
- **PR 概要欄を作成・更新する前に、必ず `$SKILL_DIR/AGENTS.md` と `$SKILL_DIR/TEMPLATE.md` を Read ツールで読み込むこと**
- 内容を推測・記憶に頼らず、毎回実際に読み込んで確認する
