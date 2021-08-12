package geacontrollers

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	"github.com/lockeysama/go-easy-admin/geadmin/utils"

	blurhash "github.com/buckket/go-blurhash"
)

// AjaxUpload 上传文件
func (c *ManageBaseController) AjaxUpload() {
	fh := c.Ctx().RequestMultipartForm().File["file"][0]
	file, _ := fh.Open()

	path := c.Ctx().RequestURL().Query().Get("path")
	if path == "" {
		c.AjaxMsg("path not found", MSG_ERR)
		return
	}
	fileName := strings.Split(fh.Filename, ".")
	path = path + fileName[0] + "/"

	size := fh.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	blur := ""
	if c.Ctx().RequestURL().Query().Get("blur") == "true" {
		f := bytes.NewReader(buffer)
		var loadedImage image.Image
		var err error
		switch utils.FileExt(fh.Filename) {
		case "png":
			if loadedImage, err = png.Decode(f); err != nil {
				c.APIRequestError(400, "生成 blur hash 失败, 图片文件解码失败")
			}
		case "jpg", "jpeg":
			if loadedImage, err = jpeg.Decode(f); err != nil {
				c.APIRequestError(400, "生成 blur hash 失败, 图片文件解码失败")
			}
		default:
			c.APIRequestError(400, "生成 blur hash 失败; 当前只 png、jpg/jpeg 图片格式")
			return
		}

		if str, err := blurhash.Encode(5, 5, loadedImage); err != nil {
			c.APIRequestError(400, "生成 blur hash 失败")
			return
		} else {
			input := []byte(str)
			blur = base64.StdEncoding.EncodeToString(input)
			blur = strings.ReplaceAll(blur, "+", "-")
			blur = strings.ReplaceAll(blur, "/", "_")
		}
	}

	fMD5 := md5.New()
	io.Copy(fMD5, bytes.NewReader(buffer))

	filePath := path + hex.EncodeToString(fMD5.Sum(nil)) + "__" + blur + "__" + fh.Filename

	data := make(map[string]interface{})
	switch utils.StorageMode {
	case utils.StorageModeLocal:
		f, err := os.OpenFile(utils.StoragePath+filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		defer f.Close()
		if err != nil {
			c.AjaxMsg("失败", MSG_ERR)
			return
		}
		f.Write(buffer)
		data["fullPath"] = utils.StoragePath + filePath
		data["path"] = utils.StoragePath + filePath
	case utils.StorageModeOS:
		if err := utils.PutFileToOS(buffer, utils.StoragePath+filePath); err != nil {
			fmt.Println(err.Error())
			c.AjaxMsg("get file url to os failed", MSG_ERR)
			return
		}
		if url, err := utils.GetFileURLFromOS(utils.StoragePath + filePath); err != nil {
			fmt.Println(err.Error())
			c.AjaxMsg("get file url from os failed", MSG_ERR)
			return
		} else {
			data["fullPath"] = url
			data["path"] = utils.StoragePath + filePath
		}
	}

	c.AjaxData(data, MSG_OK)
}

// AjaxGetFile 上传文件
func (c *ManageBaseController) AjaxGetFile() {
	if c.User == nil && c.APIUser == nil {
		c.AjaxMsg("get file failed", MSG_ERR)
		return
	}
	path := strings.Split(c.Ctx().RequestURL().RawQuery, "=")
	if len(path) < 1 {
		c.AjaxMsg("path not found", MSG_ERR)
		return
	}
	filePath := path[1]
	if filePath == "" {
		c.AjaxMsg("path empty", MSG_ERR)
		return
	}

	data := make(map[string]interface{})
	switch utils.StorageMode {
	case utils.StorageModeLocal:
		data["fullPath"] = utils.StoragePath + filePath
		data["path"] = filePath
	case utils.StorageModeOS:
		if url, err := utils.GetFileURLFromOS(filePath); err != nil {
			c.AjaxMsg("get file url from aws s3 failed", MSG_ERR)
			return
		} else {
			data["fullPath"] = url
			data["path"] = filePath
		}
	}

	c.AjaxData(data, MSG_OK)
}
