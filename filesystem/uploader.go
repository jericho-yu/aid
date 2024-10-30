package filesystem

import (
	"errors"
	"fmt"

	"github.com/jericho-yu/aid/httpClient"
)

type (
	FileManager struct {
		Err            error
		dstDir, srcDir string
		fileBytes      []byte
		fileSize       int64
		config         *FileManagerConfig
	}
	FileManagerConfigDriver = string
	FileManagerConfig       struct {
		Username  string
		Password  string
		AuthTitle string
		Driver    FileManagerConfigDriver
	}
)

var FileManagerApp FileManager

const (
	FileManagerConfigDriverLocal FileManagerConfigDriver = "LOCAL"
	FileManagerConfigDriverNexus FileManagerConfigDriver = "NEXUS"
	FileManagerConfigDriverOss   FileManagerConfigDriver = "OSS"
)

// New 初始化：文件管理
func (FileManager) New(config *FileManagerConfig) *FileManager {
	return &FileManager{
		config: config,
	}
}

// NewByLocalFile 初始化：文件管理器（通过本地文件）
func (FileManager) NewByLocalFile(srcDir, dstDir string, config *FileManagerConfig) (*FileManager, error) {
	fs := FileSystemApp.NewByAbs(srcDir)
	if !fs.IsExist {
		return nil, errors.New("目标文件不存在")
	}

	fileBytes, err := fs.Read()
	if err != nil {
		return nil, err
	}

	return &FileManager{
		dstDir:    dstDir,
		srcDir:    srcDir,
		fileBytes: fileBytes,
		fileSize:  int64(len(fileBytes)),
		config:    config,
	}, nil
}

// NewByBytes 实例化：文件管理器（通过字节）
func (FileManager) NewByBytes(srcFileBytes []byte, dstDir string, config *FileManagerConfig) *FileManager {
	return &FileManager{
		dstDir:    dstDir,
		fileBytes: srcFileBytes,
		fileSize:  int64(len(srcFileBytes)),
		config:    config,
	}
}

// SetSrcDir 设置源文件
func (r *FileManager) SetSrcDir(srcDir string) (*FileManager, error) {
	fs := FileSystemApp.NewByAbs(srcDir)
	if !fs.IsExist {
		return nil, errors.New("目标文件不存在")
	}

	fileBytes, err := fs.Read()
	if err != nil {
		return nil, err
	}

	r.srcDir = fs.GetDir()
	r.fileBytes = fileBytes
	r.fileSize = int64(len(fileBytes))

	return r, nil
}

// SetDstDir 设置目标目录
func (r *FileManager) SetDstDir(dstDir string) *FileManager {
	r.dstDir = dstDir
	return r
}

// Upload 上传文件
func (r *FileManager) Upload() (int64, error) {
	switch r.config.Driver {
	case FileManagerConfigDriverLocal:
		return r.uploadToLocal()
	case FileManagerConfigDriverNexus:
		return r.uploadToNexus()
	case FileManagerConfigDriverOss:
		return r.uploadToOss()
	}

	return 0, fmt.Errorf("不支持的驱动类型：%s", r.config.Driver)
}

// Delete 删除文件
func (r *FileManager) Delete() error {
	switch r.config.Driver {
	case FileManagerConfigDriverLocal:
		return r.deleteFromLocal()
	case FileManagerConfigDriverNexus:
		return r.deleteFromNexus()
	case FileManagerConfigDriverOss:
		return r.deleteFromOss()
	}
	return fmt.Errorf("不支持的驱动类型：%s", r.config.Driver)
}

// 上传到本地
func (r *FileManager) uploadToLocal() (int64, error) {
	dst := FileSystemApp.NewByAbs(r.dstDir)
	return dst.WriteBytes(r.fileBytes)
}

// 上传到nexus
func (r *FileManager) uploadToNexus() (int64, error) {
	client := httpClient.NewPut(r.dstDir).
		SetAuthorization(r.config.Username, r.config.Password, r.config.AuthTitle).
		AddHeaders(map[string][]string{
			"Content-Length": {fmt.Sprintf("%d", r.fileSize)},
		}).
		SetBody(r.fileBytes).
		Send()

	if client.Err != nil {
		return 0, client.Err
	}

	return int64(len(r.fileBytes)), nil
}

// 上传到oss
func (r *FileManager) uploadToOss() (int64, error) {
	return 0, errors.New("暂不支持oss方式")
}

// 从本地删除文件
func (r *FileManager) deleteFromLocal() error {
	return FileSystemApp.NewByAbs(r.dstDir).DelFile()
}

// 从nexus删除文件
func (r *FileManager) deleteFromNexus() error {
	return httpClient.NewDelete(r.dstDir).
		SetAuthorization(r.config.Username, r.config.Password, r.config.AuthTitle).
		Send().Err
}

// 从oss删除文件
func (r *FileManager) deleteFromOss() error {
	return errors.New("暂不支持oss方式")
}
