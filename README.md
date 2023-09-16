# Odori

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)

パラパラ漫画のようなシンプルなアニメーションを作成するための[Ebitengine](https://ebitengine.org/)製GUIツール。

自分用かつエイヤで作ったものなのでもろもろ雑ですが、機能としては以下のようなものがあります。

- PNG画像の読み込み
- スプライトシートの読み込み
- Import/Export機能(JSONとスプライトシート)
- GIF出力機能

スプライトシートの読み書きには拙作の[Kaban](https://github.com/aethiopicuschan/kaban)を内部的に利用しています。

## インストール

```sh
go install github.com/aethiopicuschan/odori@latest
odori
```

## 動作環境

Macでのみ動作確認しています。Windowsや各種Linuxでも動くとは思いますが、想定外の動作などをするかもしれません。
