# dev-skills

Claude Code 向けの汎用開発ワークフロースキル集。

## Skills

| スキル                                  | 概要                                         |
| --------------------------------------- | -------------------------------------------- |
| [clean-docs](./clean-docs/SKILL.md)     | `.claude/docs` のタスクドキュメント整理      |
| [codex-review](./codex-review/SKILL.md) | codex CLI によるコードレビュー               |
| [commit](./commit/skill.md)             | git コミット（段階的コミット、fixup、amend） |
| [gh-edit](./gh-edit/skill.md)           | GitHub PR/Issue の作成・更新                 |
| [link-skills](./link-skills/SKILL.md)   | skills ディレクトリのシンボリックリンク作成  |

## Setup

`~/.claude/skills` にこのリポジトリへのシンボリックリンクを作成する。

```bash
ln -s /path/to/dev-skills ~/.claude/skills
```
