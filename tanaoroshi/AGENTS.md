# tanaoroshi — Issue/PR 棚卸しツール

GitHub CLI（`gh`）を呼び出して Issue/PR のデータ収集・参照解決を行う Go プログラム。

## サブコマンド

| コマンド  | 入力                 | 出力                             |
| --------- | -------------------- | -------------------------------- |
| `collect` | owner/repo（引数）   | リポジトリ別 Issue/PR の JSON    |
| `summary` | collect の JSON      | body を除いたコンパクトな一覧    |
| `refs`    | collect の JSON      | body 内の参照パターン抽出        |
| `resolve` | owner/repo:N（引数） | 参照先の state/title/url を JSON |

## 無視リスト

`skills/tanaoroshi/ignore` に記載された Issue/PR は `summary` と `refs` の出力から除外される。`collect` の生データには引き続き含まれる。

形式: 1行に1つ `owner/repo#N`。`#` で始まる行と空行は無視。

## 既知の制約

### gh コマンド失敗時の挙動

`gh issue list` / `gh pr list` が失敗した場合（認証エラー、レート制限、リポジトリへのアクセス不可など）、該当リポジトリの結果は空配列になりプログラムは正常終了する。stderr に警告は出るが、exit code には反映されない。棚卸し結果が不完全になる可能性がある。

### collect の取得上限

`collect` は `--limit 100` 固定で GitHub CLI を呼び出す。100 件を超える Open Issue/PR があるリポジトリではデータが切り詰められる。超過時は stderr に警告が出る。
