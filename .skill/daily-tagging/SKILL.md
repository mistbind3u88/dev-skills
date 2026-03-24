---
name: daily-tagging
description: リポジトリを pull し、未タグの日付ごとにその日の最後のコミットへ daily-YYYY-MM-DD タグを作成して push する。
allowed-tools: Bash($SKILL_DIR/scripts/tag_daily_last_commits.sh:*)
---

# daily-tagging

各日付の「最後のコミット」に `daily-YYYY-MM-DD` 形式の軽量タグを付与し、リモートへ push する。

## 実行手順

次のコマンドを実行する。

```bash
"$SKILL_DIR/scripts/tag_daily_last_commits.sh"
```

## 規則

- タグ形式は `daily-YYYY-MM-DD` を使う
- 対象は現在のブランチ履歴上で、まだ `daily-*` タグがない日付のみ
- 各日付では、その日の最後のコミット（時刻が最も遅いコミット）を選ぶ
- 既存の同名タグがその日の最新コミットを指していない場合は付け替える

## 出力

- 作成したタグ名と対象コミットを一覧表示する
- 作成または付け替えがあったタグのみリモートへ push する（付け替え時は force push）
- 作成対象がない場合は何も push しない
