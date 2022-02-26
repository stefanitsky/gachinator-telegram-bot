package translations

type messageType int
type translation string
type LangCode string
type translations map[messageType]translation

var (
	Welcome        messageType = 0
	Error          messageType = 1
	SelectRussian  messageType = 2
	SelectEnglish  messageType = 3
	SelectLangMenu messageType = 4
)

var (
	Russian LangCode = "ru"
	English LangCode = "en"
)

var (
	russianMessages = translations{
		Welcome:        "Привет! Тебе был установлен русский язык. Для смены языка введи команду /lang.",
		Error:          "Ой, произошла ошибка, пожалуйста, попробуйте еще раз.",
		SelectRussian:  "Установлен русский язык.",
		SelectEnglish:  "Установлен английский язык.",
		SelectLangMenu: "Выберите язык...",
	}
	englishMessages = translations{
		Welcome:        "Hello! Default language is set to english. To change type /lang command.",
		Error:          "Oops, an error occurred, please, try again later...",
		SelectRussian:  "Russian language is set.",
		SelectEnglish:  "English language is set.",
		SelectLangMenu: "Select a language...",
	}
)

var langCodeToTranslations = map[LangCode]translations{
	Russian: russianMessages,
	English: englishMessages,
}

// GetMessage returns localized error message for a specified language code
func GetMessage(mt messageType, lc LangCode) string {
	return string(langCodeToTranslations[lc][mt])
}

// GetErrorMessage returns default error messsage for a specified language code
func GetErrorMessage(lc LangCode) string {
	return GetMessage(Error, lc)
}
