---
name: mark
description: チェック通過を示すローカル軽量タグを現在の HEAD に設置・確認・削除する。
allowed-tools: Bash(git tag:*) Bash(git rev-parse:*)
---

# mark スキル

チェック通過状態をローカル軽量タグで管理する。タグは `git push` では送信されない。

## タグ命名規則

```
check/<type>
```

| タグ名           | 意味                       |
| ---------------- | -------------------------- |
| `check/lint`     | lint 通過済み              |
| `check/test`     | テスト通過済み             |
| `check/build`    | ビルド通過済み             |
| `check/review`   | codex review 実施済み      |

## 手順

### 1. `$ARGUMENTS` から操作を判定する

| 引数パターン              | 操作                                         |
| ------------------------- | -------------------------------------------- |
| `lint`、`test` 等のタイプ | 指定タイプのタグを現在の HEAD に設置          |
| `--status`                | 現在の HEAD に設置されているタグ一覧を表示    |
| `--clean`                 | `check/` プレフィックスのタグを全て削除       |

### 2. タグを設置する

```bash
git tag -f "check/<type>" HEAD
```

`-f` で既存タグがあれば現在の HEAD に移動する。

### 3. 設置結果を報告する

```bash
# 現在の HEAD に付いている check/ タグを一覧表示
git tag --points-at HEAD | grep '^check/'
```

## `--status` の出力例

```
check/build   ✓ (現在の HEAD)
check/lint    ✓ (現在の HEAD)
check/test    ✗ (abc1234 — 2 commits behind)
check/review  ✗ (未設置)
```

各タグについて、現在の HEAD を指しているか、別のコミットを指しているか、未設置かを表示する。

## `--clean` の動作

```bash
git tag -l 'check/*' | xargs -r git tag -d
```

## 注意

- タグはローカル専用。`git push` のデフォルトでは送信されない
- コミットが進むとタグは古いコミットに残るため、再チェック後に再設置が必要
- 他のスキル（push、codex-review）から自動的に設置される場合がある
