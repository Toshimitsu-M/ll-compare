# ClaudeとChatGPTを同時に比較できるWebアプリをGoで作った

## はじめに

同じプロンプトをClaudeとChatGPTに投げて、回答を並べて比較したいと思ったことはないでしょうか。  
「どっちの説明がわかりやすいか」「どっちのコードが正確か」をいちいちブラウザを切り替えながら確認するのは面倒です。

そこで、1つのプロンプトを両モデルに同時送信して、回答を左右に並べて表示するWebアプリを作りました。

## 作ったもの

![スクリーンショット](images/screenshot.png)

- 左：Claude（Opus 4.7）
- 右：ChatGPT（GPT-4o）
- プロンプトを1回入力するだけで両方の回答が同時に表示される

## 技術スタック

- **バックエンド**：Go
- **フロントエンド**：HTML / CSS / JavaScript（フレームワークなし）
- **API**：Anthropic API / OpenAI API

## 仕組み

```
ブラウザ → Goサーバー → Claude API  ┐
                    → OpenAI API ┘ （並行処理）→ 両回答をまとめてレスポンス
```

GoのgoroutineでClaudeとGPTへのAPIリクエストを並行実行しているため、  
2つのAPIを順番に叩くより応答が速いです。

## 主要なコード

### APIの並行呼び出し（main.go）

```go
var wg sync.WaitGroup
wg.Add(2)

go func() {
    defer wg.Done()
    claudeResp, _ = callClaude(req.Message)
}()

go func() {
    defer wg.Done()
    gptResp, _ = callGPT(req.Message)
}()

wg.Wait()
```

`sync.WaitGroup` で2つのgoroutineが両方終わるまで待ってからレスポンスを返しています。

## セットアップ

### 1. APIキーの取得

| サービス | 取得先 | 料金 |
|---|---|---|
| Anthropic | console.anthropic.com | 初回$5無料クレジットあり |
| OpenAI | platform.openai.com | 従量課金（事前チャージ必要） |

### 2. 起動

```bash
git clone https://github.com/Toshimitsu-M/ll-compare.git
cd ll-compare

export ANTHROPIC_API_KEY="your-key"
export OPENAI_API_KEY="your-key"

go run main.go
```

ブラウザで `http://localhost:8080` を開いて完了です。

## 使ってみた感想

同じ質問への回答を並べると、モデルによって構成や視点が異なることがよくわかります。  
例えばGo言語について聞くと、ClaudeはコードサンプルをすぐCode blockとして見せてくれる一方、GPT-4oは箇条書きで特徴を整理してから説明する傾向がありました。

どちらが優れているというより、**用途によって使い分けるための判断材料**としてこのアプリが役立ちそうです。

## リポジトリ

https://github.com/Toshimitsu-M/ll-compare
