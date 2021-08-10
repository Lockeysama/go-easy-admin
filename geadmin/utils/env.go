package utils

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	beegoutils "github.com/beego/beego/v2/core/utils"
)

// EnvMap Config 文件变量集合
var EnvMap map[string]interface{}

// ReadLines 读取所有配置行（自动忽略注释行）
func ReadLines(filePth string) ([]string, error) {
	var err error
	lines := []string{}
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return lines, nil
			}
			return nil, err
		}
		lineStr := string(line)
		lineStr = strings.Replace(lineStr, "\n", "", -1)
		lineStr = strings.Replace(lineStr, " ", "", -1)
		if strings.Contains(lineStr, "=") && !strings.Contains(lineStr, "#") {
			lines = append(lines, lineStr)
		}
	}
}

// loadConfigFromFile 载入 Config 文件配置
func loadConfigFromFile() {
	var err error
	AppPath := ""
	WorkPath := ""
	if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	WorkPath, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	var filename = "app.conf"
	if os.Getenv("BEEGO_RUNMODE") != "" {
		filename = os.Getenv("BEEGO_RUNMODE") + ".app.conf"
	}
	appConfigPath := filepath.Join(WorkPath, "conf", filename)
	if !beegoutils.FileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, "conf", filename)
		if !beegoutils.FileExists(appConfigPath) {
			panic("config file missing")
		}
	}

	if lines, err := ReadLines(appConfigPath); err != nil {
		panic("read config failed")
	} else {
		EnvMap = make(map[string]interface{})
		for _, line := range lines {
			kvs := strings.Split(line, "=")
			if len(kvs) < 2 {
				panic("config syntax error")
			}
			EnvMap[kvs[0]] = strings.Join(kvs[1:], "=")
		}
	}
}

// Getenv 获取环境变量
func Getenv(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetenvFromConfig 从 config 文件获取环境变量
func GetenvFromConfig(key string, defaultValue interface{}) interface{} {
	if len(EnvMap) == 0 {
		loadConfigFromFile()
	}
	if value, ok := EnvMap[key]; !ok {
		return defaultValue
	} else {
		switch defaultValue.(type) {
		case int, int8, int16, int32, int64:
			if v, err := strconv.Atoi(value.(string)); err != nil {
				panic("config value error")
			} else {
				return int64(v)
			}
		default:
			return value
		}
	}
}
