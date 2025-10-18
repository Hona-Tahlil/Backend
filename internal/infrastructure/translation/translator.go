package translation

import (
	"log"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fa_IR"
	ut "github.com/go-playground/universal-translator"
)

var translationMap map[string]map[string]string

func GetTranslator(locale string) ut.Translator {
	uniTrans := ut.New(fa_IR.New(), en.New(), fa_IR.New())

	addTranslations("fa_IR", Persian, uniTrans)

	trans, found := uniTrans.GetTranslator(locale)
	if !found {
		return uniTrans.GetFallback()
	}
	return trans
}

func addTranslations(key string, translations map[string]interface{}, universalTranslator *ut.UniversalTranslator) {
	translator, found := universalTranslator.GetTranslator(key)
	if !found {
		log.Fatal("translator not found")
	}
	flattenedTranslations := loadTranslation(key, translations)

	for key, translation := range flattenedTranslations {
		translator.Add(key, translation, true)
	}
}

func loadTranslation(key string, translations map[string]interface{}) map[string]string {
	if translations, found := translationMap[key]; found {
		return translations
	}
	flattenedTranslations := make(map[string]string)
	flattenMap("", translations, flattenedTranslations)
	return flattenedTranslations
}

func flattenMap(prefix string, input map[string]interface{}, output map[string]string) {
	for key, value := range input {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		switch v := value.(type) {
		case map[string]interface{}:
			flattenMap(fullKey, v, output)
		case string:
			output[fullKey] = v
		}
	}
}
