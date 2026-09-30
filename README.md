# GoShare

GoShare は、Go 言語で作成したシンプルなファイル転送アプリケーションです。

Go のファイル I/O、`io.Reader` / `io.Writer`、TCP 通信、goroutine、独自通信プロトコルなどを、実際にアプリケーションを作りながら学ぶことを目的として開発しています。

現在は TCP を利用したファイル転送と、goroutine を利用した複数接続の同時処理に対応しています。

## 主な機能

- ファイル情報の取得
- `io.Copy` を利用したファイルコピー
- TCP によるメッセージ通信
- TCP によるファイル転送
- goroutine を利用した複数接続の並行処理
- ファイル転送の進捗表示
- ファイル名を含む簡易的な独自通信プロトコル
- 受信したファイルの自動保存

## ディレクトリ構成

```text
goshare/
├── cmd/
│   ├── fileinfo/
│   ├── filecopy/
│   ├── tcphello/
│   └── tcpfile/
│
└── internal/
    ├── file/
    │   ├── info.go
    │   └── copy.go
    │
    ├── progress/
    │   └── reader.go
    │
    └── transfer/
        ├── tcp.go
        └── file.go
```

### `internal/file`

ファイルに関する処理を担当します。

主に以下の処理を実装しています。

- ファイル名・ファイルサイズの取得
- ファイルコピー

### `internal/progress`

`io.Reader` をラップした独自の Reader を実装しています。

読み込んだバイト数を記録することで、ファイル転送の進捗を表示します。

```text
os.File
   ↓
ProgressReader
   ↓
io.Copy
   ↓
TCP Connection
```

元のファイルデータには手を加えず、データが読み込まれる途中で転送済みバイト数を計測しています。

### `internal/transfer`

TCP 通信とファイル転送を担当します。

受信側では、接続ごとに goroutine を起動することで複数の接続を並行して処理できます。

```text
listener.Accept()
       │
       ├── goroutine → Client A
       ├── goroutine → Client B
       └── goroutine → Client C
```

## ファイル転送プロトコル

GoShare では TCP の上に簡単な独自プロトコルを定義しています。

現在は以下の形式でデータを送信します。

```text
┌───────────────────────┐
│ 2 bytes               │ ファイル名の長さ（uint16）
├───────────────────────┤
│ N bytes               │ ファイル名
├───────────────────────┤
│ 残りのTCPストリーム   │ ファイル本体
└───────────────────────┘
```

例えば `photo.jpg` を転送する場合は、概念的には以下のデータが送信されます。

```text
[9][photo.jpg][ファイル本体...]
```

受信側は、

1. ファイル名の長さを取得
2. 指定された長さだけファイル名を読み込む
3. 受信したファイル名でファイルを作成
4. 残りの TCP ストリームをファイルへ書き込む

という順番で処理します。

## ファイル送信の仕組み

送信側では、最初に対象ファイルを開きます。

```text
ファイル
   ↓
os.File
```

次にファイル情報とファイル名を取得し、TCP 接続を確立します。

ファイル本体は `ProgressReader` を経由して送信します。

```text
os.File
   ↓
ProgressReader
   ↓
io.Copy
   ↓
net.Conn
   ↓
TCP
```

`ProgressReader` は、元の `os.File` の `Read()` を利用しながら、読み込んだバイト数を記録します。

これによってファイル転送中の進捗を表示できます。

## ファイル受信の仕組み

受信側では TCP サーバーを起動し、接続を待ち受けます。

```text
net.Listen()
    ↓
Accept()
```

接続を受け付けると、その接続を goroutine に渡します。

```go
go handleConnection(conn, dir)
```

そのため、一つのファイルを受信している途中でも、新しい接続を受け付けることができます。

```text
                ┌→ ファイルA受信
                │
Accept ─────────┼→ ファイルB受信
                │
                └→ ファイルC受信
```

