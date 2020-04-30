package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

// 初始文件路径
var sourceFile string

// 替换项json 配置文件路径
var replaceFile string

// 输出替换后的文件路径
var outputFile string
var sectionAlreadyAppend = make(map[string]bool)

// 替换项map
var replaceMap map[string]string

func main() {
	// 控制台接受3个文件路径的参数
	source := flag.String("source", "", "source file")
	replace := flag.String("replace", "", "replace file (k=>v)")
	output := flag.String("output", "", "output file")

	help := flag.Bool("h", false, "this is help")

	//source := flag.String("source", "./file/config.ini", "source file")
	//replace := flag.String("replace", "./file/replace.json", "replace file (k=>v)")
	//output := flag.String("output", "./file/output.ini", "output file")

	flag.Parse()
	//
	sourceFile = *source
	replaceFile = *replace
	outputFile = *output
	if *help == true || sourceFile == "" || replaceFile == "" {
		usage()
		os.Exit(0)
	}
	fmt.Println("---start to replace config items---")
	checkFileExists(sourceFile)
	checkFileExists(replaceFile)
	if err := ParsedReplaceJson(replaceFile, &replaceMap); err != nil {
		panic(err)
	}
	if err := ReadFileByLine(sourceFile, HandlerLine); err != nil {
		panic(err)
	}
}

// line 行内容, section 当前行的组, 如[database]
func HandlerLine(line string, fileExt string, section string) {
	// 注释, 直接写入
	if ok, err := regexp.MatchString(`^#`, line); err != nil {
		panic(err)
	} else if ok {
		if err := AppendToFile(outputFile, line); err != nil {
			panic(err)
		}
		return
	}

	lineKey := GetLineKey(line)
	if lineKey == "" {
		return
	}
	// 从json的找到对应的替换配置
	mapKey := lineKey
	if fileExt == ".ini" {
		// section 仅写入一次
		if section != "" && sectionAlreadyAppend[section] == false {
			if err := AppendToFile(outputFile, fmt.Sprintf("\n[%s]", section)); err != nil {
				panic(err)
			}
			sectionAlreadyAppend[section] = true
		}
		mapKey = fmt.Sprintf("%s.%s", section, lineKey)
	}
	if repValue, ok := replaceMap[mapKey]; ok {
		line = fmt.Sprintf("%s=%s", lineKey, repValue)
		fmt.Println("Replace:", section, lineKey, repValue)
	}
	if err := AppendToFile(outputFile, line); err != nil {
		panic(err)
	}
}

// 逐行读取一个文件内容，并把行内容交给 handlerFunc 处理
// 文件最后一行必须是空行，否则读取不到最后一行配置项
func ReadFileByLine(fileName string, handlerFunc func(string, string, string)) error {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	// 当前行所在的section
	var section string
	fileExt := path.Ext(fileName)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// ini文件要读取 section 如[database]
		if fileExt == ".ini" {
			// 如果当前行是[section] 则跳过不处理
			if s := GetLineSection(line); s != "" {
				section = s
				continue
			} else {
				handlerFunc(line, fileExt, section)
			}
			// 其他文件直接替换
		} else {
			handlerFunc(line, fileExt, "")
		}
	}
}

//匹配section 如 [data]
func GetLineSection(line string) (section string) {
	matches := regexp.MustCompile(`^\[([a-zA-Z0-9_]+)\]`).FindStringSubmatch(line)
	if len(matches) == 2 {
		section = matches[1]
	}
	return
}

// 解析json
func ParsedReplaceJson(sourceFile string, res interface{}) error {
	fileBytes, err := ReadAll(sourceFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileBytes, res)
}

// 读取文件内容
func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// 使用io.WriteString()函数进行数据的写入
func AppendToFile(filename, content string) error {
	content = fmt.Sprintf("%s\n", content)
	fileObj, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	defer fileObj.Close()
	if err != nil {
		return err
	}
	if _, err := io.WriteString(fileObj, content); err == nil {
		return err
	}
	return nil
}

// 判断文件夹是否存在
func IsPathExists(dirPath string) (bool, error) {
	_, err := os.Stat(dirPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetLineKey(line string) string {
	// 当前行的key
	matches := regexp.MustCompile(`^(\S+)\s*=`).FindStringSubmatch(line)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func checkFileExists(fileName string) {
	if ok, err := IsPathExists(fileName); err != nil {
		panic(err)
	} else if !ok {
		panic(fmt.Sprintf("%s  not exists!", fileName))
	}
}
func parseSectionAndKey(confKey string) (section string, key string) {
	matches := regexp.MustCompile(`(\S+)\.(\S+)`).FindStringSubmatch(confKey)
	if len(matches) < 3 {
		return
	}
	section = matches[1]
	key = matches[2]
	return
}

func usage() {
	fmt.Fprintf(os.Stderr, `replace some k=>v in [source file] depends on [replace file] then output to [output file].
Usage: replace-key -source=filename.xxx -replace=filename.xxx -output=filename.xxx

Options:
`)
	flag.PrintDefaults()
}
