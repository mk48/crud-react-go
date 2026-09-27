import i18n from "i18next"
import { initReactI18next } from "react-i18next"
import enTranslation from "@/locales/en/translation.json"

const en = {
  translation: enTranslation,
}

i18n.use(initReactI18next).init(
  {
    lng: "en",
    debug: false,
    fallbackLng: "en",
    keySeparator: ".",
    interpolation: {
      escapeValue: false,
    },
    resources: { en },
    ns: ["translation"],
    defaultNS: "translation",
    fallbackNS: "translation",
  },
  (error) => {
    if (error) console.error(error)
  }
)

export default i18n
