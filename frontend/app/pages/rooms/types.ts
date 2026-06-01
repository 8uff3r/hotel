import type {
  PaginatedResponseModelsRoomStatus,
  PaginatedResponseRoomsRackRoom,
} from "~/utils/client";

export type RoomRack = NonNullable<PaginatedResponseRoomsRackRoom["data"]>[0];
export type RoomStatus = NonNullable<PaginatedResponseModelsRoomStatus["data"]>[0];
