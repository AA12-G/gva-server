package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	uploadDir    = "uploads"
	maxFileSize  = 5 << 20 // 5MB
	allowedTypes = ".jpg,.jpeg,.png,.gif"
)

// SaveUploadedFile 保存上传的文件
func SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	// 检查文件大小
	if file.Size > maxFileSize {
		return "", fmt.Errorf("文件大小不能超过5MB")
	}

	// 检查文件类型
	ext := strings.ToLower(path.Ext(file.Filename))
	if !strings.Contains(allowedTypes, ext) {
		return "", fmt.Errorf("只支持jpg、jpeg、png、gif格式")
	}

	// 创建上传目录
	uploadPath := fmt.Sprintf("%s/%s", uploadDir, time.Now().Format("2006/01/02"))
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}

	// 生成文件名
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := fmt.Sprintf("%s/%s", uploadPath, fileName)

	// 保存文件
	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 将文件内容复制到目标文件
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// 返回相对URL路径，而不是文件系统路径
	return "/" + filePath, nil // 添加前导斜杠，使其成为URL路径
}
