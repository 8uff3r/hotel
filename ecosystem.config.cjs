module.exports = {
  apps: [
    {
      name: "hotel-new",
      script: "./.output/server/index.mjs",
      interpreter: "node",
      env: {
        PORT: 6767,
      },
    },
  ],
};
