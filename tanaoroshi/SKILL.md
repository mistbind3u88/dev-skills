---
name: tanaoroshi
description: 複数GitHubリポジトリのOpen Issue/PRを横断的に棚卸しする。body精査・コメント確認・前提作業の完了状況確認・テーマ別の構造整理・アクション提案を行う。
allowed-tools: Bash(go run ./skills/tanaoroshi:*), Bash(gh repo view:*), Agent, Read
---

# Issue/PR 棚卸し

複数の GitHub リポジトリを横断して Open Issue/PR を精査し、テーマ別に構造化した棚卸しレポートを作成する。

## 入力

引数でリポジトリをスペース区切りで指定する。

```
/tanaoroshi owner/repo1 owner/repo2
```

引数なしの場合は、カレントリポジトリに加えて CLAUDE.md または AGENTS.md の「関連リポジトリ」セクションに記載されたリポジトリも対象とする。関連リポジトリセクションがなければカレントリポジトリのみを対象とする。

## 出力

以下の構成の棚卸しレポートを会話に出力する。

1. テーマ別の構造整理（Open/Closed の依存関係ツリー）
2. 全体の進捗サマリー
3. アクション提案（マージ可能・クローズ可能・方針合意が必要なもの）

## 手順

### Phase 0: 対象リポジトリの決定

引数でリポジトリが指定されていない場合、以下の手順で対象を決定する。

1. `gh repo view --json nameWithOwner` でカレントリポジトリの `owner/repo` を取得する
2. CLAUDE.md または AGENTS.md の「関連リポジトリ」セクションからリポジトリ名を抽出する
3. カレントリポジトリ + 関連リポジトリを対象とする

### 無視リスト

`skills/tanaoroshi/ignore` にリストされた Issue/PR は `summary`・`refs` の出力から自動的に除外される。棚卸し対象外にしたい Issue/PR がある場合はこのファイルに追記する。形式は `.gitignore` のコメント規則に準拠し、`#` で始まる行はコメントとして扱う。

### Phase 1: データ収集

`collect` で各リポジトリの Open Issue/PR を一括取得し、ファイルに保存する。

```bash
go run ./skills/tanaoroshi collect <owner/repo> [owner/repo2 ...] > .results/tanaoroshi.json
```

レスポンスはリポジトリごとに `issues` と `prs` を持つ JSON オブジェクト（body 含む全フィールド）。

データが大きい場合は `summary` で body を除いた一覧を取得する。

```bash
go run ./skills/tanaoroshi summary .results/tanaoroshi.json
```

### Phase 2: 参照先の解決

`refs` で collect 結果から body 内の参照を自動抽出する（重複排除済み）。

```bash
go run ./skills/tanaoroshi refs .results/tanaoroshi.json
```

出力は `{source, ref}` のペア配列。検出する参照パターン:

- `#123` — 同一リポジトリ内の参照
- `owner/repo#123` — クロスリポジトリ参照
- `closes #123`, `fixes #123` — クローズ対象の明示
- `https://github.com/owner/repo/issues/123` — URL 形式の参照
- `https://github.com/owner/repo/pull/123` — URL 形式の参照

抽出した参照のうち Open 一覧にないものを `resolve` でステータス確認する。
引数では `#` の代わりに `:` を区切り文字として使用する（シェルでのコメント解釈を回避するため）。

```bash
go run ./skills/tanaoroshi resolve \
  --issues <owner/repo:N> [owner/repo:N ...] \
  --prs <owner/repo:N> [owner/repo:N ...]
```

### Phase 2a: 直近コメントの確認

直近で更新のあった Open Issue/PR について、最新のコメントを確認する。

1. `summary` の出力から `updatedAt` が直近2週間以内の Issue/PR を抽出する
2. 該当する Issue/PR のコメントを `comments` サブコマンドで取得する

```bash
go run ./skills/tanaoroshi comments <owner/repo:N> [owner/repo:N ...]
```

- Issue コメントと PR レビューコメントを時系列でマージして返す
- 出力の各コメントには `type` フィールド（`"comment"` = Issue コメント、`"review"` = PR レビューコメント）が付与される
- コメント内容は Phase 3 のテーマ整理で「進捗」「ブロッカー」「方針変更」の判断材料にする

### Phase 2b: リリース PR の分離

テーマ整理の前に、リリース PR を識別して分離する。以下のいずれかに該当する PR はリリース PR として扱う:

