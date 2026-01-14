package utils

import (
	"crypto/md5"
	"encoding/hex"
	"image"
	"io"
	"os"

	"github.com/disintegration/imaging"
)

func FileMd5(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(md5hash.Sum(nil)), nil
}

func CutImage(oldFileName, newFileName string, width int) (err error) {
	file, err := os.Open(oldFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	c, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}

	if c.Width <= (width) {
		_, err = CopyFile(newFileName, oldFileName)
		return err
	}

	//按照宽度进行等比例缩放
	src, err := imaging.Open(oldFileName, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	src = imaging.Resize(src, width, 0, imaging.Lanczos)
	return imaging.Save(src, newFileName)
}

// CopyFile 文件复制
func CopyFile(dstFileName string, srcFileName string) (written int64, err error) {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return
	}
	defer srcFile.Close()

	//打开dstFileName
	dstFile, err := os.OpenFile(dstFileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, nil
}
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func WriteFile(name string, content []byte) error {
	return os.WriteFile(name, content, 0644)
}
func ReadFile(name string) (data []byte, err error) {
	return os.ReadFile(name)
}
