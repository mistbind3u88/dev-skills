---
name: codex-review
description: codex CLI を使って変更差分のコードレビューを実行する。設計完了後や実装・テスト完了後のレビューに使う。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(codex review:*)
---

# codex review スキル

変更差分を codex CLI でレビューに出す。

## 手順

1. git の状態を確認し、レビュー対象の変更範囲を特定する

`$ARGUMENTS` が指定されていれば base ref として使う。未指定なら main を使う。
`git status -s`、`git log`、`git diff --stat` をそれぞれ個別に実行し、未コミットの変更とコミット済みの差分を把握する。

2. 差分の内容からレビュータイトルとプロンプトを作成する

- **タイトル**: PR タイトル相当の簡潔な説明 (conventional commits 形式)
- **プロンプト**: 以下を含む
  - **コミット範囲**: コミット済みの差分がある場合、BASE_REF と HEAD のハッシュ値を必ず明記する (例: `レビュー対象: abc1234..def5678`)。未コミットの変更がある場合はその旨も記載する
  - 変更概要 (何を、なぜ変更したか)
  - 主な変更点 (箇条書き)
  - 設計判断 (トレードオフや選択の根拠)
  - レビュー観点 (特に見てほしいポイント)

3. codex review を実行する

プロンプト付きで実行すれば、コミット済み・未コミットの変更を問わずレビュー対象になる。

```bash
codex review --title "<タイトル>" "<プロンプト>"
```

### codex review の引数仕様

```
codex review [OPTIONS] [PROMPT]
```

- `--title <TITLE>`: レビュータイトル (必須)
- `--uncommitted`: ステージ済み・未ステージ・未追跡の変更をレビュー (**`[PROMPT]` と併用不可**)
- `--base <BRANCH>`: 指定ブランチとの差分をレビュー (**`[PROMPT]` と併用不可**)
- `--commit <SHA>`: 特定コミットの変更をレビュー
- `[PROMPT]`: レビュー指示 (自由記述)。`-` で stdin から読み込み

**制約**: `--uncommitted`、`--base` はいずれも `[PROMPT]` と同時に使えない。stdin (`-`) もプロンプト扱いで同様に併用不可。

**推奨**: プロンプト付きでレビューしたい場合は `--uncommitted` や `--base` を付けず、プロンプトのみで実行する。codex が自動的に未コミットの変更も含めてレビュー対象を検出する。

## 注意

- プロンプトが長い場合はヒアドキュメントを `$()` で渡す
- codex の実行は時間がかかるため、バックグラウンドで進行する可能性がある
