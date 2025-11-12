package domaintranslation

import ut "github.com/go-playground/universal-translator"

type Translator interface {
	GetTranslator(locale string) ut.Translator
}
