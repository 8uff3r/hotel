<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("guests.title") }}</h1>
      <UButton to="/guests/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        {{ t("guests.addGuest") }}
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('guests.searchPlaceholder')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <UButton variant="outline" @click="clearFilters"> {{ t("actions.clear") }} </UButton>
      </div>
    </UCard>

    <!-- Guests Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("guests.list") }}</span>
          <span class="text-sm text-gray-500">{{
            t("guests.count", { count: pagination.total })
          }}</span>
        </div>
      </template>

      <UTable :data="guests" :columns="columns" :loading="loading" striped>
        <template #id-cell="{ row }">
          <NuxtLink
            :to="`/guests/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            #{{ row.original.id }}
          </NuxtLink>
        </template>

        <template #name-cell="{ row }">
          <div>
            <p class="font-medium">{{ row.original.firstName }} {{ row.original.lastName }}</p>
          </div>
        </template>

        <template #room-cell="{ row }">
          {{ row.original.roomNumber || "-" }}
        </template>

        <template #phone-cell="{ row }">
          {{ row.original.phone || "-" }}
        </template>

        <template #reservation-cell="{ row }">
          <div v-if="row.original.reservationCode || row.original.roomType">
            <p>{{ row.original.reservationCode || "-" }}</p>
            <p class="text-sm text-gray-500">{{ row.original.roomType || "-" }}</p>
          </div>
          <span v-else class="text-gray-400">-</span>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/guests/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton variant="ghost" size="sm" :to="`/guests/${row.original.id}/edit`">
              <UIcon name="i-lucide-pencil" class="h-4 w-4" />
            </UButton>
          </div>
        </template>
      </UTable>

      <template #footer>
        <div class="flex items-center justify-between">
          <span class="text-sm text-gray-500">
            {{
              t("pagination.pageOf", { page: pagination.page, totalPages: pagination.totalPages })
            }}
          </span>
          <UPagination
            v-model="page"
            :page-count="pagination.limit"
            :total="pagination.total"
            @change="fetchGuests"
          />
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";

definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

interface GuestRow {
  id: number;
  roomNumber: string | null;
  roomType: string | null;
  reservationCode: string | null;
  firstName: string;
  lastName: string;
  phone: string | null;
}

const { t } = useI18n();
const columns = computed<TableColumn<GuestRow>[]>(() => [
  { accessorKey: "id", header: t("guests.columns.id") },
  { accessorKey: "name", header: t("guests.columns.name") },
  { accessorKey: "room", header: "Room" },
  { accessorKey: "phone", header: t("guests.columns.phone") },
  { accessorKey: "reservation", header: "Reservation" },
  { accessorKey: "actions", header: t("guests.columns.actions") },
]);

const guests = ref<GuestRow[]>([]);
const loading = ref(false);
const page = ref(1);

const filters = reactive({
  search: "",
});

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  totalPages: 0,
});

let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pagination.page = 1;
    fetchGuests();
  }, 300);
};

const fetchGuests = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    params.append("page", pagination.page.toString());
    params.append("limit", pagination.limit.toString());

    if (filters.search) params.append("search", filters.search);
    const response = await $api("/api/guests/", {
      query: {
        page: pagination.page,
        limit: pagination.limit,
      },
    });
    guests.value = response.data;
    pagination.total = response.pagination.total;
    pagination.totalPages = response.pagination.totalPages;
  } catch (error) {
    console.error("Failed to fetch guests:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  pagination.page = 1;
  fetchGuests();
};

onMounted(fetchGuests);
</script>
