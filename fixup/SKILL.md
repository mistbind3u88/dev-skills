---
name: fixup
description: コードを修正し対象コミットへの fixup コミットを作成する。レビュー指摘の対応やコミットの修正漏れの追加に使う。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git add:*) Bash(git commit:*) Bash(git show:*) Bash(git rev-parse:*) Read Edit
---

# fixup スキル

指示された修正を行い、適切なコミットへの fixup コミットを作成する。

## フラグ

| フラグ     | 動作                       | ファイル編集       | ステージング     | 修正指示 |
| ---------- | -------------------------- | ------------------ | ---------------- | -------- |
| (なし)     | 現行どおり                 | する (Read + Edit) | する (`git add`) | 必要     |
| `--staged` | ステージ済み変更から fixup | しない             | しない（済み）   | 不要     |

## 手順

### 1. フラグと修正指示を確認する

`$ARGUMENTS` から `--staged` の有無を確認する。

- `--staged` なし: 修正指示を `$ARGUMENTS` またはユーザーの口頭指示から把握する
- `--staged` あり: 修正指示があっても無視し、ステージ済み変更をそのまま使う旨を警告する

### 2. ステージ済み変更を確認する（`--staged` 時のみ）

```bash
git diff --staged --name-only
```

- ステージ済み変更がない場合はエラーで終了する
- 変更があるファイル一覧を把握する

### 3. 対象コミットを特定する

```bash
git log --oneline main..HEAD
```

- `--staged` なし: 修正内容から fixup 先のコミットを推定する
- `--staged` あり: `git diff --staged --name-only` のファイルと各コミットの変更ファイルを照合して推定する
- 修正内容が明らかに特定のコミットに属する場合はそのコミットを対象にする
- 対象が不明確な場合はユーザーに確認する

### 4. 修正を実施する（`--staged` 時はスキップ）

Read でファイルを確認し、Edit で修正する。

### 5. fixup コミットを作成する

- `--staged` なし:

```bash
git add <修正したファイル>
git commit --fixup=<対象コミットのSHA>
```

- `--staged` あり（`git add` をスキップ）:

```bash
git commit --fixup=<対象コミットのSHA>
```

**fixup メッセージの規則**: `fixup!` プレフィックスは常に 1 つだけにする。対象コミットが既に `fixup!` 付きでも積み重ねない。

### 6. 確認

```bash
git log --oneline -3
git status -s
```

## 注意

- `git add -A` や `git add .` は使わない。修正したファイルを明示的に指定する
- `$ARGUMENTS` に `--autosquash` が指定された場合は fixup コミット作成後に autosquash を実行する。指定がない場合は fixup コミットの作成のみで終了する
- autosquash 後はコミットメッセージの見直しと品質チェック（lint・build・test）を行い、成功したらスキル `/mark` を実行してタグを設置する
