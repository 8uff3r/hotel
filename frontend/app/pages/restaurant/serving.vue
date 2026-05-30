<template>
  <div class="space-y-6">
    <!-- New Bill Section -->
    <UCard>
      <template #header>
        <span class="text-lg font-semibold">{{ t("restaurant.newBill") }}</span>
      </template>

      <UForm :state="billForm" class="space-y-4">
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <UFormGroup :label="t('restaurant.selectGuest')" required>
            <USelect
              v-model="billForm.guestId"
              :items="guestOptions"
              :placeholder="t('restaurant.selectGuest')"
              searchable
            />
          </UFormGroup>

          <UFormGroup :label="t('restaurant.selectRoom')">
            <USelect
              v-model="billForm.roomId"
              :items="roomOptions"
              :placeholder="t('restaurant.selectRoom')"
              searchable
            />
          </UFormGroup>
        </div>

        <UFormGroup>
          <UCheckbox v-model="billForm.isExternal">
            <div class="flex items-center gap-2">
              <UIcon name="i-lucide-shop" class="h-4 w-4" />
              {{ t("restaurant.markAsExternal") }}
            </div>
          </UCheckbox>
        </UFormGroup>

        <UFormGroup v-if="billForm.isExternal" :label="t('restaurant.externalRestaurant')">
          <UInput
            v-model="billForm.externalRestaurant"
            :placeholder="t('restaurant.externalRestaurant')"
          />
        </UFormGroup>
      </UForm>
    </UCard>

    <!-- Active Bills -->
    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <span class="text-lg font-semibold">{{ t("restaurant.transactions") }}</span>
          <UButton color="primary" @click="createBill">
            <UIcon name="i-lucide-plus" class="mr-2" />
            {{ t("restaurant.createBill") }}
          </UButton>
        </div>
      </template>

      <div class="mb-4 flex flex-wrap items-center gap-4">
        <UInput
          v-model="filters.search"
          :placeholder="t('restaurant.searchPlaceholder')"
          icon="i-lucide-search"
          class="w-full sm:w-64"
          @input="debouncedSearch"
        />
        <USelect
          v-model="filters.status"
          :items="statusOptions"
          :placeholder="t('common.all_statuses')"
          class="w-full sm:w-40"
        />
        <UButton variant="outline" @click="clearFilters">{{ t("common.clear") }}</UButton>
      </div>

      <div v-if="bills.length === 0 && !pending" class="py-8 text-center text-gray-500">
        {{ t("restaurant.noBills") }}
      </div>

      <div v-else class="grid grid-cols-1 gap-4 lg:grid-cols-2">
        <UCard
          v-for="bill in bills"
          :key="bill.id"
          :class="{ 'ring-2 ring-primary-500': selectedBill?.id === bill.id }"
          @click="selectBill(bill)"
        >
          <template #header>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="font-semibold">#{{ bill.id }}</span>
                <UBadge v-if="bill.isExternal" color="warning" variant="soft">
                  <UIcon name="i-lucide-shop" class="mr-1 h-3 w-3" />
                  {{ t("restaurant.external") }}
                </UBadge>
                <UBadge v-else color="success" variant="soft">
                  <UIcon name="i-lucide-building" class="mr-1 h-3 w-3" />
                  {{ t("restaurant.internal") }}
                </UBadge>
              </div>
              <UBadge
                :style="{
                  backgroundColor: bill.status?.colorHex,
                }"
                variant="soft"
              >
                {{ bill.status?.label }}
              </UBadge>
            </div>
          </template>

          <div class="space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">{{ t("forms.firstName") }}</span>
              <span class="font-medium">{{ bill.guestId ? `#${bill.guestId}` : "-" }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">{{ t("restaurant.total") }}</span>
              <span class="font-semibold">${{ bill.totalAmount?.toFixed(2) }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-500">{{ t("common.created") }}</span>
              <span>{{ formatDate(bill.billDate) }}</span>
            </div>
          </div>

          <template #footer>
            <div class="flex items-center justify-between">
              <UButton
                v-if="!bill.settled"
                variant="ghost"
                size="sm"
                color="success"
                @click.stop="settleBill(bill)"
              >
                <UIcon name="i-lucide-check-circle" class="mr-1 h-4 w-4" />
                {{ t("restaurant.settle") }}
              </UButton>
              <UButton variant="ghost" size="sm" @click.stop="viewTransactions(bill)">
                <UIcon name="i-lucide-list" class="mr-1 h-4 w-4" />
                {{ t("restaurant.transactions") }}
              </UButton>
            </div>
          </template>
        </UCard>
      </div>

      <template #footer>
        <div class="flex items-center justify-between">
          <span class="text-sm text-gray-500">
            {{
              t("pagination.pageOf", { page: pagination.page, totalPages: pagination.totalPages })
            }}
          </span>
          <UPagination v-model="page" :page-count="pagination.limit" :total="pagination.total" />
        </div>
      </template>
    </UCard>

    <!-- Add Meal Modal -->
    <UModal v-model="mealModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">{{ t("restaurant.addMeal") }}</h2>
      </template>

      <template #body>
        <UForm :state="mealForm" class="space-y-4">
          <UFormGroup :label="t('restaurant.item')" required>
            <USelect
              v-model="mealForm.inventoryItemId"
              :items="inventoryOptions"
              :placeholder="t('restaurant.item')"
              searchable
              label-key="name"
              value-key="id"
              @update:model-value="onItemSelect"
            />
          </UFormGroup>

          <div class="grid grid-cols-2 gap-4">
            <UFormGroup :label="t('restaurant.quantity')" required>
              <UInput
                v-model="mealForm.quantity"
                type="number"
                min="1"
                @update:model-value="calculateTotal"
              />
            </UFormGroup>

            <UFormGroup :label="t('restaurant.price')">
              <UInput v-model="mealForm.unitPrice" type="number" min="0" step="0.01" readonly />
            </UFormGroup>
          </div>

          <div class="flex justify-between border-t pt-4">
            <span class="font-semibold">{{ t("restaurant.total") }}</span>
            <span class="font-bold">${{ mealTotal.toFixed(2) }}</span>
          </div>

          <UFormGroup :label="t('forms.notes')">
            <UTextarea v-model="mealForm.notes" :rows="2" />
          </UFormGroup>
        </UForm>
      </template>

      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton variant="outline" @click="mealModalOpen = false">{{
            t("actions.cancel")
          }}</UButton>
          <UButton color="primary" :loading="savingMeal" @click="addMeal">{{
            t("restaurant.addMeal")
          }}</UButton>
        </div>
      </template>
    </UModal>

    <!-- Transactions Modal -->
    <UModal v-model="transactionsModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold">
          {{ t("restaurant.transactions") }} - #{{ selectedBill?.id }}
        </h2>
      </template>

      <template #body>
        <div v-if="selectedBill" class="space-y-4">
          <div class="mb-4 flex justify-end">
            <UButton color="primary" :disabled="selectedBill.settled" @click="openMealModal">
              <UIcon name="i-lucide-plus" class="mr-2" />
              {{ t("restaurant.addMeal") }}
            </UButton>
          </div>

          <div v-if="transactions.length === 0" class="py-8 text-center text-gray-500">
            {{ t("restaurant.noItems") }}
          </div>

          <UTable v-else :data="transactions" :columns="transactionColumns" striped>
            <template #itemName-cell="{ row }">
              <span>{{ row.original.itemName }}</span>
              <UBadge v-if="row.original.isExternal" color="warning" variant="soft" class="mr-2">
                {{ t("restaurant.external") }}
              </UBadge>
            </template>

            <template #quantity-cell="{ row }">
              {{ row.original.quantity }}
            </template>

            <template #totalPrice-cell="{ row }">
              ${{ row.original.totalPrice?.toFixed(2) }}
            </template>

            <template #actions-cell="{ row }">
              <UButton
                v-if="!selectedBill?.settled"
                variant="ghost"
                size="sm"
                color="error"
                @click="deleteTransaction(row.original)"
              >
                <UIcon name="i-lucide-trash-2" class="h-4 w-4" />
              </UButton>
            </template>
          </UTable>

          <div class="flex justify-between border-t pt-4">
            <span class="font-semibold">{{ t("restaurant.total") }}</span>
            <span class="font-bold">${{ selectedBill.totalAmount?.toFixed(2) }}</span>
          </div>
        </div>
      </template>

      <template #footer>
        <div class="flex justify-end">
          <UButton variant="outline" @click="transactionsModalOpen = false">{{
            t("actions.back")
          }}</UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type {
  PaginatedResponseModelsRestaurantBill,
  RestaurantBill,
  PaginatedResponseModelsInventoryItem,
  InventoryItem,
  PaginatedResponseModelsMealTransaction,
  MealTransaction,
  PaginatedResponseModelsRoom,
  Room,
  PaginatedResponseModelsGuest,
  Guest,
} from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.restaurant.restaurantServing.read,
});

