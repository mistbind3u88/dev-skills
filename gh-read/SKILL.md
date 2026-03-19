---
name: gh-read
description: PR や Issue の情報（タイトル・概要・レビュー状態・差分量など）を取得する。作業開始時の要件把握や PR の状態確認に使う。
allowed-tools: Bash(gh pr view:*) Bash(gh issue view:*) Bash(gh api:*)
---

# gh-read スキル

GitHub の PR や Issue の主要な情報を取得し、JSON で出力する。

## 手順

### 1. 対象を特定する

`$ARGUMENTS` から PR 番号、Issue 番号、または URL を受け取る。

- `#123` や数字のみの場合は、コンテキストから PR か Issue かを判断する（不明ならユーザーに確認）
- URL の場合はパスから PR (`/pull/`) か Issue (`/issues/`) かを判別する

### 2. 情報を取得する

#### PR の場合

```bash
gh pr view <number> --json number,url,title,body,state,author,baseRefName,headRefName,labels,isDraft,reviewDecision,additions,deletions,changedFiles,files,closingIssuesReferences
```

取得フィールド:

| フィールド | 内容 |
|-----------|------|
| `number`, `url` | 識別情報 |
| `title`, `body` | タイトルと概要欄 |
| `state` | open / closed / merged |
| `author` | 作成者 |
| `baseRefName`, `headRefName` | ベースブランチとヘッドブランチ |
| `labels` | ラベル |
| `isDraft` | ドラフト状態 |
| `reviewDecision` | レビュー状態（APPROVED / CHANGES_REQUESTED / REVIEW_REQUIRED） |
| `additions`, `deletions`, `changedFiles` | 差分の量 |
| `files` | 変更ファイルリスト（パス・追加行・削除行） |
| `closingIssuesReferences` | リンク済み Issue |

#### Issue の場合

```bash
gh issue view <number> --json number,url,title,body,state,author,labels,assignees
```

取得フィールド:

| フィールド | 内容 |
|-----------|------|
| `number`, `url` | 識別情報 |
| `title`, `body` | タイトルと概要欄 |
| `state` | open / closed |
| `author` | 作成者 |
| `labels` | ラベル |
| `assignees` | アサイン |

Issue にリンク済みの PR を取得する場合:

```bash
gh api graphql -f query='
  query($owner: String!, $repo: String!, $number: Int!) {
    repository(owner: $owner, name: $repo) {
      issue(number: $number) {
        timelineItems(itemTypes: [CROSS_REFERENCED_EVENT], first: 10) {
          nodes {
            ... on CrossReferencedEvent {
              source {
                ... on PullRequest { number title state url }
              }
            }
          }
        }
      }
    }
  }' -f owner=<owner> -f repo=<repo> -F number=<number>
```

### 3. 結果を出力する

取得した JSON をそのまま出力する。

## 注意

- このスキルは情報の取得と出力のみを行う。PR や Issue の編集は `/gh-edit` に委ねる
- コメントやレビューコメントの収集は `/collect-feedback` に委ねる
