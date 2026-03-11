---
name: fixup
description: 指示された内容に基づいてコードを修正し、対象コミットへの fixup コミットを作成する。
allowed-tools: Bash(git status:*) Bash(git log:*) Bash(git diff:*) Bash(git add:*) Bash(git commit:*) Bash(git show:*) Bash(git rev-parse:*) Read Edit
---

# fixup スキル

指示された修正を行い、適切なコミットへの fixup コミットを作成する。

## 手順

### 1. 修正指示を確認する

`$ARGUMENTS` またはユーザーの口頭指示から、何をどう修正するかを把握する。

### 2. 対象コミットを特定する

修正対象のファイル・内容から、fixup 先となるコミットを特定する。

```bash
git log --oneline main..HEAD
```

- 修正内容が明らかに特定のコミットに属する場合はそのコミットを対象にする
- 対象が不明確な場合はユーザーに確認する

### 3. 修正を実施する

Read でファイルを確認し、Edit で修正する。

### 4. fixup コミットを作成する

```bash
git add <修正したファイル>
git commit --fixup=<対象コミットのSHA>
```

**fixup メッセージの規則**: `fixup!` プレフィックスは常に 1 つだけにする。対象コミットが既に `fixup!` 付きでも積み重ねない。

### 5. 確認

```bash
git log --oneline -3
git status -s
```

## 注意

- `git add -A` や `git add .` は使わない。修正したファイルを明示的に指定する
- `$ARGUMENTS` に `--autosquash` が指定された場合は fixup コミット作成後に autosquash を実行する。指定がない場合は fixup コミットの作成のみで終了する
- autosquash 後はコミットメッセージの見直しと品質チェック（lint・build・test）を行い、成功したら `/mark` でタグを設置する
