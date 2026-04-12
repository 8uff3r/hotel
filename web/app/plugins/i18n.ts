export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook("i18n:localeSwitched", (options) => {
    console.log("switched");
    localStorage.setItem("language", options.newLocale);
  });
});
