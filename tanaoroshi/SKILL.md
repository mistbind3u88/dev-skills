---
name: tanaoroshi
description: 複数GitHubリポジトリのOpen Issue/PRを横断的に棚卸しする。body精査・前提作業の完了状況確認・テーマ別の構造整理・アクション提案を行う。
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

  [OPEN] タイトル owner/repo#125
    └─ [OPEN/Draft] タイトル owner/repo#126
  [CLOSED] タイトル（完了日） owner/repo#123
    ├─ [CLOSED] タイトル（完了日） owner/repo#124

**進捗**: 要約と次のアクション
```

- Open は詳しく（タイトル・残件内容・ブロッカーの有無・`owner/repo#N`）
- Closed は簡潔に（タイトル・完了日・`owner/repo#N`）
- リンク（`owner/repo#N`）はタイトルの後に置く

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

以下の 3 カテゴリに分けて提案する。

**すぐ対応できるもの**:

- APPROVED 済みでマージ可能な PR
- 前提作業の完了によりクローズ可能な Issue
- 上位 Issue に統合してクローズ可能な Issue

**着手可能になったもの**:

- 前提の Closed により実装基盤が整った Issue/PR
- Draft PR のうち、依存がすべて解決済みのもの

**方針合意が必要なもの**:

- 議論 Issue で方針未決のもの
- クロスリポジトリで方針が対立しているもの
- 長期放置（3ヶ月以上）で必要性の再評価が必要なもの

## 出力の共通ルール

- Issue/PR への言及は必ず `owner/repo#123` 形式で記載する。`#123` だけではリンクにならない
- リンクはタイトルの後に置く
- テーマ別ツリー、個別 Issue 一覧、サマリー表、アクション提案のいずれでも同じ形式を適用する
- Open を先に、Closed を後に記載する
- PR の状態は Closed に統一する（Merged も Closed として扱う）

## 注意

- Closed の前提作業は「何が完了したか」の事実のみ記載し、詳細な説明は省く
- リリース PR（bot 自動生成）は個別テーマに含めず、リリース運用として別セクションにまとめる
