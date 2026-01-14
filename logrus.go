package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func InitializeLog(env string) {
	//创建日志文件路径
	logPath := "logs/" + env + "/app.log"

	//初始化 file-rotatelogs
	logWriter, err := rotatelogs.New(
		filepath.Join("logs", env, "app_%Y%m%d.log"), // 日志文件名格式
		rotatelogs.WithLinkName(logPath),             // 生成软连接指向最新日志文件
		rotatelogs.WithMaxAge(30*24*time.Hour),       // 最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour),    // 每天切割一次
	)
	if err != nil {
		log.Fatal("fail to initialize log writer:", err)
	}

	log.SetFormatter(&MyFormatter{})

	// 配置 logrus 的日志格式
	//log.SetFormatter(&log.TextFormatter{
	//	TimestampFormat: "2006-01-02 15:04:05",
	//})

	//配置日志级别
	log.SetLevel(log.DebugLevel)

	gin.DefaultWriter = io.MultiWriter(logWriter, os.Stdout)
	log.SetOutput(io.MultiWriter(logWriter, os.Stdout))
}

// MyFormatter 自定义格式器
type MyFormatter struct{}

func (f *MyFormatter) Format(entry *log.Entry) ([]byte, error) {
	// 格式化时间
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	// 构建日志内容（格式：[时间] [级别] 消息 [额外字段]）
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, entry.Level, entry.Message)

	// 附加额外字段
	for k, v := range entry.Data {
		logLine += fmt.Sprintf(" %s=%v", k, v)
	}
	logLine += "\n" // 换行

	return []byte(logLine), nil
}
