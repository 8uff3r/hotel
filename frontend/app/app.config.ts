export default defineAppConfig({
  ui: {
    input: {
      slots: {
        root: "w-full",
        base: "w-full rounded-md border-0 appearance-none placeholder:text-dimmed focus:outline-none disabled:cursor-not-allowed disabled:opacity-75",
      },
    },
    select: {
      slots: {
        base: "w-full",
      },
    },
    selectMenu: {
      slots: {
        base: "w-full",
      },
    },
  },
});
