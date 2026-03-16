package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	utils "github.com/liushuojia/open"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {

	// 1. 初始化 Bundle（文案容器，默认语言：简体中文）
	bundle := i18n.NewBundle(language.Chinese)
	// 注册 TOML 解析器（修复 UnmarshalTOML 找不到问题）
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	d := utils.LocalDir("open")

	test2(d + "/demo/locales")
	return

	// 2. 加载所有文案文件
	localesDir := d + "/demo/locales"
	files, _ := os.ReadDir(localesDir)
	for _, f := range files {
		filePath := filepath.Join(localesDir, f.Name())
		bundle.LoadMessageFile(filePath)
	}

	// 3. 创建翻译器（指定语言）
	zhLocalizer := i18n.NewLocalizer(bundle, "zh-CN") // 简体
	enLocalizer := i18n.NewLocalizer(bundle, "en-US") // 英文

	// 4. 基础翻译（带变量）
	zhMsg, _ := zhLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "greet",
		TemplateData: map[string]interface{}{"name": "张三"},
	})
	enMsg, _ := enLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "greet",
		TemplateData: map[string]interface{}{"name": "Zhang San"},
	})

	fmt.Println("简体：", zhMsg) // 输出：你好，张三！
	fmt.Println("英文：", enMsg) // 输出：Hello, Zhang San!

	// 5. 嵌套文案翻译
	errMsg, _ := zhLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "error.not_found",
		TemplateData: map[string]interface{}{"id": 1001},
	})
	fmt.Println("错误提示：", errMsg) // 输出：资源 1001 不存在

	// 4. 复数翻译（关键：指定 PluralCount）
	// 测试数量=0
	zh0, _ := zhLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "apple_count",
		TemplateData: map[string]interface{}{"count": 0},
		PluralCount:  0, // 必须指定复数计数
	})
	en0, _ := enLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "apple_count",
		TemplateData: map[string]interface{}{"count": 0},
		PluralCount:  0,
	})

	// 测试数量=1
	zh1, _ := zhLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "apple_count",
		TemplateData: map[string]interface{}{"count": 1},
		PluralCount:  1,
	})
	en1, _ := enLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "apple_count",
		TemplateData: map[string]interface{}{"count": 1},
		PluralCount:  1,
	})

	// 测试数量=5
	zh5, _ := zhLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "apple_count",
		TemplateData: map[string]interface{}{"count": 5},
		PluralCount:  5,
	})
	en5, _ := enLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "apple_count",
		TemplateData: map[string]interface{}{"count": 5},
		PluralCount:  5,
	})

	// 输出结果
	fmt.Println("中文 0 个：", zh0) // 你有 0 个苹果
	fmt.Println("英文 0 个：", en0) // You have 0 apples
	fmt.Println("中文 1 个：", zh1) // 你有 1 个苹果
	fmt.Println("英文 1 个：", en1) // You have 1 apple
	fmt.Println("中文 5 个：", zh5) // 你有 5 个苹果
	fmt.Println("英文 5 个：", en5) // You have 5 apples

	loginSuccess, _ := zhLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "user.login.success",
		TemplateData: map[string]interface{}{"username": "张三"},
	})
	fmt.Println(loginSuccess) // 登录成功，张三！

	// 调用 user.profile.avatar_error
	avatarErr, _ := zhLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "user.profile.avatar_error",
		TemplateData: map[string]interface{}{"maxSize": 5},
	})
	fmt.Println(avatarErr) // 头像上传失败：文件大小超过 5MB

	loginSuccessen, _ := enLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "user.login.success",
		TemplateData: map[string]interface{}{"username": "张三"},
	})
	fmt.Println(loginSuccessen) // 登录成功，张三！

	// 调用 user.profile.avatar_error
	avatarErren, _ := enLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID:    "user.profile.avatar_error",
		TemplateData: map[string]interface{}{"maxSize": 5},
	})
	fmt.Println(avatarErren) // 头像上传失败：文件大小超过 5MB

}

func test2(localeDir string) {

	v, _ := utils.NewI18n("zh-CN", []string{
		"zh-CN",
		"en-US",
	}, localeDir)

	fmt.Println(v.Translate("zh-CN", "user.login.success", map[string]interface{}{"username": "张三"}))
	fmt.Println(v.Translate("en-US", "user.login.success", map[string]interface{}{"username": "zhang san"}))

}
