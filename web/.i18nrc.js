const { defineConfig } = require('@lobehub/i18n-cli');

module.exports = defineConfig({
  entry: 'src/i18n/locales/zh_CN.json',
  entryLocale: 'zh_CN',
  output: 'src/i18n/locales',
  outputLocales: ['zh_TW', 'en_US', 'zh_HK'],
  temperature: 0,
  modelName: 'gpt-4o-mini',
  splitToken: 2048,
  experimental: {
    jsonMode: true
  }
});