- タイトルに `release` / `Release` / `v0.0.0` 等のバージョン番号パターンを含む
- 作成者が bot（`dependabot`, `renovate`, `github-actions`, `release-please` 等）
- ラベルに `release` / `autorelease` を含む

リリース PR は Phase 3 のテーマ整理に含めず、Phase 4 のサマリーで「リリース運用」として別セクションにまとめる。

### Phase 3: テーマ別構造整理

収集したデータを以下の観点で分析し、テーマ（上位目標）ごとにグループ化する。

#### 3-1. グループ化の基準

- **body 内の相互参照**: Issue/PR が互いに言及し合っている場合は同一テーマ
- **closes/fixes 関係**: PR が Issue を closes している場合は親子関係
- **ブランチ名の共通プレフィックス**: 同一機能の連続 PR
- **タイトルの共通キーワード**: 同一施策の Issue 群
- **クロスリポジトリの対応関係**: 同名ブランチ、同一テーマの Issue が複数リポジトリに存在

#### 3-2. 各テーマの出力形式

テーマごとに以下を整理する。Open を先に、Closed を後に記載する:

```
## N. テーマ名

**ゴール**: このテーマで達成しようとしていること（1-2文）

  [PR/OPEN] タイトル @author → レビュワー: レビュー [owner/repo#125](https://github.com/owner/repo/issues/125)
    └─ [PR/Draft] タイトル @author → 作者: 下書き完成 [owner/repo#126](https://github.com/owner/repo/issues/126)
  [Issue/OPEN] タイトル @author → @assignee: 実装 [owner/repo#127](https://github.com/owner/repo/issues/127)
  [PR/CLOSED] タイトル（完了日） [owner/repo#123](https://github.com/owner/repo/issues/123)
    ├─ [Issue/CLOSED] タイトル（完了日） [owner/repo#124](https://github.com/owner/repo/issues/124)

**進捗**: 要約と次のアクション
```

- 行頭の種別プレフィックスは `[Issue/OPEN]` / `[Issue/CLOSED]` / `[PR/OPEN]` / `[PR/Draft]` / `[PR/CLOSED]` のいずれかを必ず付ける（merged PR は `[PR/CLOSED]` として扱う）
- Open の PR はタイトル直後に `@author` を付ける。Issue は assignee がいる場合 `@author → @assignee` と書く
- Open の行末（リンクの前）には `→ <主体>: <アクション>` を付ける。判定ルールは末尾「ネクストアクション主体の判定ルール」節を参照
- Closed は主体注釈を付けない（完了済みのため）
- Open は詳しく（タイトル・残件内容・ブロッカーの有無・`owner/repo#N`）
- Closed は簡潔に（タイトル・完了日・`owner/repo#N`）
- リンク（`owner/repo#N`）は行末に置く

#### 3-3. テーマに属さない個別 Issue

テーマにまとめられない独立した Issue/PR は以下のカテゴリで分類する:

- **インフラ・運用改善**: デプロイ、監視、コスト管理
- **コード品質**: リファクタ、リネーム、設定整理
- **管轄外の可能性**: 他リポジトリや他チームの対応が必要なもの

### Phase 4: 全体サマリーとアクション提案

#### 4-1. 進捗サマリー

テーマごとの Open 数・Closed 数・進捗状況を表形式で出力する。

| テーマ   | Open | Closed | 進捗 |
| -------- | :--: | :----: | ---- |
| テーマ名 |  N   |   M    | 要約 |

#### 4-2. アクション提案

以下の 3 カテゴリに分けて提案する。各箇条書きには `[主体]` タグ（`[マージ担当]` / `[レビュワー]` / `[作者]` / `[@assignee]` / `[チーム]` / `[未アサイン]` 等）を前置し、その後に `[Issue/…]` または `[PR/…]` と `@author`・リンクを続ける。例:

```
- [マージ担当] [PR/OPEN] タイトル @author [owner/repo#N](https://github.com/owner/repo/issues/N)
- [レビュワー] [PR/OPEN] タイトル @author [owner/repo#N](https://github.com/owner/repo/issues/N)
- [作者] [PR/OPEN] タイトル @author [owner/repo#N](https://github.com/owner/repo/issues/N)
```

**すぐ対応できるもの**:

- `[マージ担当]` APPROVED 済みでマージ可能な PR
- `[クローズ担当]` 前提作業の完了によりクローズ可能な Issue
- `[クローズ担当]` 上位 Issue に統合してクローズ可能な Issue

