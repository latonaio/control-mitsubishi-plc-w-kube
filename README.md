# control-mitsubishi-plc-w-kube
## 概要
control-mitsubishi-plc-w-kubeは、kanbanから取得したデータを元に、三菱電機製のPLCのレジスタにメッセージを送信するマイクロサービスです。
メッセージの送受信方法およびフォーマットは**MCプロトコル**に準じています。

### MCプロトコルとは
三菱電機製レジスタに採用されている、三菱電機独自のプロトコルです。
16進数のバイナリで構成された電文を送受信し、レジスタに対して操作を行うメッセージングプロトコルです。
[MCプロトコルのマニュアル（三菱電機のHPに遷移します）](https://www.mitsubishielectric.co.jp/fa/download/search.do?mode=keymanual&q=sh080003)


## 動作環境
動作には以下の環境であることを前提とします。   
* OS: Linux   
* CPU: Intel64/AMD64/ARM64   
最低限スペック   
* CPU: 2 core   
* memory: 4 GB   

### 対応している接続方式
* Ethernet接続


## I/O
### Input
kanbanからデータを受信します。
受け取れるkanbanのパラメータは以下の通りです。
```
status: IO(IN: 0, OUT: 1)
```

### Output
kanbanのデータを元に、PLCへデータの書き込みを行います。

## セットアップ
### 電文フォーマット仕様
PLCへの書き込みの仕様は下記の通りです。
* 対応フォーマット：3Eフレーム（固定）
* 接続先ネットワーク：マルチドロップ局（固定）
* 読み取り方式：バイト単位の一括書き込み（固定）
* 自局番号：00（固定）

### デバイス番号
書き込むデバイスのデバイス番号はyamlファイルで設定が可能です。   
yamlファイルは`/var/lib/aion/default/config/`へ設置してください。

### yamlファイルの書き方
```
strContent: デバイス名
iDataSize: データ長
strDevNo: デバイス番号
iReadWrite: IO（IN:0, OUT: 1）
iFlowNo: 実行フロー番号
```

例:
```yaml
settings:
  - strContent: "sample"
    iDataSize: 16
    strDevNo: X9000
    iReadWrite: 0
    iFlowNo: 0
  - strContent: "sample2"
    iDataSize: 16
    strDevNo: X9020
    iReadWrite: 0
    iFlowNo: 0
```

### セットアップ手順
```shell
mv nis_settings.yaml.sample nis_settings.yaml

# nis_settings.yamlを上記のように書き換え

cp nis_settings.yaml /var/lib/aion/default/config/nis_setting.yaml
```

## 関連するマイクロサービス
* [control-mitsubishi-plc-r-kube](https://github.com/latonaio/control-mitsubishi-plc-r-kube)