# OASC (OpenAPI Specification Combiner)

OpenAPI仕様書をマージするCLIツールです。

## 特徴

- OpenAPI 3.x仕様書のマージをサポート
- パス、オペレーション、パラメータ、レスポンス、リクエストボディの適切なマージ
- 重複するパラメータやタグの自動処理
- バージョン互換性の検証

## インストール

```bash
go install github.com/haryoiro/oasc@latest
```

## 使用方法

```bash
# 基本的な使用方法
oasc -file spec1.yaml -file spec2.yaml -output merged.yaml

# 複数のファイルをマージ
oasc -file spec1.yaml -file spec2.yaml -file spec3.yaml -output merged.yaml

# JSON形式で出力
oasc -file spec1.yaml -file spec2.yaml -output merged.json -format json
```

### オプション

- `-file`, `-f`: 入力ファイル（複数指定可能）
- `-output`, `-o`: 出力ファイルパス（デフォルト: merged.yaml）
- `-format`, `-F`: 出力形式（json または yaml）

## マージルール

- パス: 同じパスが存在する場合、HTTPメソッドごとにマージ
- オペレーション: パラメータ、レスポンス、リクエストボディを適切にマージ
- パラメータ: 重複するパラメータ名は1つに統合
- コンポーネント: スキーマ、レスポンス、パラメータなどを適切にマージ
- タグ: 重複するタグは1つに統合
