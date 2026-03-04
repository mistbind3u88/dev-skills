---
name: commit
description: 変更をコミットする。変更が大きい場合はレイヤ構成に応じて段階的にコミットし、fixupやamendも適切に使い分ける。
allowed-tools: Bash(git status:*) Bash(git diff:*) Bash(git log:*) Bash(git add:*) Bash(git commit:*) Bash(git push:*) Bash(gh pr view:*) Bash(gh pr create:*) Bash(gh pr edit:*)
---

# commit スキル

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
- lint が通る（プロジェクトに lint がある場合）
- 可能な限りテストが通る

段階的コミットの際は `git add <ファイル>` で対象を選択し、1 段階ずつコミットする。

#### B. fixup（PR 内の既存コミットへの修正・漏れ追加）

PR 内の既存コミットに対する修正やコミット漏れの追加の場合:

```bash
# 修正対象のコミットを特定
git log --oneline main..HEAD

# fixup コミットを作成
git commit --fixup=<対象コミットのSHA>
```

#### C. amend（直前のコミットへの修正）

直前のコミットの変更内容に対する修正の場合:

```bash
git add <修正ファイル>
git commit --amend --no-edit
```

コミットメッセージも変更する場合は `--no-edit` を外す。

### 3. コミットメッセージを作成する

conventional commits 形式で記述する。

```
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

### 5. push 時の PR 管理

ユーザーから push を指示された場合、push 後に PR の状態を確認・管理する。

#### PR が存在しない場合 — 新規作成

```bash
gh pr create --draft --title "<conventional commits 形式のタイトル>" --body "$(cat <<'EOF'
## Summary
<変更の要約を箇条書き>

## Test plan
<テスト方針>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

#### PR が既に存在する場合 — 概要欄を更新

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

## 注意

- `git add -A` や `git add .` は使わない。ファイルを明示的に指定する
- 段階的コミットの各段階で、可能であればコンパイル・lint を実行して壊れていないことを確認する
- fixup コミットの autosquash はユーザーに委ねる（自動で rebase しない）
- amend 後に force push が必要な場合はユーザーに確認する
- push はユーザーが明示的に指示しない限り行わない
