import { zGuestWithStayRequest, zPostApiGuestsWithReservationBody } from "~/utils/client/zod.gen";
import type z from "zod";

export const companionSchema = zPostApiGuestsWithReservationBody.shape.companions.unwrap().element;
export type Companion = NonNullable<z.output<typeof companionSchema>>;

export const createSchema = zGuestWithStayRequest;
export type CreateRequest = z.output<typeof createSchema>;
