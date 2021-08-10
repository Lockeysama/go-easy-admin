package utils

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	// StorageMode 文件存储模式 Local | S3
	StorageMode string
	// StorageTarget 对象存储对象(aliyun、aws 等)
	StorageTarget string
	// StorageTimeout 存储超时时间
	StorageTimeout time.Duration
	// StorageEndpoint 存储端点
	StorageEndpoint string
	// StoragePath 存储路径
	StoragePath string
	// StorageBucket 存储 Bucket
	StorageBucket string
	// StorageRegion 存储 Region
	StorageRegion string
	// StorageKey 存储 AccessKey
	StorageKey string
	// StorageSecret 存储 AccessSecret
	StorageSecret string
	// StoragePresignTime 获取存储文件授权有效时间
	StoragePresignTime time.Duration
)

var (
	// StorageModeLocal 本地存储模式
	StorageModeLocal string = "local"
	// StorageModeOS 第三方对象存储
	StorageModeOS string = "os"
)

func init() {
	StorageMode = GetenvFromConfig("storage.mode", "local").(string)
	if StorageMode == StorageModeLocal {
		StoragePath = GetenvFromConfig("storage.path", "static/upload/").(string)
	} else if StorageMode == StorageModeOS {
		StorageTarget = GetenvFromConfig("storage.target", "").(string)
		StoragePath = GetenvFromConfig("storage.path", "").(string)
		StorageEndpoint = GetenvFromConfig("storage.endpoint", "").(string)
		StorageBucket = GetenvFromConfig("storage.bucket", "").(string)
		StorageKey = GetenvFromConfig("storage.key", "").(string)
		StorageSecret = GetenvFromConfig("storage.secret", "").(string)
		if StorageBucket == "" || StorageKey == "" || StorageSecret == "" {
			panic("storage setting missing")
		}
		StorageTimeout = time.Duration(GetenvFromConfig("storage.timeout", int64(15)).(int64)) * time.Second
		StoragePresignTime = time.Duration(GetenvFromConfig("storage.presigntime", int64(3600)).(int64)) * time.Second
	} else {
		panic("unknown storage mode")
	}
}

// FileExt 判断文件后缀
func FileExt(file string) string {
	ext := strings.Split(file, ".")
	if len(ext) > 1 {
		return ext[1]
	}
	return ""
}

// PutFileToOS 上传文件到 OS
func PutFileToOS(buffer []byte, path string) error {
	client, err := oss.New(StorageEndpoint, StorageKey, StorageSecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(StorageBucket)
	if err != nil {
		return err
	}
	return bucket.PutObject(path, bytes.NewReader(buffer))
}

// GetFileURLFromOS 上传文件到 OS
func GetFileFromOS(path string) (io.ReadCloser, error) {
	client, err := oss.New(StorageEndpoint, StorageKey, StorageSecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(StorageBucket)
	if err != nil {
		return nil, err
	}
	return bucket.GetObject(path)
}

// GetFileURLFromOS 上传文件到 OS
func GetFileURLFromOS(path string) (string, error) {
	return PresignRequest("GET", path)
}

// PresignRequest 从 AWS S3 获取文件授权 URL 或 上传
func PresignRequest(action string, path string) (url string, err error) {
	action = strings.ToUpper(action)
	if action != "PUT" && action != "GET" && action != "POST" {
		err = errors.New("method not allowed")
		return "", err
	}

	client, err := oss.New(StorageEndpoint, StorageKey, StorageSecret)
	if err != nil {
		return "", err
	}

	bucket, err := client.Bucket(StorageBucket)
	if err != nil {
		return "", err
	}

	url, err = bucket.SignURL(path, oss.HTTPMethod(action), int64(StoragePresignTime/time.Second))
	url = strings.Replace(url, "http://", "https://", -1)
	return url, err
}
