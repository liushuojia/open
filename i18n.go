package utils

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type I18n struct {
	defaultLanguage string
	languages       []string
	bundle          *i18n.Bundle
	localizers      map[string]*i18n.Localizer
	mu              sync.RWMutex
}

/*	NewI18n		defaultLanguage 默认语言		languages 语言列表	locales 语言文件	*/
func NewI18n(defaultLanguage string, languages []string, localesDir string) (*I18n, error) {
	conn := &I18n{
		defaultLanguage: defaultLanguage,
		languages:       languages,
		bundle:          i18n.NewBundle(language.MustParse(defaultLanguage)),
		localizers:      make(map[string]*i18n.Localizer),
	}

	conn.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 2. 加载所有文案文件
	files, _ := os.ReadDir(localesDir)
	for _, f := range files {
		filePath := filepath.Join(localesDir, f.Name())
		_, _ = conn.bundle.LoadMessageFile(filePath)
	}

	for _, lang := range languages {
		conn.localizers[lang] = i18n.NewLocalizer(conn.bundle, lang)
	}

	return conn, nil
}

func (conn *I18n) Translate(lang, key string, data map[string]interface{}) (string, error) {
	localizerTranslate, ok := conn.localizers[lang]
	if !ok {
		return "", errors.New("not found language")
	}
	return localizerTranslate.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
}
func (conn *I18n) TranslatePlural(lang, key string, data map[string]interface{}, count int) (string, error) {
	localizerTranslate, ok := conn.localizers[lang]
	if !ok {
		return "", errors.New("not found language")
	}
	return localizerTranslate.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
		PluralCount:  count,
	})
}
