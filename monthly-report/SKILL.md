---
name: monthly-report
description: GitHub の活動データから月次活動報告書を作成する。テーマ軸で整理し、同僚向けの説明として構成する。
allowed-tools: Bash(gh search:*) Bash(gh api:*) Bash(date:*) Read Write Edit
---

# monthly-report スキル

指定した GitHub org での月次活動を集計し、報告書を作成する。

## 引数

```
$ARGUMENTS: <org> <YYYY-MM> [出力ファイルパス]
```

- `<org>`: GitHub organization 名（必須）
- `<YYYY-MM>`: 対象月（必須）
- `[出力ファイルパス]`: 省略時はユーザーに確認する

## 手順

### 1. データ収集

以下の 4 つの検索を並列実行する。

```bash
# 当月作成の Issue
gh search issues --author=@me --owner=<org> --created=<YYYY-MM-01>..<YYYY-MM-末日> --json number,title,repository,state,closedAt,createdAt --limit 200

# 当月作成の PR
gh search prs --author=@me --owner=<org> --created=<YYYY-MM-01>..<YYYY-MM-末日> --json number,title,repository,state,closedAt,createdAt --limit 200

# 前月以前作成で当月中に merge/close された PR
gh search prs --author=@me --owner=<org> --merged-at=<YYYY-MM-01>..<YYYY-MM-末日> --created=..<前月末日> --json number,title,repository,state,closedAt,createdAt --limit 200

# 前月以前作成で当月中に close された Issue
gh search issues --author=@me --owner=<org> --closed=<YYYY-MM-01>..<YYYY-MM-末日> --created=..<前月末日> --json number,title,repository,state,closedAt,createdAt --limit 200
```

### 2. 月末時点のステータスを特定する

- 各 PR/Issue の `closedAt` を確認する
- `closedAt` が月末より後、または `0001-01-01`（open）のものは月末時点では open
- `closedAt` と `state` の組み合わせで merged/closed/open を判定する（`gh search prs` に `mergedAt` フィールドは存在しない）
- 月末時点の merged/closed/open を集計する

### 3. テーマを抽出・構成する

PR/Issue のタイトルと対象リポジトリから主要テーマを特定する。

- PR 数が多い・影響範囲が広いテーマを「主な作業」として先頭に配置
- 関連するが独立したテーマ（CI 改善、運用対応等）を次に配置
- 単発の対応は「その他」にまとめる

### 4. 報告書を作成する

以下の構成で作成する。

#### 4-1. 概況（月末時点）

Issue/PR 数とステータス内訳、前月からの持ち越し。

#### 4-2. 何をしていたか

テーマ別の活動説明。各テーマで以下を記述する:

- 何をしたか
- なぜそうしたか
- 結果どうなったか
- 月末時点の状況

#### 4-3. 時期別の活動

時系列での PR/Issue 一覧（補足的な位置づけ）。月末に open/未 merge の PR には `→{merge日} merge` や `月末時点 open` の注記をつける。

#### 4-4. 主な対象リポジトリ

リポジトリごとの PR 数と主な内容。

### 5. ユーザーに確認する

作成した報告書を提示し、修正指示を受け付ける。

## 注意

- 読み手は「同じ org で日常的に開発している同僚」を想定する。前提知識はあるが、この月に何をしていたかは知らない人
- テーマ軸で整理する。時系列での記述は補足に留める
- 前月からの持ち越しや翌月への引き継ぎは明示する
- `gh search prs` の `--merged` フラグは bool 型。日付範囲には `--merged-at` を使う
- `gh search` の json フィールドに `mergedAt` は存在しない
