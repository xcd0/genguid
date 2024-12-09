# genguid

GUIDを生成する。


```sh
$ ./genguid.exe -h
Usage: genguid [--genversion GENVERSION] [--ulid] [--xid] [--win] [--lower] [--upper] [--output OUTPUT] [--version] [--debug] [NUM [NUM ...]]

Positional arguments:
  NUM                    生成個数。

Options:
  --genversion GENVERSION, -v GENVERSION
                         生成するGUIDのバージョンを指定する。4 or 1
  --ulid                 ULIDを生成する。生成順でソートできるUUIDのようなもの。
  --xid                  xidを生成する。生成順でソートでき、かつURLエンコードが不要なUUIDのようなもの。
  --win, -w              Windowsのレジストリで使用される形式で出力する。{}で囲まれる。
  --lower                小文字で出力する。
  --upper                大文字で出力する。
  --output OUTPUT, -o OUTPUT
                         ファイル出力する。
  --version              バージョン情報を出力する。
  --debug, -d            デバッグ用。ログが詳細になる。
  --help, -h             ヘルプを出力する。

```

## 使い方

単に呼ぶと1つ生成する。

```sh
$ ./genguid.exe
73fb8ff8-7234-4e8b-9076-794fea4d8489
```

数値を与えるとその個数生成する。

```sh
$ ./genguid.exe 10
0b6b6f82-dd3d-41ee-a148-32534ba8b121
1234a65a-8b1b-42aa-ac10-e5f97f136c63
4b323d0e-433d-43b1-994d-fd7ccb108e99
f7a8b58e-2c2a-4d3a-8c0b-f403e211db9f
71934380-86e7-4f70-b6b3-0e0732536f45
30b3006d-c8ff-4ed2-9e7d-392ca1dc4df5
cff90b6d-889a-407c-b8a9-5fb1265d5a60
a10d0071-0a08-4643-b583-795f1d2176ac
3b48c35b-8b10-41a0-ae10-7ded891c1a2c
84dfc24e-4b93-47dc-a53e-4bb46ad6bb65
```

`-w`か`--win`をつけると、windowsのレジストリ形式で出力する。

```sh
$ ./genguid.exe -w
{7472910e-aaa2-48ea-9638-a1f44b18be2e}
```

`--upper`や`--lower`で大文字小文字を制御できる。

```sh
$ ./genguid.exe --upper
B1750630-3E89-4098-866D-F074C8CFD17D
$ ./genguid.exe --lower
1a10fa6e-5b88-45b3-b448-94260e701570
```

その他の機能として
[ULID](https://ja.wikipedia.org/wiki/UUID#ULID) と
[xid](https://github.com/rs/xid) の生成機能がある。

```sh
$ ./genguid.exe --ulid
01JEN25R8TJ9S2KZN7W103WWFJ
$ ./genguid.exe --xid
ctb9058aslujrt6u2c30
```

