# X-Ray-Test プロジェクト

X-Rayのテスト用に作成したWebアプリケーションです。

## 主な使用言語

### Backend

- Go
- SQL

### Frontend

- HTML
- CSS
- JavaScript

## 主な使用パッケージ/ライブラリ等

### Backend

- Gin Web Framework
- GORM
- guregu/dynamo
- AWS X-Ray SDK for Go

### Frontend

- 特になし

## 仕様

### 概要

- ユーザーの管理（新規登録、編集、削除、一覧表示）を行うWebアプリケーション
- 登録先のデータベースはMySQL又はDynamoDBを任意選択できる
- X-Rayによるトレースを行う

### メイン画面

- 登録済ユーザーを一覧表示する画面
- メニューからユーザーの新規登録、登録済ユーザーの編集/削除を選択できる

### 登録画面

- ユーザーを新規登録するための画面
- IDを発行し、選択されたDBに登録する

### 編集画面

- 登録済ユーザーの編集を行うための画面
- ID、登録先DBは変更不可

## デプロイ手順

### RDS

- MySQL作成
- インバウンドルールにマイIPを追加
- SSHログインし、データベース/テーブルを作成する ※docs/users.sqlを参照

### DynamoDB

- テーブル作成

### Elastic Beanstalk

- Goのウェブサーバー環境を作成
- RDSのインバウンドルールに本環境を追加
- X-Rayデーモンを有効化
- 以下の環境変数を登録
    - MYSQL_USER：MySQLのユーザー名
    - MYSQL_PASSWORD：MySQLのパスワード
    - MYSQL_HOST：MySQLのホスト名
    - MYSQL_DBNAME：MySQLのデータベース名
    - DYNAMO_REGION：DynamoDBのリージョン
    - DYNAMO_TABLENAME：DynamoDBのテーブル名
    - AWS_XRAY_CONTEXT_MISSING: IGNORE_ERROR
- ビルド
```
cd /path/to/X-Ray-Test
go build -o bin/application application.go
zip ../X-Ray-Test.zip -r *
```
- X-Ray-Test.zipをデプロイ
