package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Reader interface {
	Read(rc chan LineMessage, wg *sync.WaitGroup)
}

type Writer interface {
	Write()
}

type ReadFormLogFile struct {
	path          string ""
	TotalLineCont int
	decodeErrCont int
}

var wg sync.WaitGroup

func (r *ReadFormLogFile) Read(rc chan LineMessage, wg *sync.WaitGroup) {

	defer wg.Done()

	fmt.Println(r.path)
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}
	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\s-\s([^\s]*|\-)\s(\[[^\[\]]+\])\s([^\s]*)\s(\"(?:[^"]|\")+|-\")\s([^\s]*)\s(\"\d{3}\")\s(\d+|-)\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")`)
	rd := bufio.NewReader(f)

	defer f.Close()

	//var decodeErrCont = 0
	//var totalLineCont = 0
	for {
		line, err := rd.ReadBytes('\n')
		var lineMessage LineMessage

		if err == io.EOF {
			fmt.Println("最后一行！")
			break

		} else if err != nil {
			fmt.Println(err.Error())
			break
		}
		r.TotalLineCont++

		result := re.FindStringSubmatch(string(line))

		if len(result) != 12 {
			//glog.Info("FindStringSubmatch error:", string(line))
			//logs.Error("FindStringSubmatch error:", string(line))
			continue
		}

		if strings.Contains(result[5], "POST") {
			//fmt.Println(string(line))
			m, _ := url.ParseQuery(result[6])
			if m.Get("data") == "" {
				//fmt.Println("error:", result[5])

				aslice := strings.Split(result[5], " ")
				if len(aslice) >= 3 {
					t, _ := url.ParseQuery(aslice[1])
					if t.Get("data") != "" {
						//fmt.Println("error:", t.Get("data"))
						//hasKeyword(t.Get("data"), string(line))

						lineMessage.data = t.Get("data")
						lineMessage.line = string(line)
						rc <- lineMessage

						continue
					}
				}

				//fmt.Println("error:", result[6])
				// get 参数
			} else {
				if strings.Contains(m.Get("data"), "{") {

					lineMessage.data = m.Get("data")
					lineMessage.line = string(line)
					rc <- lineMessage
					continue

				} else {

					encrypt := strings.ReplaceAll(m.Get("data"), "x0A", "---")
					encrypt = strings.ReplaceAll(encrypt, `\`, "-")
					encrypt = strings.ReplaceAll(encrypt, `----`, "")

					var missing = (4 - len(encrypt)%4) % 4
					encrypt += strings.Repeat("=", missing)
					//res, err := base64.URLEncoding.DecodeString(data)
					//fmt.Println("  decodebase64urlsafe is :", string(res), err)
					dst, _ := base64.URLEncoding.DecodeString(encrypt)

					//dst,_ := Base64URLDecode(encrypt)
					//fmt.Println("Base64URLDecode:", string(dat))
					//
					//encrypt = strings.ReplaceAll(encrypt, `===`, "+++")
					//dst, err := base64.StdEncoding.DecodeString(encrypt) // 490  v 5000
					//

					if err != nil {
						r.decodeErrCont++
						fmt.Println("替换后:", encrypt)
						fmt.Println("解码败:", m.Get("data"))
						fmt.Println(string(line))
						//logs.Error("base64Decode error:", string(line))
						continue
						//fmt.Println("error:", err)
					}
					//hasKeyword(string(dst), string(line))

					lineMessage.data = string(dst)
					lineMessage.line = string(line)
					rc <- lineMessage
					continue
					//fmt.Println("解码后的数据为:", string(dst))
				}

				// body 参数
			}
		} else if strings.Contains(result[5], "GET") {

			m, _ := url.ParseQuery(result[5])
			if m.Get("data") == "" {

				aslice := strings.Split(result[5], "data=")
				if len(aslice) >= 2 {
					arslice := strings.Split(aslice[1], " ")

					if len(arslice) >= 2 {
						getParam := strings.ReplaceAll(arslice[0], "\\x22", "\"")
						getParam = strings.ReplaceAll(getParam, "\\x5C", "\\")
						//fmt.Println("解码后:", getParam)
						//hasKeyword(getParam, string(line))

						lineMessage.data = getParam
						lineMessage.line = string(line)
						rc <- lineMessage
						continue
					} else {
						// error
						//logs.Error("params format error:", string(line))
					}
					//Str ,_:= hex.DecodeString(arslice[0])
					//fmt.Println(result[5])
					//fmt.Println("解码前：",arslice[0])
					//fmt.Println("解码后：",Str)
				} else {
					// error
					//logs.Error("params format error:", string(line))
				}
				//Str ,_:= hex.DecodeString(result[5])
				//fmt.Println(result[5])
				//fmt.Println("解码后：",Str)
			} else {
				//hasKeyword(m.Get("data"), string(line))
				lineMessage.data = m.Get("data")
				lineMessage.line = string(line)
				rc <- lineMessage
				continue
			}
			//fmt.Println(m.Get("data"))
		}
	}

	fmt.Println("读取结束。")
	var endLineMessage LineMessage
	endLineMessage.data = "stop"
	endLineMessage.line = "stop"
	rc <- endLineMessage

	//line := "messages"
	//rc <- line
	//rc <- "stop"
}

type LineMessage struct {
	data string
	line string
	idx  int64
}

type LogProcess struct {
	rc    chan LineMessage
	read  Reader
	write Writer
	hit   int
}

func (l *LogProcess) Process(wg *sync.WaitGroup) {
	//time.Sleep(2 * time.Second)
	defer wg.Done()
	for v := range l.rc {
		if v.data == "stop" {
			// 需要停掉协程
			break
		}
		var has = true
		for _, word := range Config.keywords {
			if strings.Contains(v.data, word) == false {
				has = false
				break
			}
		}
		if has {
			l.hit++
			fmt.Print(v.line)
			if Config.logDecode {
				fmt.Println(v.data, "\n")
			}
		}
	}
}

//https://github.com/steveoc64/jass/blob/47a385d110c0ee17a0781006d3cb9913279df199/server/config.go
type ConfigType struct {
	path      string
	keyword   string
	keywords  []string
	regexp    string
	logDecode bool
}

var Config ConfigType

func LoadConfig() ConfigType {
	//flag.BoolVar(&logDecodefailure, "l", false, "logDecodefailure")
	flag.BoolVar(&Config.logDecode, "d", false, "logDecode")
	flag.StringVar(&Config.keyword, "k", "", "keyword")
	flag.StringVar(&Config.path, "f", "/Users/xxxx/data.log", "logfile")
	flag.Parse()
	if strings.Contains(Config.keyword, "|") {
		Config.keywords = strings.Split(Config.keyword, "|")
	} else {
		Config.keywords = append(Config.keywords, Config.keyword)
	}
	Config.regexp = `(\d+\.\d+\.\d+\.\d+)\s-\s([^\s]*|\-)\s(\[[^\[\]]+\])\s([^\s]*)\s(\"(?:[^"]|\")+|-\")\s([^\s]*)\s(\"\d{3}\")\s(\d+|-)\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")`
	return Config
}

func init() {
	LoadConfig()
}

func main() {

	var wg sync.WaitGroup
	beginTime := time.Now().Unix()
	r := &ReadFormLogFile{
		path: Config.path,
	}
	lp := &LogProcess{
		rc:   make(chan LineMessage, 200),
		read: r,
	}
	wg.Add(1)
	go lp.read.Read(lp.rc, &wg)
	//for p := 0; p < 10; p++ {
	wg.Add(1)
	go lp.Process(&wg)
	//}
	wg.Wait()
	endTime := time.Now().Unix()
	usedTime := endTime - beginTime
	fmt.Printf("find %d rows\n", lp.hit)
	fmt.Printf("total line cont %d rows\n", r.TotalLineCont)
	fmt.Printf("Used for %d second\n", usedTime)
	fmt.Printf("avg 1 second proess %d row\n", r.TotalLineCont/int(usedTime))
	fmt.Println("decodeErrCont：" + strconv.Itoa(r.decodeErrCont))
}
