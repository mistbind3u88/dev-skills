---
name: mark
description: チェック通過を示すローカル軽量タグを現在の HEAD に設置・確認・削除する。
allowed-tools: Bash($SKILL_DIR/mark.sh:*)
---

# mark スキル

チェック通過状態をローカル軽量タグで管理する。タグは `git push` では送信されない。

## タグ命名規則

```
check/<type>
```

| タグ名         | 意味                  |
| -------------- | --------------------- |
| `check/lint`   | lint 通過済み         |
| `check/test`   | テスト通過済み        |
| `check/build`  | ビルド通過済み        |
| `check/review` | codex review 実施済み |

## 手順

`$SKILL_DIR/mark.sh` を使って操作する。

### タグを設置する

```bash
"$SKILL_DIR/mark.sh" <type>
```

### 状態を確認する

```bash
"$SKILL_DIR/mark.sh" --status
```

### 全タグを削除する

```bash
"$SKILL_DIR/mark.sh" --clean
```

## 注意

- タグはローカル専用。`git push` のデフォルトでは送信されない
- コミットが進むとタグは古いコミットに残るため、再チェック後に再設置が必要
- 以下のスキルが関連処理の成功後に `/mark` を呼び出してタグを設置する:
  - **codex-review**: レビュー完了後に `check/review`
  - **push**: build/lint/test 実行成功後に `check/build`、`check/lint`、`check/test`
  - **commit**: fixup 後の品質チェック成功後に `check/lint`、`check/build`、`check/test`
