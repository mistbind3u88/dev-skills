---
name: gh-edit
description: GitHubのPRやIssueを作成・更新する。既存の内容を読み込んだ上で編集する。
allowed-tools: Bash(gh pr:*) Bash(gh issue:*) Bash(git log:*) Bash(git diff:*) Read
---

# gh-editスキル

GitHubのPRやIssueを作成・更新する。

**概要欄の読み手は人間である。** データの羅列や機械的な転記ではなく、読み手が背景・目的・判断ポイントを理解できる文章を書くこと。

## 手順

### 1. 引数から対象と操作を特定する

- PR番号、Issue番号、またはURLが指定される
- 未指定の場合は現在のブランチに紐づくPRを対象とする
- 対象が存在しない場合は新規作成する

### 2. テンプレートを探す

PR/Issueの概要欄を作成・更新する前に、テンプレートを探してReadツールで読み込む。

#### テンプレートの優先順位

1. **リポジトリのテンプレート**（優先）: `.github/PULL_REQUEST_TEMPLATE.md` や `.github/ISSUE_TEMPLATE/` 配下のファイルが存在すればそちらを使う
2. **スキル同梱のテンプレート**（フォールバック）: リポジトリにテンプレートがない場合、このスキルの配置ディレクトリにあるテンプレートを使う
   - PR: [AGENTS.md](AGENTS.md) と [pull_request_template.md](pull_request_template.md)
   - Issue: [issue_template.md](issue_template.md)

### 概要欄の構成

PR・Issueを問わず、概要欄の最初の見出しは `# Subject` とし、そのPR/Issueの趣旨を要約する文章を記載する。タイトルの繰り返しではなく、背景・目的・動機が伝わる内容にすること。

### 3-A. PRの新規作成

ブランチのコミット履歴からタイトルと概要を作成する。PRは常にdraftで作成する。

```bash
git log --oneline main..HEAD
git diff --stat main..HEAD

gh pr create --draft --title "<タイトル>" --body "$(cat <<'EOF'
<テンプレートに従って概要欄を作成>
EOF
)"
```

### 3-B. Issueの新規作成

Issueは直接 `gh issue create` するのではなく、まずローカルに設計ドキュメントを書き出し、レビューを経てから投稿する。

1. `.claude/docs/` 配下にIssueの概要欄となるドキュメントを作成する
2. ユーザーに内容を確認してもらう
3. 確認後、Issueを作成する

```bash
gh issue create --title "<タイトル>" --body "$(cat <<'EOF'
<ドキュメントの内容>
EOF
)"
```

### 3-C. 既存の更新

#### タイトルと概要欄を読み込む

```bash
# PRの場合
gh pr view <番号> --json title,body --jq '.title,.body'

# Issueの場合
gh issue view <番号> --json title,body --jq '.title,.body'
```

**重要**: 既存の内容を必ず確認してから編集すること。白紙から書き直さない。

#### コミット履歴を確認する（PRの場合）

```bash
git log --oneline main..HEAD
git diff --stat main..HEAD
```

#### タイトルと概要欄を更新する

既存の内容をベースに、修正・追記・削除を行う。
タイトルはコミット履歴と概要欄の内容を踏まえて、PR/Issueの現在のスコープを正確に反映しているか見直す。

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

- **概要欄は人間が読むものである。** 調査データや分析結果をそのまま貼るのではなく、読み手が順を追って理解できるように構成・要約すること
- コミットハッシュはコードブロック（`` ` ``）で囲まない。GitHub UI 上でコミットへのリンクとして自動認識させるため
- PRは常にdraftで作成する
- 既存の概要欄が空でない場合、必ず既存の内容を読み込んでから更新する
- ユーザーが明示的に全面書き換えを指示しない限り、既存の構造を維持する
- **概要欄を作成・更新する前に、必ずテンプレートを探して読み込むこと**（手順2参照）
- 内容を推測・記憶に頼らず、毎回実際に読み込んで確認する