**着手可能になったもの**:

- `[@assignee]`/`[未アサイン]` 前提の Closed により実装基盤が整った Issue/PR
- `[作者]` Draft PR のうち、依存がすべて解決済みのもの

**方針合意が必要なもの**:

- `[チーム]` 議論 Issue で方針未決のもの
- `[チーム]` クロスリポジトリで方針が対立しているもの
- `[オーナー要判断]` 長期放置（3ヶ月以上）で必要性の再評価が必要なもの

## 出力の共通ルール

- **Issue/PR への言及は、タイトル末尾・本文中・括弧内を問わず、すべて `[owner/repo#N](https://github.com/owner/repo/issues/N)` 形式のクリッカブルリンクで記載する。`#N` や `owner/repo#N` だけの素の表記ではターミナル上でリンクにならない。** 範囲指定（`#256〜#258`）や列挙（`#204, #205`）も個別に完全形式で書く
- **種別プレフィックス** `[Issue/OPEN]` / `[Issue/CLOSED]` / `[PR/OPEN]` / `[PR/Draft]` / `[PR/CLOSED]` を必ず行頭に付ける（merged PR は `[PR/CLOSED]`）
- **Open の PR** はタイトル直後に `@author` を付ける。**Open の Issue** は assignee がいる場合 `@author → @assignee`、いない場合 `@author` を付ける
- **Open の行末** には `→ <主体>: <アクション>` を付ける（判定ルールは下の「ネクストアクション主体の判定ルール」を参照）。Closed は付けない
- リンクは行末に置く
- テーマ別ツリー、個別 Issue 一覧、サマリー表、アクション提案のいずれでも同じ形式を適用する
- Open を先に、Closed を後に記載する
- PR の状態は Closed に統一する（Merged も Closed として扱う）

## ネクストアクション主体の判定ルール

Open の Issue/PR 各行の行末に付ける `→ <主体>: <アクション>` は以下の基準で決める。判定に必要な `isDraft` / `reviewDecision` / `assignees` / ラベルは `summary` の出力に含まれる。必要なら `comments` で直近コメントの向き先を補足する。

### PR

**主判定**（`summary` だけで決定する。上の行から評価し、最初に一致した行を採用する）:

| 優先 | 条件                                                | 主体: アクション       |
| ---- | --------------------------------------------------- | ---------------------- |
| 1    | `isDraft=true`                                      | `作者: 下書き完成`     |
| 2    | `reviewDecision="APPROVED"`                         | `マージ担当: マージ可` |
| 3    | `reviewDecision="CHANGES_REQUESTED"`                | `作者: 指摘対応`       |
| 4    | `reviewDecision="REVIEW_REQUIRED"` または空・未設定 | `レビュワー: レビュー` |

**補足判定**（任意）: `comments` を取得済みの場合に限り、以下で主判定の結果を上書きしてよい。判別できない場合は主判定の結果をそのまま使う。

- 直近コメントの投稿者が作者以外 → `作者: 再プッシュ`
- 直近コメントの投稿者が作者 → `レビュワー: 再レビュー`

### Issue

上の行から評価し、最初に一致した行を採用する（`assignees` が優先）。

| 優先 | 条件                                  | 主体: アクション         |
| ---- | ------------------------------------- | ------------------------ |
| 1    | `assignees` あり                      | `@assignee: 実装/調査`   |
| 2    | 議論系ラベル（`discussion` 等）を持つ | `チーム: 方針合意`       |
| 3    | 上記以外                              | `未アサイン: トリアージ` |

### アクション提案（Phase 4-2）の主体タグ

アクション提案セクションでは行頭に `[主体]` タグを付ける（上記の「主体」を `[…]` で囲んだ形）。代表例:

- `[マージ担当]` — APPROVED 済み PR のマージを待っている
- `[レビュワー]` — レビュー・再レビューを待っている
- `[作者]` — 下書き完成・指摘対応・再プッシュを待っている
- `[@assignee]` — 指名済みの担当者の実装/調査を待っている
- `[チーム]` — 方針合意を待っている
- `[クローズ担当]` — クローズ判断を待っている
- `[オーナー要判断]` — 長期放置で必要性の再評価が必要

## 注意

- Closed の前提作業は「何が完了したか」の事実のみ記載し、詳細な説明は省く
- リリース PR は Phase 2b で分離し、個別テーマに含めない。Phase 4 のサマリーで「リリース運用」として別セクションにまとめる
