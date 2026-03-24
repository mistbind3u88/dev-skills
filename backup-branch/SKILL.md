---
name: backup-branch
description: 現在の HEAD のバックアップブランチを作成する。rebase や autosquash の前に使う。
allowed-tools: Bash($SKILL_DIR/backup.sh:*)
---

# backup スキル

現在の HEAD を指すバックアップブランチを作成する。

## 実行

```bash
"$SKILL_DIR/backup.sh"
```

作成されたバックアップブランチ名が標準出力に返される。

## 命名規則

`<ブランチ名>-<9桁ハッシュ>`

- ブランチ名の末尾に既にバックアップ接尾辞（ハイフン + 9桁 hex）がある場合は除去してから付け直す
- 同名のバックアップブランチが既にある場合は上書きする（`git branch -f`）
