---
name: commit
description: 変更をコミットする。変更が大きい場合はレイヤ構成に応じて段階的にコミットし、fixupやamendも適切に使い分ける。
allowed-tools: Bash(git status:*) Bash(git diff:*) Bash(git log:*) Bash(git add:*) Bash(git commit:*) Bash(git show:*) Bash(git rev-parse:*) Bash(git stash:*) Bash(git restore:*) Bash(git branch:*) Bash(GIT_SEQUENCE_EDITOR=:*)
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

段階的コミットの際は `git add <ファイル>` で対象を選択し、1 段階ずつコミットする。各コミットの検証手順:

1. `git stash --keep-index --include-untracked` で未ステージの変更を退避
2. `/check --skip-review` で lint・build・test を実行（成功時に `/mark` でタグ設置される）
3. コミット
4. `git stash pop` で退避した変更を復元（コンフリクトした場合は `git restore` で HEAD に戻し stash を drop）

#### B. fixup（PR 内の既存コミットへの修正・漏れ追加）

`/fixup` スキルに委ねる。

#### C. amend（直前のコミットへの修正）

直前のコミットの変更内容に対する修正の場合:

1. amend 後のコミットに含まれる全変更を把握する（直前コミットの内容 + 今回の変更）
2. 全変更を踏まえて、コミットメッセージが適切か判断する
3. メッセージ変更が不要な場合は `--no-edit`、必要な場合は新しいメッセージを指定する

```bash
# amend 後の全体像を確認
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

conventional commits 形式で記述する。

```
<type>: <簡潔な説明>
```

- `feat`: 新機能（外部から見た振る舞いが変わる変更）
- `fix`: バグ修正
- `refactor`: リファクタリング（外部から見た振る舞いが変わらない内部変更）
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

## autosquash

- main の取り込み（rebase）は `/catch-up` スキルに委ねる。autosquash 時に main を取り込まない
- autosquash の起点には main からブランチを切ったコミットハッシュを指定する

1. `/backup-branch` でバックアップブランチを作成する
2. autosquash を実行する

```bash
BASE=$(git merge-base main HEAD)
GIT_SEQUENCE_EDITOR=: git rebase --autosquash --rebase-merges "$BASE"
```

- コンフリクト解消時は、あるコミットに対する全ての fixup が squash された段階（= そのコミットが完成した状態）で lint・build・test を実行して通過を確認する。途中の fixup 適用中はビルドが通らない場合があるため、同一コミットへの fixup が連続する間はスキップしてよい

## 注意

- main/master ブランチ上で直接コミットしない。コミット前に現在のブランチを確認し、main/master であれば作業ブランチを切ってから作業する
- `git add -A` や `git add .` は使わない。ファイルを明示的に指定する
- 段階的コミットの品質検証は手順2-Aに従う。`/check --skip-review` が lint・build・test の実行と `/mark` でのタグ設置を行う
- amend 後に force push が必要な場合はユーザーに確認する
- push はユーザーが明示的に指示しない限り行わない
- コミット後の PR 作成・更新は gh-edit スキル、push は push スキルに委ねる
