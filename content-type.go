package utils

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// GetContentTypeByContent 通过文件内容检测 MIME 类型
func GetContentTypeByContent(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败：%w", err)
	}
	defer file.Close()

	// 读取文件前 512 字节（DetectContentType 仅需前 512 字节）
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("读取文件内容失败：%w", err)
	}

	// 检测 MIME 类型
	contentType := http.DetectContentType(buffer[:n])
	return contentType, nil
}

// GetContentTypeByExtension 通过文件扩展名推断 MIME 类型
func GetContentTypeByExtension(filePath string) string {
	// 获取文件扩展名（含 .）
	ext := filepath.Ext(filePath)
	// 根据扩展名映射 MIME 类型（内部维护了常见扩展名的映射表）
	return mime.TypeByExtension(ext)
}
