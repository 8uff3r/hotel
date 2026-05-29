import { fa } from "zod/v4/locales";
import z from "zod";

export default defineNuxtPlugin(() => {
  z.config({
    ...fa(),
  });
});
