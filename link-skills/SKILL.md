---
name: link-skills
description: リポジトリのトップレベルに skills ディレクトリがある場合、.claude/skills にシンボリックリンクを作成する。
allowed-tools: Bash(ls:*) Bash(mkdir:*) Bash(ln:*) Bash(readlink:*)
---

# link-skills

リポジトリルートの `skills` ディレクトリを `.claude/skills` にシンボリックリンクする。

## 手順

1. リポジトリルートに `skills` ディレクトリが存在するか確認する

```bash
ls -d skills 2>/dev/null
```

存在しない場合はエラーメッセージを出して終了する。

2. `.claude/skills` の現在の状態を確認する

```bash
ls -la .claude/skills 2>/dev/null
readlink .claude/skills 2>/dev/null
```

- 既に `../skills` へのシンボリックリンクが存在する場合は「既にリンク済み」と報告して終了する
- `.claude/skills` がシンボリックリンクでないディレクトリとして存在する場合は、上書きせず警告を出して終了する

3. `.claude` ディレクトリを作成する（存在しない場合）

```bash
mkdir -p .claude
```

4. シンボリックリンクを作成する

```bash
ln -s ../skills .claude/skills
```

5. 結果を確認して報告する

```bash
ls -la .claude/skills
```
