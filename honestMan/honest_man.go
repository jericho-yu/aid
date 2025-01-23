package honestMan

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"

	"gopkg.in/yaml.v2"
)

type HonestMan struct{ dir string }

var App HonestMan

func (HonestMan) New(dirs ...string) *HonestMan {
	return &HonestMan{dir: path.Join(dirs...)}
}

// 读取文件
func (my *HonestMan) readFile() []byte {
	var (
		fileContent []byte
		err         error
	)
	fileContent, err = os.ReadFile(my.dir)
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败(%s)：%s", my.dir, err.Error()))
	}

	return fileContent
}

// 检查参数是否是一个指针
func (my *HonestMan) isPtr(target any) {
	// 使用反射检查target是否为指针类型
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		panic(errors.New("参数必须是一个指针"))
	}
}

// LoadYaml 读取Yaml配置文件
func (my *HonestMan) LoadYaml(target any) (err error) {
	my.isPtr(target)
	return yaml.Unmarshal(my.readFile(), target)
}

// LoadJson 读取Json配置文件
func (my *HonestMan) LoadJson(target any) (err error) {
	my.isPtr(target)

	return json.Unmarshal(my.readFile(), target)
}

// SaveYaml 写入Yaml文件
func (my *HonestMan) SaveYaml(target any) (err error) {
	my.isPtr(target)
	out, err := yaml.Marshal(target)
	if err != nil {
		return err
	}

	return os.WriteFile(my.dir, out, os.ModePerm)
}

// SaveJson 写入Json文件
func (my *HonestMan) SaveJson(target any) (err error) {
	my.isPtr(target)
	out, err := json.Marshal(target)
	if err != nil {
		return err
	}

	return os.WriteFile(my.dir, out, os.ModePerm)
}
