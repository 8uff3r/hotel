import { getApiRooms, getApiRoomsStatuses, getApiRoomsTypes, getApiRoomsAmenities } from '~/utils/client/sdk.gen'
import { getApiReservation } from '~/utils/client/sdk.gen'
import { getApiTravelAgencies } from '~/utils/client/sdk.gen'
import { getApiCommonCountries } from '~/utils/client/sdk.gen'
import type { Room } from '~/utils/client/types.gen'

export interface EnrichedRoom extends Room {
  currentReservation?: any
}

export interface RackFiltersState {
  roomTypeId: number | null
  nationalityId: number | null
  agencyId: number | null
  entryDateFrom: string
  entryDateTo: string
  departureDateFrom: string
  departureDateTo: string
}

export function useRoomRackData() {
  const { data: allRooms, isLoading: roomsLoading } = useQuery({
    key: ['rooms', 'list-all'],
    query: async () => {
      const response = await getApiRooms({ query: { limit: -1 } })
      return response.data?.data || []
    },
    placeholderData: [],
  })

  const { data: statuses } = useQuery({
    key: ['rooms', 'status', 'list'],
    query: async () => {
      const response = await getApiRoomsStatuses({})
      return response.data?.data || []
    },
    placeholderData: [],
  })

  const { data: roomTypes } = useQuery({
    key: ['rooms', 'types', 'list'],
    query: async () => {
      const response = await getApiRoomsTypes({})
      return response.data?.data || []
    },
    placeholderData: [],
  })

  const { data: amenities } = useQuery({
    key: ['rooms', 'amenities', 'list'],
    query: async () => {
      const response = await getApiRoomsAmenities({})
      return response.data?.data || []
    },
    placeholderData: [],
  })

  const { data: activeReservations } = useQuery({
    key: ['reservations', 'active'],
    query: async () => {
      const response = await getApiReservation({
        query: {
          limit: -1,
          filters: 'status:eq:checked_in',
        },
      })
      return response.data?.data || []
    },
    placeholderData: [],
  })

  const roomReservationMap = computed(() => {
    const map = new Map<number, any>()
    activeReservations.value?.forEach((res: any) => {
      res.rooms?.forEach((room: any) => {
        if (room.id) {
          map.set(room.id, res)
        }
      })
    })
    return map
  })

  const enrichedRooms = computed(() => {
    return (allRooms.value || []).map((room: any) => ({
      ...room,
      currentReservation: roomReservationMap.value.get(room.id!) || null,
    }))
  })

  const { data: agencies } = useQuery({
    key: ['travel-agencies', 'list-all'],
    query: async () => {
      const response = await getApiTravelAgencies({ query: { limit: -1 } })
      return response.data?.data || []
    },
    placeholderData: [],
  })

  const { data: countries } = useQuery({
    key: ['countries'],
    query: async () => {
      const response = await getApiCommonCountries({})
      return response.data?.data ?? []
    },
    placeholderData: [],
  })

  return {
    allRooms: enrichedRooms,
    pending: roomsLoading,
    statuses,
    roomTypes,
    amenities,
    agencies,
    countries,
    activeReservations,
  }
}
