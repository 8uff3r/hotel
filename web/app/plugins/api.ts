export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.hook("openFetch:onRequest", (ctx) => {
    ctx.options.headers.set("Accept-Language", localStorage.getItem("language") ?? "fa");
  });
  nuxtApp.hook("openFetch:onRequestError", async (ctx) => {
    if (ctx.response?.status === 401) {
      await nuxtApp.runWithContext(() => navigateTo("/login"));
    }
  });
});