各 goroutine では、

```text
TCP
 ↓
ファイル名を取得
 ↓
保存先ファイルを作成
 ↓
io.Copy
 ↓
ファイル保存
```

という処理を行います。

## このプロジェクトで学んだGoの機能

### `io.Reader` / `io.Writer`

Go ではデータの読み書きを `io.Reader` / `io.Writer` という共通インターフェースで扱えます。

例えば、

```go
io.Copy(dst, src)
```

は `src` が `io.Reader`、`dst` が `io.Writer` を満たしていれば利用できます。

そのため、同じ `io.Copy` で、

```text
File → File

File → TCP

TCP → File
```

といった処理ができます。

### 独自 `io.Reader`

GoShare では進捗表示のために独自の Reader を実装しています。

`io.Reader` は次のインターフェースです。

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

そのため、

```go
func (r *Reader) Read(p []byte) (int, error)
```

を実装することで、独自の型を `io.Reader` として扱えるようになります。

GoShare では、

```text
io.Copy
   ↓
ProgressReader.Read()
   ↓
os.File.Read()
```

という形で元の Reader をラップしています。

### goroutine

Go では、

```go
go function()
```

とすることで処理を並行実行できます。

GoShare の受信側では、

```go
go handleConnection(conn, dir)
```

とすることで、TCP 接続ごとに独立した goroutine を起動しています。

### `defer`

ファイルやネットワーク接続などのリソース解放には `defer` を利用しています。

```go
src, err := os.Open(path)
if err != nil {
	return err
}
defer src.Close()
```

リソースの取得に成功した直後に `defer` を登録することで、関数がどの `return` から終了しても `Close()` が実行されます。

### ポインタ

状態を変更する必要がある構造体ではポインタレシーバを利用しています。

```go
func (r *Reader) Read(p []byte) (int, error)
```

`*Reader` とすることで、コピーではなく同じ `Reader` の状態を変更できます。

例えば、

```go
r.Current += int64(n)
```

によって、読み込むたびに転送済みバイト数を蓄積しています。

## 実行

Go Modules を利用しています。

```bash
go mod init github.com/kuwaharu-git/goshare
```

各学習用コマンドは `cmd` 以下から実行できます。

### ファイル情報取得

```bash
go run ./cmd/fileinfo <ファイル>
```

### ファイルコピー

```bash
go run ./cmd/filecopy <コピー元> <コピー先>
```

### TCP メッセージ通信

受信：

```bash
go run ./cmd/tcphello receive
```

送信：

```bash
go run ./cmd/tcphello send
```

### TCP ファイル転送

受信側を起動した後、別のターミナルから送信します。

```bash
go run ./cmd/tcpfile receive
```

```bash
go run ./cmd/tcpfile send <ファイル>
```

## 今後実装したい機能

- ファイルサイズをプロトコルへ追加
- 受信データサイズの検証
- UDP を利用した LAN 内端末探索
- GoShare が動作している端末の自動検出
- CLI の整理
- 進捗表示の改善
- 転送完了表示
- SHA-256 などを利用したファイル整合性確認
- 複数ファイルの転送
- ディレクトリ転送
- 通信の暗号化

最終的には、

```bash
goshare receive
```

で受信待機し、

```bash
goshare send photo.jpg
```

とするだけで LAN 内の GoShare 端末を自動検出し、ファイルを送信できる形を目指します。

## 開発目的

このプロジェクトでは、単に Go で Web API を作るのではなく、Go が得意とする分野を実際に使いながら学ぶことを目的としています。

主な学習テーマは以下です。

- ファイル I/O
- `io.Reader` / `io.Writer`
- インターフェース
- TCP 通信
- バイナリプロトコル
- goroutine
- ポインタ
- `defer`
- パッケージ設計
- CLI アプリケーション
- ネットワークプログラミング

## License

現在は Go の学習・実験を目的としたプロジェクトです。
