<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton to="/parking" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />
          Back to Parking
        </UButton>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Vehicles</h1>
      </div>
      <UButton to="/parking/vehicles/create" color="primary">
        <UIcon name="i-lucide-plus" class="mr-2" />
        Add Vehicle
      </UButton>
    </div>

    <UCard class="mb-4">
      <div class="flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          placeholder="Search by license plate..."
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.vehicleType"
          :items="typeOptions"
          placeholder="All Types"
          class="w-full sm:w-40"
          @change="fetchVehicles"
        />
        <UButton variant="outline" @click="clearFilters"> Clear </UButton>
      </div>
    </UCard>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">Registered Vehicles</span>
          <span class="text-sm text-gray-500">{{ pagination.total }} vehicles</span>
        </div>
      </template>

      <UTable :data="vehicles" :columns="columns" :loading="loading" striped>
        <template #licensePlate-cell="{ row }">
          <NuxtLink
            :to="`/parking/vehicles/${row.original.id}`"
            class="font-medium text-primary hover:underline"
          >
            {{ row.original.licensePlate }}
          </NuxtLink>
        </template>

        <template #guestId-cell="{ row }">
          {{ getGuestName(row.original.guestId) }}
        </template>

        <template #vehicleType-cell="{ row }">
          <UBadge :color="getTypeColor(row.original.vehicleType) as any" variant="soft">
            {{ row.original.vehicleType }}
          </UBadge>
        </template>

        <template #isRegistered-cell="{ row }">
          <UBadge :color="row.original.isRegistered ? 'success' : 'warning'" variant="soft">
            {{ row.original.isRegistered ? "Registered" : "Guest" }}
          </UBadge>
        </template>

        <template #actions-cell="{ row }">
          <div class="flex items-center gap-2">
            <UButton variant="ghost" size="sm" :to="`/parking/vehicles/${row.original.id}`">
              <UIcon name="i-lucide-eye" class="h-4 w-4" />
            </UButton>
            <UButton variant="ghost" size="sm" color="error" @click="confirmDelete(row.original)">
              <UIcon name="i-lucide-trash-2" class="h-4 w-4" />
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
            @change="fetchVehicles"
          />
        </div>
      </template>
    </UCard>

    <UModal v-model="deleteModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">Confirm Delete</h2>
      </template>
      <template #body>
        <p>Are you sure you want to delete vehicle "{{ selectedVehicle?.licensePlate }}"?</p>
      </template>
      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="deleteModalOpen = false">Cancel</UButton>
          <UButton color="error" :loading="deleting" @click="deleteVehicle">Delete</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";

definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

interface Vehicle {
  id: number;
  licensePlate: string;
  guestId: number | null;
  vehicleType: string;
  make: string | null;
  model: string | null;
  color: string | null;
  isRegistered: number;
}

const columns: TableColumn<Vehicle>[] = [
  { accessorKey: "licensePlate", header: "License Plate" },
  { accessorKey: "guestId", header: "Guest" },
  { accessorKey: "vehicleType", header: "Type" },
  { accessorKey: "make", header: "Make" },
  { accessorKey: "model", header: "Model" },
  { accessorKey: "color", header: "Color" },
  { accessorKey: "isRegistered", header: "Status" },
  { accessorKey: "actions", header: "Actions" },
];

const vehicles = ref<Vehicle[]>([]);
const guests = ref<any[]>([]);
const loading = ref(false);
const deleting = ref(false);
const deleteModalOpen = ref(false);
const selectedVehicle = ref<Vehicle | null>(null);
const page = ref(1);

const filters = reactive({
  search: "",
  vehicleType: "",
});

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  totalPages: 0,
});

const typeOptions = [
  { value: "all", label: "All Types" },
  { value: "car", label: "Car" },
  { value: "motorcycle", label: "Motorcycle" },
  { value: "truck", label: "Truck" },
  { value: "van", label: "Van" },
  { value: "other", label: "Other" },
];

let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pagination.page = 1;
    fetchVehicles();
  }, 300);
};

const fetchGuests = async () => {
  try {
    const res = await $fetch("/api/guests");
    guests.value = res.data;
  } catch (error) {
    console.error("Failed to fetch guests:", error);
  }
};

const getGuestName = (guestId: number | null) => {
  if (!guestId) return "-";
  const guest = guests.value.find((g) => g.id === guestId);
  return guest ? `${guest.firstName} ${guest.lastName}` : "-";
};

const fetchVehicles = async () => {
  loading.value = true;
  try {
    const params = new URLSearchParams();
    params.append("page", pagination.page.toString());
    params.append("limit", pagination.limit.toString());

    if (filters.search) params.append("search", filters.search);
    if (filters.vehicleType && filters.vehicleType !== "all")
      params.append("vehicleType", filters.vehicleType);

    const response = await $fetch(`/api/parking/vehicles?${params.toString()}`);
    vehicles.value = response.data;
    pagination.total = response.pagination.total ?? 0;
    pagination.totalPages = response.pagination.totalPages ?? 0;
  } catch (error) {
    console.error("Failed to fetch vehicles:", error);
  } finally {
    loading.value = false;
  }
};

const clearFilters = () => {
  filters.search = "";
  filters.vehicleType = "";
  pagination.page = 1;
  fetchVehicles();
};

const confirmDelete = (vehicle: Vehicle) => {
  selectedVehicle.value = vehicle;
  deleteModalOpen.value = true;
};

const deleteVehicle = async () => {
  if (!selectedVehicle.value) return;

  deleting.value = true;
  try {
    await $fetch(`/api/parking/vehicles/${selectedVehicle.value.id}`, { method: "DELETE" });
    deleteModalOpen.value = false;
    await fetchVehicles();
  } catch (error) {
    console.error("Failed to delete vehicle:", error);
  } finally {
    deleting.value = false;
  }
};

const getTypeColor = (type: string) => {
  const colors: Record<string, string> = {
    car: "primary",
    motorcycle: "info",
    truck: "warning",
    van: "success",
    other: "neutral",
  };
  return colors[type] || "neutral";
};

onMounted(() => {
  fetchGuests();
  fetchVehicles();
});
</script>
