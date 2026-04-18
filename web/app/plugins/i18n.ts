export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook("i18n:localeSwitched", (options) => {
    localStorage.setItem("language", options.newLocale);
  });
});
