---
name: commit
description: 変更をコミットする。変更が大きい場合はレイヤ構成に応じて段階的にコミットし、fixupやamendも適切に使い分ける。
allowed-tools: Bash(git status:*) Bash(git diff:*) Bash(git log:*) Bash(git add:*) Bash(git commit:*) Bash(git push:*) Bash(gh pr view:*) Bash(gh pr create:*) Bash(gh pr edit:*)
---

# commitスキル

変更内容を適切な粒度でコミットする。

## 手順

### 1. 変更の全体像を把握する

```bash
git status -s
git diff --stat
git diff --staged --stat
git log --oneline -5
```

未コミットの変更内容、ステージ済みの変更、直近のコミット履歴を確認する。

### 2. コミット戦略を決定する

変更内容に応じて以下のいずれかの戦略を選択する。

#### A. 通常のコミット（新規の変更）

変更が小さい、または単一の関心事に収まる場合はそのままコミットする。

変更が大きい場合は、レイヤ構成（domain → infra → usecase → entrypoint、または設定 → ロジック → テスト）に応じて段階的にコミットする。各コミットは以下を満たすこと:

- コンパイルが通る
- lintが通る（プロジェクトにlintがある場合）
- 可能な限りテストが通る

段階的コミットの際は `git add <ファイル>` で対象を選択し、1段階ずつコミットする。

#### B. fixup（PR内の既存コミットへの修正・漏れ追加）

PR内の既存コミットに対する修正やコミット漏れの追加の場合:

```bash
# 修正対象のコミットを特定
git log --oneline main..HEAD

# fixupコミットを作成
git commit --fixup=<対象コミットのSHA>
```

あるコミットに対するfixupがすべて終わったらautosquashする:

1. autosquashを実行する
2. squash後のコミットメッセージが全変更を適切に反映しているか見直し、必要なら `git commit --amend` で修正する
3. lint・build・testを実行して壊れていないことを確認する

```bash
# autosquash実行
GIT_SEQUENCE_EDITOR=: git rebase --autosquash main

# squash後のコミットメッセージを確認・見直し
git log --oneline main..HEAD
git show --stat <squashされたコミット>
# 必要ならgit commit --amendで修正

# 品質チェック
make lint
```

#### C. amend（直前のコミットへの修正）

直前のコミットの変更内容に対する修正の場合:

1. amend後のコミットに含まれる全変更を把握する（直前コミットの内容 + 今回の変更）
2. 全変更を踏まえて、コミットメッセージが適切か判断する
3. メッセージ変更が不要な場合は `--no-edit`、必要な場合は新しいメッセージを指定する

```bash
# amend後の全体像を確認
git log -1 --stat HEAD
git diff --staged --stat

# メッセージ変更不要の場合
git commit --amend --no-edit

# メッセージも更新する場合
git commit --amend -m "$(cat <<'EOF'
<type>: <全変更を反映した説明>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

### 3. コミットメッセージを作成する

conventional commits形式で記述する。

```text
<type>: <簡潔な説明>
```

- `feat`: 新機能
- `fix`: バグ修正
- `refactor`: リファクタリング
- `docs`: ドキュメント
- `test`: テスト
- `ci`: CI/CD
- `chore`: その他

メッセージはヒアドキュメントで渡す:

```bash
git commit -m "$(cat <<'EOF'
<type>: <説明>

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
EOF
)"
```

### 4. コミット後の確認

```bash
git log --oneline -3
git status -s
```

### 5. push時のPR管理

ユーザーからpushを指示された場合、push後にPRの状態を確認・管理する。

#### PRが存在しない場合 — 新規作成

```bash
gh pr create --draft --title "<conventional commits形式のタイトル>" --body "$(cat <<'EOF'
## Summary
<変更の要約を箇条書き>

## Test plan
<テスト方針>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

#### PRがすでに存在する場合 — 概要欄を更新

既存の概要欄を必ず読み込んでから、新しいコミットの内容を反映する。白紙から書き直さない。

```bash
# 既存の概要欄を取得
gh pr view --json body --jq .body

# 更新
gh pr edit --body "$(cat <<'EOF'
<既存の内容をベースに更新>
EOF
)"
```

## rebase時のバックアップ

mainの取り込みなど差分が大きくなるrebaseや、結果の同一性を担保する必要があるrebaseを行う前に、バックアップブランチを作成する。

```bash
# 現在のブランチ名末尾の既存バックアップ接尾辞を置換して作成
git branch -f "$(git rev-parse --abbrev-ref HEAD | sed 's/-[0-9a-f]\{9\}$//')-$(git rev-parse --short=9 HEAD)"
```

fixupのautosquashのみの場合はバックアップ不要。

## 注意

- `git add -A` や `git add .` は使わない。ファイルを明示的に指定する
- 段階的コミットの各段階で、可能であればコンパイル・lintを実行して壊れていないことを確認する
- fixupがすべて終わったらautosquashし、コミットメッセージの見直しと品質チェック（lint・build・test）を行う
- amend後にforce pushが必要な場合はユーザーに確認する
- pushはユーザーが明示的に指示しない限り行わない
