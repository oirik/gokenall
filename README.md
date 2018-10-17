# gokenall

[![Build Status](https://travis-ci.org/oirik/gokenall.svg?branch=master)](https://travis-ci.org/oirik/gokenall)
[![GoDoc](https://godoc.org/github.com/oirik/gokenall?status.svg)](https://godoc.org/github.com/oirik/gokenall)
[![apache license](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE)

日本郵便の郵便番号データ（ken_all.csv）をGoから扱うためのライブラリ及びコマンドラインツールです。


以下のような機能があります。

* 最新のken_all.csvを日本郵便のサイトからダウンロード・解凍する。（コマンド名: Download）
* 前回ダウンロード時から更新があるか確認する。（コマンド名: Updated)
* データの使いづらい部分を加工する。（コマンド名: Normalize）
    * sjis→utf8
    * 半角カナ→全角カナ。ASCII文字→半角
    * 複数行に分割された行をマージ
    * 地名項目の運用上邪魔になる文字を修正
        * 「以下に掲載がない場合」は除去
        * 「～の次に番地がくる場合」は除去
        * 「～一円」は除去（但し地名が「一円」の場合は除く）
        * （）で囲まれている部分は除去。但し、可能なものは残すようにする。
            *  (除去) その他
            *  (除去) 地階・階層不明
            *  (除去) *を除く
            *  (残す) ○階
            *  (分割) ○～○(丁目|番地|番)
            *  (分割) ○、○、○(丁目|番地|番)
            *  (分割) 地名、地名、地名

# Usage

## ライブラリ利用

```go
import "github.com/oirik/gokenall"
```

[GoDoc](https://godoc.org/github.com/oirik/gokenall)

## コマンド利用

### install

```sh
$ go get github.com/oirik/gokenall/cmd/kenall
```

Or download binaries from [github releases](https://github.com/oirik/gokenall/releases)


### usage

```
$ kenall <command> [arguments]
```

（例）日本郵便のデータが更新された時だけダウンロードする。

1. `UPDATED`という名前でファイルを作り、中身を以下のようにします。
```text:UPDATED
00010101
```

2. 下記のようなコマンドを叩けば、更新されたときだけカレントディレクトリのken_all.csvが上書きされます。
```sh
$ kenall updated -p `cat UPDATED` > UPDATED && kenall download -x | kenall normalize -o ken_all.csv
```

詳しくはヘルプを参考にしてください。

```
$ kenall help
kenall is a tool for managing ken_all.csv

Usage:

  kenall <command> [arguments]

The commands are:

  download   Download ken_all.zip from japanpost website
  help       Show help information
  normalize  Normalize -make easy to use- input (file or standard input if no argument)
  updated    Read updated date of data from japanpost website. Exit status 0 if later than [argument](yyyyMMdd) or exit status 1.
  version    Show version information

Use "kenall help <command>" for more information about a command.
```

# Todo

* [ ] 単一JSONファイルへの変換
* [ ] 郵便番号ごとのJSONファイルへの変換
* [ ] READMEをもっと親切にする。
