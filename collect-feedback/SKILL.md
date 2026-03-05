---
name: collect-feedback
description: 開発中のPR/Issueからコメントを収集し、要対応事項・検討事項・ドキュメント化すべき事項を整理して報告する。
allowed-tools: Bash(gh api:*) Bash(gh pr:*) Bash(gh issue:*) Bash(git log:*) Read
---

# collect-feedback スキル

PR/Issue のコメントを収集し、対応が必要な事項を整理する。

## 手順

### 1. 対象の PR/Issue を特定する

`$ARGUMENTS` で PR/Issue 番号が指定されていればそれを使う。未指定なら現在のブランチに紐づく PR を対象とする。

```bash
# 現在のブランチの PR を取得
gh pr view --json number,title --jq '{number, title}'
```

### 2. コメントを収集する

PR のレビューコメントと通常コメントを取得する。

```bash
# レビューコメント（コード上のコメント）
gh api repos/{owner}/{repo}/pulls/{number}/comments --paginate

# PR の通常コメント
gh api repos/{owner}/{repo}/issues/{number}/comments --paginate
```

### 3. コメントを分類・整理する

収集したコメントを以下のカテゴリに分類する。

#### 要対応

- 明示的な修正依頼や指摘
- 合意された変更方針で未実装のもの
- バグや問題の報告

#### 検討事項

- 意見が分かれたまま結論が出ていない議論
- question タグ付きで回答後にアクションが不明なもの

#### ドキュメント化すべき事項

- note タグ付きの設計判断や注意事項
- 将来の対応方針（「別PRで対応」「次のフェーズで」等）
- コードコメントやAGENTS.mdに残すべき知見

#### 対応済み

- 既にコードに反映された指摘
- 議論の結果対応不要と合意されたもの

### 4. ユーザーに報告する

分類結果をユーザーに報告する。各項目について以下を含める。

- コメントの要約
- 誰のコメントか
- 対象ファイル・行番号（レビューコメントの場合）
- 現在の対応状況の判断根拠
