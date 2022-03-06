package translations

type messageType int
type translation string
type LangCode string
type translations map[messageType]translation

var (
	Welcome        messageType = 0
	Error          messageType = 1
	SelectNewLang  messageType = 2
	SelectLangMenu messageType = 3
)

var (
	Russian LangCode = "ru"
	English LangCode = "en"
)

var (
	russianMessages = translations{
		Welcome:        "Привет! Тебе был установлен русский язык. Для смены языка введи команду /lang.",
		Error:          "Ой, произошла ошибка, пожалуйста, попробуйте еще раз.",
		SelectNewLang:  "Язык изменен.",
		SelectLangMenu: "Выберите язык...",
	}
	englishMessages = translations{
		Welcome:        "Hello! Default language is set to english. To change type /lang command.",
		Error:          "Oops, an error occurred, please, try again later...",
		SelectNewLang:  "Language is changed.",
		SelectLangMenu: "Select a language...",
	}
)

var langCodeToTranslations = map[LangCode]translations{
	Russian: russianMessages,
	English: englishMessages,
}

var DefaultLangCode = Russian

// GetMessage returns localized error message for a specified language code
func GetMessage(mt messageType, lc LangCode) string {
	return string(langCodeToTranslations[lc][mt])
}

// GetErrorMessage returns default error messsage for a specified language code
func GetErrorMessage(lc LangCode) string {
	return GetMessage(Error, lc)
}