const { t } = useI18n();

const transactionColumns = computed<TableColumn<MealTransaction>[]>(() => [
  { accessorKey: "itemName", header: t("restaurant.item") },
  { accessorKey: "quantity", header: t("restaurant.quantity") },
  { accessorKey: "unitPrice", header: t("restaurant.price") },
  { accessorKey: "totalPrice", header: t("restaurant.total") },
  { accessorKey: "actions", header: t("restaurant.columns.actions") },
]);

const statusOptions = computed(() => [
  { value: "", label: t("common.all_statuses") },
  { value: "open", label: t("common.status") },
  { value: "settled", label: t("accounting.settle") },
  { value: "cancelled", label: t("parking.status_cancelled") },
]);

const pagination = reactive({ page: 1, limit: 20, total: 0, totalPages: 0 });
const page = computed({
  get: () => pagination.page,
  set: (val) => {
    pagination.page = val;
    fetchBills();
  },
});

const filters = reactive({ search: "", status: "" });
let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pagination.page = 1;
    fetchBills();
  }, 300);
};

const clearFilters = () => {
  filters.search = "";
  filters.status = "";
  pagination.page = 1;
  fetchBills();
};

const {
  data: billsData,
  pending,
  refresh: fetchBills,
} = useAsyncData(
  "restaurant-bills-serving",
  async () => {
    const params: any = { page: pagination.page, limit: pagination.limit };
    if (filters.search) params.search = filters.search;
    if (filters.status) params.status = filters.status;
    else params.status = "open";

    const response = await $fetch<PaginatedResponseModelsRestaurantBill>("/api/restaurant/bills", {
      query: params,
    });
    pagination.total = response.data?.data?.total ?? 0;
    pagination.totalPages = response.data?.data?.totalPages ?? 0;
    return response.data?.data;
  },
  { watch: [() => pagination.page], immediate: true }
);

