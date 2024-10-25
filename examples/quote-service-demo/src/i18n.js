import { createI18n } from 'vue-i18n'
import en from './locales/en.json'
import ja from './locales/ja.json'
import de from './locales/de.json'
import lo from './locales/lo.json'
import es from './locales/es.json'
import fr from './locales/fr.json'
import bo from './locales/bo.json'
import ug from './locales/ug.json'
import sn from './locales/sn.json'
import zh from './locales/zh.json'

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en,
    ja,
    de,
    lo,
    es,
    fr,
    bo,
    ug,
    sn,
    zh
  }
})

export default i18n
