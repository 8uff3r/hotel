import {
  getApiRoomsStatuses,
  getApiRoomsTypes,
  getApiRoomsAmenities,
} from "~/utils/client/sdk.gen";
import { getApiTravelAgencies } from "~/utils/client/sdk.gen";
import { getApiCommonCountries } from "~/utils/client/sdk.gen";

export interface RackFiltersState {
  roomTypeId: number | null;
  nationalityId: number | null;
  agencyId: number | null;
  entryDateFrom: string;
  entryDateTo: string;
  departureDateFrom: string;
  departureDateTo: string;
}

export function useRoomRackData() {
  const { data: allRooms, isLoading: roomsLoading } = useQuery({
    key: ["rooms", "rack"],
    query: async () => {
      const response = await getApiRoomsRack({});
      return response.data?.data || [];
    },
    placeholderData: [],
  });

  const { data: statuses } = useQuery({
    key: ["rooms", "status", "list"],
    query: async () => {
      const response = await getApiRoomsStatuses({});
      return response.data?.data || [];
    },
    placeholderData: [],
  });

  const { data: roomTypes } = useQuery({
    key: ["rooms", "types", "list"],
    query: async () => {
      const response = await getApiRoomsTypes({});
      return response.data?.data || [];
    },
    placeholderData: [],
  });

  const { data: amenities } = useQuery({
    key: ["rooms", "amenities", "list"],
    query: async () => {
      const response = await getApiRoomsAmenities({});
      return response.data?.data || [];
    },
    placeholderData: [],
  });

  const { data: agencies } = useQuery({
    key: ["travel-agencies", "list-all"],
    query: async () => {
      const response = await getApiTravelAgencies({ query: { limit: -1 } });
      return response.data?.data || [];
    },
    placeholderData: [],
  });

  const { data: countries } = useQuery({
    key: ["countries"],
    query: async () => {
      const response = await getApiCommonCountries({});
      return response.data?.data ?? [];
    },
    placeholderData: [],
  });

  return {
    allRooms,
    pending: roomsLoading,
    statuses,
    roomTypes,
    amenities,
    agencies,
    countries,
  };
}
