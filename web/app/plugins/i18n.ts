import { en, fa } from "zod/v4/locales";
import z from "zod";

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook("i18n:localeSwitched", (options) => {
    localStorage.setItem("language", options.newLocale);
    if (options.newLocale === "en") {
      z.config(en());
    } else if (options.newLocale === "fa") {
      z.config(fa());
    }
  });
});
