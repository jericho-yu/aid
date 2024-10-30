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
func (r *HonestMan) readFile() []byte {
	var (
		fileContent []byte
		err         error
	)
	fileContent, err = os.ReadFile(r.dir)
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败(%s)：%s", r.dir, err.Error()))
	}

	return fileContent
}

// 检查参数是否是一个指针
func (r *HonestMan) isPtr(target any) {
	// 使用反射检查target是否为指针类型
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr {
		panic(errors.New("参数必须是一个指针"))
	}
}

// LoadYaml 读取Yaml配置文件
func (r *HonestMan) LoadYaml(target any) (err error) {
	r.isPtr(target)
	return yaml.Unmarshal(r.readFile(), target)
}

// LoadJson 读取Json配置文件
func (r *HonestMan) LoadJson(target any) (err error) {
	r.isPtr(target)

	return json.Unmarshal(r.readFile(), target)
}
