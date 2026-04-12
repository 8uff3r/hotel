export default defineNuxtPlugin((nuxtApp) => {
  const api = $fetch.create({
    onRequest({ options }) {
      const { locale } = useI18n();
      options.headers.set("Accept-Language", locale.value);
    },
    async onResponseError({ response }) {
      if (response.status === 401) {
        await nuxtApp.runWithContext(() => navigateTo("/login"));
      }
    },
  });

  return {
    provide: {
      api,
    },
  };
});
