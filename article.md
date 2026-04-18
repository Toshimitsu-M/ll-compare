# Go言語もClaude Codeも初めてだったけど、LLM比較アプリを1日で作れた話

## はじめに

最近、社内でもAIツールの話題がよく出るようになりました。ClaudeもChatGPTも名前は聞くけど、正直「何がどう違うの？」がよくわかっていませんでした。

そこで、同じ質問を両方に投げて回答を並べて比較できるアプリを作ってみました。

![スクリーンショット](images/screenshot.png)

## Go言語を選んだ理由（でも触ったことなかった）

バックエンドの言語はGoを使いました。理由は単純で、**社内の話題にも上がっていたから**。

ただ、私自身Goは未経験です。構文もSDKの使い方も全く知らない状態からスタートしました。

## Claude Codeで爆速開発できた

今回の一番の発見は、**Claude Codeを使うと知らない言語でもアプリが作れる**ということです。

「GoでClaude APIとOpenAI APIを並行呼び出しするWebアプリを作りたい」と伝えたら、バックエンドもフロントエンドも一気に書いてくれました。

自分でGoのコードを書いた部分は正直ほぼゼロです。それでも動くものが完成しました。

## 作ったもの

- 1つのプロンプトをClaudeとChatGPTに同時送信
- 左右に並べて回答を表示
- Goのgoroutineで並行処理しているので、2つのAPIを順番に叩くより速い

## セットアップ

APIキーが必要です。

| サービス | 取得先 | 料金 |
|---|---|---|
| Anthropic | console.anthropic.com | 初回$5無料クレジットあり |
| OpenAI | platform.openai.com | 従量課金（事前チャージ必要） |

```bash
git clone https://github.com/Toshimitsu-M/ll-compare.git
cd ll-compare

export ANTHROPIC_API_KEY="your-key"
export OPENAI_API_KEY="your-key"

go run main.go
```

`http://localhost:8080` で起動します。

## 使ってみた感想

同じ質問への回答を並べると、確かに構成や視点が違うのがわかります。ただ、まだ「こういうときはClaudeの方がいい」みたいな使い分けの感覚は自分の中でできていません。

このアプリを使いながら、少しずつ違いを掴んでいけたらと思っています。

## リポジトリ

https://github.com/Toshimitsu-M/ll-compare
