hotel/app/pages/guests/index.vue ``` ```vue
<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Guests</h1>
      <UButton to="/guests/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        Add Guest
      </UButton>
    </div>

    <!-- Filters -->
    <UCard class="mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          placeholder="Search by name, email, or phone..."
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.idType"
          :items="idTypeOptions"
          placeholder="All ID Types"
          class="w-full sm:w-40"
          @change="fetchGuests"
        />
        <UButton variant="outline" @click="clearFilters"> Clear </UButton>
      </div>
    </UCard>

    <!-- Guests Table -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">Guest List</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} guests</span>
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
            <p v-if="row.original.email" class="text-sm text-gray-500">
              {{ row.original.email }}
            </p>
          </div>
        </template>

        <template #phone-cell="{ row }">
          {{ row.original.phone || "-" }}
        </template>

        <template #idInfo-cell="{ row }">
          <div v-if="row.original.idType">
            <UBadge variant="soft" size="sm">
              {{ formatIdType(row.original.idType) }}
            </UBadge>
            <p v-if="row.original.idNumber" class="text-sm text-gray-500">
              {{ row.original.idNumber }}
            </p>
          </div>
          <span v-else class="text-gray-400">-</span>
        </template>

        <template #location-cell="{ row }">
          <div v-if="row.original.city || row.original.country">
            <p>
              {{ row.original.city }}{{ row.original.city && row.original.country ? ", " : ""
              }}{{ row.original.country }}
            </p>
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
            Page {{ pagination.page }} of {{ pagination.totalPages }}
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
  firstName: string;
  lastName: string;
  email: string | null;
  phone: string | null;
  idType: string | null;
  idNumber: string | null;
  city: string | null;
  country: string | null;
}

const columns: TableColumn<GuestRow>[] = [
  { accessorKey: "id", header: "ID" },
  { accessorKey: "name", header: "Name" },
  { accessorKey: "phone", header: "Phone" },
  { accessorKey: "idInfo", header: "ID Info" },
  { accessorKey: "location", header: "Location" },
  { accessorKey: "actions", header: "Actions" },
];

const idTypeOptions = [
  { value: "all", label: "All ID Types" },
  { value: "passport", label: "Passport" },
  { value: "national_id", label: "National ID" },
  { value: "driver_license", label: "Driver License" },
  { value: "other", label: "Other" },
];

const guests = ref<GuestRow[]>([]);
const loading = ref(false);
const page = ref(1);

const filters = reactive({
  search: "",
  idType: "",
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
    if (filters.idType && filters.idType !== "all") params.append("idType", filters.idType);

    const response = await $fetch(`/api/guests?${params.toString()}`);
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
  filters.idType = "";
  pagination.page = 1;
  fetchGuests();
};

const formatIdType = (idType: string | null): string => {
  if (!idType) return "";
  const types: Record<string, string> = {
    passport: "Passport",
    national_id: "National ID",
    driver_license: "Driver License",
    other: "Other",
  };
  return types[idType] || idType;
};

onMounted(fetchGuests);
</script>
