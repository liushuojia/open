package gin_middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	utils "github.com/liushuojia/open"
)

// 定义上下文键名
const (
	LangKey  = "lang"  // 存储语言标识
	TransKey = "trans" // 存储翻译函数
)

var i18n *utils.I18n

// I18nMiddleware Gin 国际化中间件：解析语言标识并挂载翻译函数
func I18nMiddleware(defaultLanguage string, supportLanguages []string, localeDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err  error
			lang string
		)
		if i18n == nil {
			i18n, err = utils.NewI18n(defaultLanguage, supportLanguages, localeDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "请设置语言文件夹 - " + err.Error(),
				})
				c.Abort()
			}
		}

		// 1. 优先解析 URL 参数（?lang=en-US）
		lang = c.Query("lang")
		if isSupportLang(lang, supportLanguages) {
			setLangToContext(c, lang)
			c.Next()
			return
		}

		// 2. 解析 Header（Accept-Language: en-US,en;q=0.9）
		acceptLang := c.GetHeader("Accept-Language")
		if acceptLang != "" {
			// 提取第一个语言标识（如 en-US）
			lang = strings.Split(acceptLang, ",")[0]
			lang = strings.Split(lang, ";")[0]
			// 标准化语言标识（如 en → en-US，zh → zh-CN）
			lang = normalizeLang(lang)
			if isSupportLang(lang, supportLanguages) {
				setLangToContext(c, lang)
				c.Next()
				return
			}
		}

		// 3. 解析 Cookie
		cookieLang, _ := c.Cookie("lang")
		if isSupportLang(cookieLang, supportLanguages) {
			setLangToContext(c, cookieLang)
			c.Next()
			return
		}

		// 4. 使用默认语言
		setLangToContext(c, defaultLanguage)
		c.Next()
	}
}

// setLangToContext 将语言标识和翻译函数挂载到上下文
func setLangToContext(c *gin.Context, lang string) {
	c.Set(LangKey, lang)
	// 挂载翻译函数（简化接口内调用）
	c.Set(TransKey, func(key string, data map[string]interface{}) (string, error) {
		return i18n.Translate(lang, key, data)
	})
}

// isSupportLang 检查语言是否受支持
func isSupportLang(language string, supportLanguages []string) bool {
	if language == "" {
		return false
	}
	for _, s := range supportLanguages {
		if s == language {
			return true
		}
	}
	return false
}

// normalizeLang 标准化语言标识
func normalizeLang(language string) string {
	switch language {
	case "en":
		return "en-US"
	case "zh":
		return "zh-CN"
	case "zh-TW", "zh-HK":
		return "zh-TW"
	default:
		return language
	}
}
