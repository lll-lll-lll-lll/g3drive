```sh
export GOOGLE_APPLICATION_CREDENTIALS="./credentials.json"
```

# Google Driveへのファイルのアップロードをコマンド一発で管理しようってやつ

#### 概要
- ダウンロードしたファイルの管理が面倒
- わざわざブラウザーを開いてまで、アップロードするのも面倒
- コマンド1つでファイルをアップロードして欲しい。

- 目的は並行処理の学習とアウトプット

#### TODO 
- [x] 並行処理の簡易的
- [ ] 並行処理のエラー処理
- [ ] goroutineのトレース

#### DONE
- 特定のディレクトリ内にあるファイルを全てアップロードする
- 既にアップロードされているファイルはアップロードしない
- ディレクトリがない場合は作成する

サンプルのawsのPDFのリンク
- https://andresriancho.com/wp-content/uploads/2019/06/whitepaper-internet-scale-analysis-of-aws-cognito-security.pdf