const bills = computed(() => billsData.value ?? []);

const { data: guestsData } = useAsyncData("guests-select", () => getApiGuests({}), {
  transform: (response) => response.data?.data,
});

const guestOptions = computed(() =>
  (guestsData.value ?? []).map((g) => ({
    value: g.id,
    label: `${g.firstName} ${g.lastName}`,
  }))
);

const { data: roomsData } = useAsyncData("rooms-select", () => getApiRooms({}), {
  transform: (response) => response.data?.data,
});

const roomOptions = computed(() =>
  (roomsData.value ?? []).map((r) => ({
    value: r.id,
    label: r.roomNumber,
  }))
);

const { data: inventoryData } = useAsyncData(
  "inventory-select",
  () => getApiRestaurantInventory({}),
  {
    transform: (response) => response.data?.data,
  }
);

const inventoryOptions = computed(() =>
  (inventoryData.value ?? [])
    .map((i) => ({
      id: i.id,
      name: i.name,
      unitCost: i.unitCost,
    }))
    .filter((i) => i.id)
);

const billForm = reactive({
  guestId: undefined as number | undefined,
  roomId: undefined as number | undefined,
  isExternal: false,
  externalRestaurant: "",
});

const createBill = async () => {
  try {
    await $fetch("/api/restaurant/bills", {
      method: "POST",
      body: {
        billDate: new Date().toISOString(),
        guestId: billForm.guestId,
        roomId: billForm.roomId,
        isExternal: billForm.isExternal,
        externalRestaurant: billForm.externalRestaurant,
      },
    });
    billForm.guestId = undefined;
    billForm.roomId = undefined;
    billForm.isExternal = false;
    billForm.externalRestaurant = "";
    await fetchBills();
  } catch (error) {
    console.error("Failed to create bill:", error);
  }
};

