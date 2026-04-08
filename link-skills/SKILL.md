---
name: link-skills
description: リポジトリ内のスキルを Codex または Claude Code から使えるようにリンクする。Codex では ~/.codex/skills に各スキルディレクトリへのジャンクションを作成し、Claude Code では .claude/skills へのリンクを作成する。
allowed-tools: Bash(ls:*) Bash(mkdir:*) Bash(ln:*) Bash(readlink:*)
---

# link-skills

Codex または Claude Code からこのリポジトリのスキルを使えるようにリンクする。

## 手順

1. 対象エージェントと OS を確認する

- Codex on Windows の場合: `~/.codex/skills` 配下に各スキルディレクトリへのジャンクションを作成する
- Claude Code on macOS / Linux の場合: `.claude/skills` へのシンボリックリンクを作成する

2. リポジトリ内のスキルディレクトリを確認する

```bash
find . -name SKILL.md
```

Codex では各 `SKILL.md` の親ディレクトリをリンク対象にする。`.skill/daily-tagging` のようにルート直下でないスキルも見落とさない。

3. Codex on Windows の場合は `~/.codex/skills` の現在の状態を確認する

```bash
ls ~/.codex/skills
```

- 同名エントリが既にある場合はリンク先を確認する
- 想定外の既存ディレクトリやファイルがある場合は上書きしない

4. Codex on Windows の場合は各スキルディレクトリへのジャンクションを作成する

```bash
mklink /J %USERPROFILE%\.codex\skills\<skill-name> C:\path\to\dev-skills\<skill-name>
```

- 既に正しいジャンクションがある場合はそのままにする

5. Claude Code on macOS / Linux の場合は `.claude` ディレクトリを作成する（存在しない場合）

```bash
mkdir -p .claude
```

6. Claude Code on macOS / Linux の場合は `.claude/skills` の現在の状態を確認する

```bash
ls -la .claude/skills 2>/dev/null
readlink .claude/skills 2>/dev/null
```

- すでに正しいリンクが存在する場合はそのまま終了する
- `.claude/skills` がシンボリックリンクでないディレクトリとして存在する場合は、上書きせず警告を出して終了する

7. Claude Code on macOS / Linux の場合は `.claude/skills` へのリンクを作成する

```bash
ln -s ../skills .claude/skills
```

8. 結果を確認して報告する

Codex on Windows:

```bash
dir %USERPROFILE%\.codex\skills
```

Claude Code on macOS / Linux:

```bash
ls -la .claude/skills
```
