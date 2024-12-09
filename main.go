package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/google/uuid"
	"github.com/oklog/ulid"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

var (
	version  string = "debug build" // makefileからビルドされると上書きされる。
	revision string
	Revision = func() string { // {{{
		revision := ""
		modified := false
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					//return setting.Value
					revision = setting.Value
					if len(setting.Value) > 7 {
						revision = setting.Value[:7] // 最初の7文字にする
					}
				}
				if setting.Key == "vcs.modified" {
					modified = setting.Value == "true"
				}
			}
		}
		if modified {
			revision = "develop+" + revision
		}
		return revision
	}() // }}}
	parser *arg.Parser // ShowHelp() で使う

	toJsonFromXml = true // jsonからxmlに変換する。
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile) // ログの出力書式を設定する

	args := ArgParse()

	num := 1
	if len(args.Num) > 0 {
		num = args.Num[0]
	}
	if args.Debug {
		log.Printf("num: %v", num)
	}
	str := ""
	for i := 0; i < num; i++ {
		id := ""
		if args.Ulid {
			id = GenULID()
		} else if args.Xid {
			id = GenXid()
		} else if args.GenVersion == 1 {
			id = GenGUIDv1()
		} else {
			id = GenGUIDv4()
		}
		if args.Win {
			str += fmt.Sprintf("{%s}\n", id)
		} else {
			str += fmt.Sprintf("%s\n", id)
		}
	}

	if args.ToUpper {
		str = strings.ToUpper(str)
	} else if args.ToLower {
		str = strings.ToLower(str)
	}

	fmt.Printf("%s\n", str)
	if len(args.Output) > 0 {
		WriteText(args.Output, str)
	}
}

func GenGUIDv4() string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", uuid)
}

func GenGUIDv1() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", uuid)
}

func GenULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return fmt.Sprintf("%s", id)
}

func GenXid() string {
	guid := xid.New()
	return guid.String()
}

func WriteText(file, str string) {
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		panic(errors.Errorf("%v", err))
	} else {
		if _, err := f.Write([]byte(str)); err != nil {
			panic(errors.Errorf("%v", err))
		}
	}
}

type Args struct {
	Num        []int  `arg:"positional"   help:"生成個数。" default: [1]`
	GenVersion int    `arg:"-v"           help:"生成するGUIDのバージョンを指定する。4 or 1" default:4`
	Ulid       bool   `arg:"--ulid"       help:"ULIDを生成する。生成順でソートできるUUIDのようなもの。"`
	Xid        bool   `arg:"--xid"        help:"xidを生成する。生成順でソートでき、かつURLエンコードが不要なUUIDのようなもの。"`
	Win        bool   `arg:"-w,--win"     help:"Windowsのレジストリで使用される形式で出力する。{}で囲まれる。" default:false`
	ToLower    bool   `arg:"--lower"      help:"小文字で出力する。" default:false`
	ToUpper    bool   `arg:"--upper"      help:"大文字で出力する。" default:false`
	Output     string `arg:"-o,--output"  help:"ファイル出力する。"`
	Version    bool   `arg:"--version" help:"バージョン情報を出力する。"`
	Debug      bool   `arg:"-d,--debug"      help:"デバッグ用。ログが詳細になる。"`
}
type ArgsVersion struct {
}

func (args *Args) Print() {
	log.Printf(`
	Num         : %v
	GenVersion  : %v
	Ulid        : %v
	Xid         : %v
	ToLower     : %v
	ToUpper     : %v
	Output      : %v
	Win         : %v
	Version     : %v
	`,
		args.Num,        //[]int `arg:"positional"   help:"生成個数。未指定の時1つ生成する。"`
		args.GenVersion, //int   `arg:"-v"           help:"生成するGUIDのバージョンを指定する。未指定の時4。4 or 1"`
		args.Ulid,       //bool  `arg:"-l"           help:"ULIDを生成する。生成順でソートできるUUIDのようなもの。"`
		args.Xid,        //bool  `arg:"-x"           help:"xidを生成する。生成順でソートでき、かつURLエンコードが不要なUUIDのようなもの。"`
		args.ToLower,    //bool   `arg:"--lower"      help:"小文字で出力する。" default:false`
		args.ToUpper,    //bool   `arg:"--upper"      help:"大文字で出力する。" default:false`
		args.Output,     //string `arg:"-o"           help:"ファイル出力する。"`
		args.Win,        //bool  `arg:"-w"           help:"Windowsのレジストリで使用される形式で出力する。{}で囲まれる。" default:false`
		args.Version,    //bool  `arg:"-V,--version" help:"バージョン情報を出力する。"`
	)
}

func ShowHelp(post string) {
	buf := new(bytes.Buffer)
	parser.WriteHelp(buf)
	fmt.Printf("%v\n", strings.ReplaceAll(buf.String(), "display this help and exit", "ヘルプを出力する。"))
	if len(post) != 0 {
		fmt.Println(post)
	}
	os.Exit(1)
}
func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
func ShowVersion() {
	if len(Revision) == 0 {
		// go installでビルドされた場合、gitの情報がなくなる。その場合v0.0.0.のように末尾に.がついてしまうのを避ける。
		fmt.Printf("%v version %v\n", GetFileNameWithoutExt(os.Args[0]), version)
	} else {
		fmt.Printf("%v version %v.%v\n", GetFileNameWithoutExt(os.Args[0]), version, revision)
	}
	os.Exit(0)
}

func ArgParse() *Args {
	args := &Args{
		//Output: "./generated_guid.txt", //      string `arg:"-o"           help:"ファイル出力する。"`
	}

	var err error
	parser, err = arg.NewParser(arg.Config{Program: GetFileNameWithoutExt(os.Args[0]), IgnoreEnv: false}, args)
	if err != nil {
		ShowHelp(fmt.Sprintf("%v", errors.Errorf("%v", err)))
	}

	if err := parser.Parse(os.Args[1:]); err != nil {
		if err.Error() == "help requested by user" {
			ShowHelp("")
		} else if err.Error() == "version requested by user" {
			ShowVersion()
		} else {
			panic(errors.Errorf("%v", err))
		}
	}

	if args.Version {
		ShowVersion()
		os.Exit(0)
	}

	if args.Debug {
		args.Print()
	}
	return args
}