const selectedBill = ref<RestaurantBill | null>(null);
const selectBill = (bill: RestaurantBill) => {
  selectedBill.value = bill;
};

const settleBill = async (bill: RestaurantBill) => {
  try {
    await $fetch(`/api/restaurant/bills/${bill.id}/settle`, { method: "POST" });
    await fetchBills();
  } catch (error) {
    console.error("Failed to settle bill:", error);
  }
};

const mealModalOpen = ref(false);
const savingMeal = ref(false);
const mealForm = reactive({
  inventoryItemId: undefined as number | undefined,
  itemName: "",
  quantity: 1,
  unitPrice: 0,
  notes: "",
});

const mealTotal = computed(() => mealForm.quantity * mealForm.unitPrice);

const openMealModal = () => {
  if (!selectedBill.value) return;
  mealForm.inventoryItemId = undefined;
  mealForm.itemName = "";
  mealForm.quantity = 1;
  mealForm.unitPrice = 0;
  mealForm.notes = "";
  mealModalOpen.value = true;
};

const onItemSelect = (item: any) => {
  if (item) {
    mealForm.unitPrice = item.unitCost || 0;
    mealForm.itemName = item.name;
  }
};

const calculateTotal = () => {
  // Auto-calculate
};

const addMeal = async () => {
  if (!selectedBill.value || !mealForm.inventoryItemId) return;

  savingMeal.value = true;
  try {
    await $fetch("/api/restaurant/transactions", {
      method: "POST",
      body: {
        billId: selectedBill.value.id,
        inventoryItemId: mealForm.inventoryItemId,
        itemName: mealForm.itemName,
        quantity: mealForm.quantity,
        unitPrice: mealForm.unitPrice,
        totalPrice: mealTotal.value,
        isExternal: selectedBill.value.isExternal,
        notes: mealForm.notes,
      },
    });
    mealModalOpen.value = false;
    await fetchTransactions();
    await fetchBills();
  } catch (error) {
    console.error("Failed to add meal:", error);
  } finally {
    savingMeal.value = false;
  }
};

const transactionsModalOpen = ref(false);
const transactions = ref<MealTransaction[]>([]);

const viewTransactions = async (bill: RestaurantBill) => {
  selectedBill.value = bill;
  await fetchTransactions();
  transactionsModalOpen.value = true;
};

const fetchTransactions = async () => {
  if (!selectedBill.value) return;
  const response = await $fetch<PaginatedResponseModelsMealTransaction>(
    "/api/restaurant/transactions",
    {
      query: { billId: selectedBill.value.id },
    }
  );
  transactions.value = response.data?.data ?? [];
};

const deleteTransaction = async (transaction: MealTransaction) => {
  if (!transaction) return;
  try {
    await $fetch(`/api/restaurant/transactions/${transaction.id}`, { method: "DELETE" });
    await fetchTransactions();
    await fetchBills();
  } catch (error) {
    console.error("Failed to delete transaction:", error);
  }
};

const formatDate = (date: string | Date | undefined) => {
  if (!date) return "-";
  return new Date(date).toLocaleDateString();
};

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    open: "info",
    settled: "success",
    cancelled: "error",
  };
  return colors[status] || "neutral";
};
</script>
