package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var keyword string
var logfile string
var keywords []string
var hit int64
var logDecode bool
var logDecodefailure bool

func init() {
	//flag.StringVar(&keyword, "k", "34D82150-53D9-489C-A1E1-7D9FB45D0495", "keyword")
	//flag.StringVar(&keyword, "k", "47820001000603", "keyword")
	//flag.StringVar(&logfile, "f", "/Users/xxxx/data.log-20190904.gz", "keyword")
	flag.StringVar(&keyword, "k", "OS_RAPK3K6W1EKPL3R6C", "keyword")
	flag.StringVar(&logfile, "f", "logs/data.log", "logfile")
	flag.BoolVar(&logDecode, "d", false, "logDecode")
	flag.BoolVar(&logDecodefailure, "l", false, "logDecodefailure")
	flag.Parse()
}

var logs = logrus.New()

func hasKeyword(txt string, line string) {
	//fmt.Println(txt)
	var has = true
	for _, word := range keywords {
		if strings.Contains(txt, word) == false {
			has = false
			break
		}
	}
	if has {
		hit++
		fmt.Print(line)
		if logDecode {
			fmt.Println(txt, "\n")
		}
	}
}

//type LogProcess struct {
//	path string
//	dc chan string
//}

// https://www.cnblogs.com/lavin/p/5373188.html
func Base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	//res, err := base64.URLEncoding.DecodeString(data)
	//fmt.Println("  decodebase64urlsafe is :", string(res), err)
	return base64.URLEncoding.DecodeString(data)
}

//func (l *LogProcess) Read()  {
//
//}
//
//func (l *LogProcess) Process(dc chan string)  {
//
//}


func main() {


}




func test()  {
	beginTime := time.Now().Unix()
	//logs.SetFormatter(&logrus.JSONFormatter{})
	logs.SetFormatter(&logrus.TextFormatter{})

	file, err := os.OpenFile("icat.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logs.Out = file
	} else {
		logs.Info("Failed to log to file, using default stderr")
	}
	if keyword == "" {
		return
	}
	if strings.Contains(keyword, "|") {
		keywords = strings.Split(keyword, "|")
	} else {
		keywords = append(keywords, keyword)
	}
	defer glog.Flush()

	f, err := os.Open(logfile)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}
	//f.Seek(2);
	r := regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)\s-\s([^\s]*|\-)\s(\[[^\[\]]+\])\s([^\s]*)\s(\"(?:[^"]|\")+|-\")\s([^\s]*)\s(\"\d{3}\")\s(\d+|-)\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")\s(\"(?:[^"]|\")+|-\")`)
	rd := bufio.NewReader(f)
	defer f.Close()

	var decodeErrCont = 0
	var totalLineCont = 0
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			glog.Info("readbytes error:", err.Error())
			break
		} else if err != nil {
			glog.Info("readbytes error:", err.Error())
		}
		totalLineCont++

		result := r.FindStringSubmatch(string(line))
		if len(result) != 12 {
			glog.Info("FindStringSubmatch error:", string(line))
			logs.Error("FindStringSubmatch error:", string(line))
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
						hasKeyword(t.Get("data"), string(line))
					}
				}

				//fmt.Println("error:", result[6])
				// get 参数
			} else {
				if strings.Contains(m.Get("data"), "{") {
					hasKeyword(m.Get("data"), string(line))
				} else {



					encrypt := strings.ReplaceAll(m.Get("data"), "x0A", "---")
					encrypt = strings.ReplaceAll(encrypt, `\`, "-")
					encrypt = strings.ReplaceAll(encrypt, `----`, "")

					var missing = (4 - len(encrypt)%4) % 4
					encrypt += strings.Repeat("=", missing)
					//res, err := base64.URLEncoding.DecodeString(data)
					//fmt.Println("  decodebase64urlsafe is :", string(res), err)
					dst,_ :=  base64.URLEncoding.DecodeString(encrypt)

					//dst,_ := Base64URLDecode(encrypt)
					//fmt.Println("Base64URLDecode:", string(dat))
					//
					//encrypt = strings.ReplaceAll(encrypt, `===`, "+++")
					//dst, err := base64.StdEncoding.DecodeString(encrypt) // 490  v 5000
					//

					if err != nil {
						decodeErrCont++
						fmt.Println("替换后:", encrypt)
						fmt.Println("解码败:", m.Get("data"))
						fmt.Println(string(line))
						break


						logs.Error("base64Decode error:", string(line))
						continue

						//fmt.Println("error:", err)
					}
					hasKeyword(string(dst), string(line))
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
						hasKeyword(getParam, string(line))
					} else {
						// error
						logs.Error("params format error:", string(line))
					}
					//Str ,_:= hex.DecodeString(arslice[0])
					//fmt.Println(result[5])
					//fmt.Println("解码前：",arslice[0])
					//fmt.Println("解码后：",Str)
				} else {
					// error
					logs.Error("params format error:", string(line))
				}
				//Str ,_:= hex.DecodeString(result[5])
				//fmt.Println(result[5])
				//fmt.Println("解码后：",Str)
			} else {
				hasKeyword(m.Get("data"), string(line))
			}
			//fmt.Println(m.Get("data"))
		}
		//fmt.Println(len(result))
	}
	endTime := time.Now().Unix()
	usedTime := endTime - beginTime
	fmt.Printf("find %d rows\n", hit)
	fmt.Printf("total line cont %d rows\n", totalLineCont)
	fmt.Printf("Used for %d second\n", usedTime)
	fmt.Printf("avg 1 second proess %d row\n", totalLineCont/int(usedTime))
	fmt.Println("decodeErrCont：" + strconv.Itoa(decodeErrCont))
}