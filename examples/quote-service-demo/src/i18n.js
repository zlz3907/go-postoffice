import { createI18n } from 'vue-i18n'
import zh from './locales/zh.json'
import en from './locales/en.json'
import sn from './locales/sn.json'
import lo from './locales/lo.json'
import ug from './locales/ug.json'
import bo from './locales/bo.json'
import de from './locales/de.json'
import fr from './locales/fr.json'
import es from './locales/es.json'

const messages = {
  zh,
  en,
  sn,
  lo,
  ug,
  bo,
  de,
  fr,
  es
}

export const i18n = createI18n({
  legacy: false,
  locale: 'zh',
  fallbackLocale: 'en',
  messages,
})
