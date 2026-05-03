<template>
  <UCard>
    <template #header>
      <div class="flex items-center justify-between">
        <span class="text-lg font-semibold">{{ t("restaurant.inventory") }}</span>
        <UButton v-if="canCreate" color="primary" @click="openCreateModal">
          <UIcon name="i-lucide-plus" class="mr-2" />
          {{ t("restaurant.addItem") }}
        </UButton>
      </div>
    </template>

    <!-- Filters -->
    <div class="mb-4 flex flex-wrap items-center gap-4">
      <UInput
        v-model="filters.search"
        :placeholder="t('restaurant.searchPlaceholder')"
        icon="i-lucide-search"
        class="w-full sm:w-64"
        @input="debouncedSearch"
      />
      <USelect
        v-model="filters.category"
        :items="categoryOptions"
        :placeholder="t('restaurant.allCategories')"
        class="w-full sm:w-40"
      />
      <UButton variant="outline" @click="clearFilters">{{ t("common.clear") }}</UButton>
    </div>

    <!-- Items Table -->
    <UTable :data="items" :columns="columns" :loading="pending" striped>
      <template #name-cell="{ row }">
        <NuxtLink
          :to="`/restaurant/inventory/${row.original.id}`"
          class="font-medium text-primary hover:underline"
        >
          {{ row.original.name }}
        </NuxtLink>
      </template>

      <template #category-cell="{ row }">
        <UBadge variant="soft" color="primary">{{ row.original.category }}</UBadge>
      </template>

      <template #quantity-cell="{ row }">
        <span :class="{ 'text-red-500': row.original.quantity <= row.original.reorderLevel }">
          {{ row.original.quantity }} {{ row.original.unit }}
        </span>
      </template>

      <template #unitCost-cell="{ row }"> ${{ row.original.unitCost?.toFixed(2) }} </template>

      <template #isActive-cell="{ row }">
        <UBadge :color="row.original.isActive ? 'success' : 'error'" variant="soft">
          {{ row.original.isActive ? t("restaurant.active") : t("restaurant.inactive") }}
        </UBadge>
      </template>

      <template #actions-cell="{ row }">
        <div class="flex items-center gap-2">
          <UButton variant="ghost" size="sm" @click="openEditModal(row.original)">
            <UIcon name="i-lucide-pencil" class="h-4 w-4" />
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
          {{ t("pagination.pageOf", { page: pagination.page, totalPages: pagination.totalPages }) }}
        </span>
        <UPagination v-model="page" :page-count="pagination.limit" :total="pagination.total" />
      </div>
    </template>
  </UCard>

  <!-- Create/Edit Modal -->
  <UModal v-model="modalOpen">
    <template #header>
      <h2 class="text-lg font-semibold">
        {{ isEditing ? t("restaurant.editItem") : t("restaurant.addItem") }}
      </h2>
    </template>

    <template #body>
      <UForm :state="form" class="space-y-4">
        <UFormGroup :label="t('restaurant.itemName')" required>
          <UInput v-model="form.name" />
        </UFormGroup>

        <UFormGroup :label="t('restaurant.category')" required>
          <USelect v-model="form.category" :items="categorySelectOptions" />
        </UFormGroup>

        <UFormGroup :label="t('restaurant.unit')" required>
          <USelect v-model="form.unit" :items="unitOptions" />
        </UFormGroup>

        <div class="grid grid-cols-2 gap-4">
          <UFormGroup :label="t('restaurant.quantity')" required>
            <UInput v-model="form.quantity" type="number" min="0" />
          </UFormGroup>

          <UFormGroup :label="t('restaurant.unitCost')" required>
            <UInput v-model="form.unitCost" type="number" min="0" step="0.01" />
          </UFormGroup>
        </div>

        <UFormGroup :label="t('restaurant.reorderLevel')">
          <UInput v-model="form.reorderLevel" type="number" min="0" />
        </UFormGroup>

        <UFormGroup :label="t('restaurant.description')">
          <UTextarea v-model="form.description" :rows="3" />
        </UFormGroup>

        <UFormGroup :label="t('common.status')">
          <UCheckbox v-model="form.isActive">{{ t("restaurant.active") }}</UCheckbox>
        </UFormGroup>
      </UForm>
    </template>

    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton variant="outline" @click="modalOpen = false">{{ t("actions.cancel") }}</UButton>
        <UButton color="primary" :loading="saving" @click="saveItem">
          {{ t("actions.saveChanges") }}
        </UButton>
      </div>
    </template>
  </UModal>

  <!-- Delete Modal -->
  <UModal v-model="deleteModalOpen">
    <template #header>
      <h2 class="text-lg font-semibold">{{ t("actions.confirmDelete") }}</h2>
    </template>
    <template #body>
      <p>{{ t("restaurant.confirmDelete", { name: selectedItem?.name }) }}</p>
    </template>
    <template #footer>
      <div class="flex justify-end gap-3">
        <UButton variant="outline" @click="deleteModalOpen = false">{{
          t("actions.cancel")
        }}</UButton>
        <UButton color="error" :loading="deleting" @click="deleteItem">{{
          t("actions.delete")
        }}</UButton>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import type { PaginatedResponseModelsInventoryItem, InventoryItem } from "~/utils/client";

definePageMeta({
  requiresPermission: PERMISSIONS.restaurant.restaurantInventory.read,
});

const canCreate = useAuthStore().can(PERMISSIONS.restaurant.restaurantInventory.create);

const { t } = useI18n();
const columns = computed<TableColumn<InventoryItem>[]>(() => [
  { accessorKey: "name", header: t("restaurant.columns.name") },
  { accessorKey: "category", header: t("restaurant.columns.category") },
  { accessorKey: "quantity", header: t("restaurant.columns.quantity") },
  { accessorKey: "unitCost", header: t("restaurant.columns.price") },
  { accessorKey: "isActive", header: t("restaurant.columns.status") },
  { accessorKey: "actions", header: t("restaurant.columns.actions") },
]);

const categoryOptions = computed(() => [
  { value: "", label: t("restaurant.allCategories") },
  { value: "food", label: t("restaurant.categories.food") },
  { value: "beverage", label: t("restaurant.categories.beverage") },
  { value: "dessert", label: t("restaurant.categories.dessert") },
  { value: "other", label: t("restaurant.categories.other") },
]);

const categorySelectOptions = computed(() => [
  { value: "food", label: t("restaurant.categories.food") },
  { value: "beverage", label: t("restaurant.categories.beverage") },
  { value: "dessert", label: t("restaurant.categories.dessert") },
  { value: "other", label: t("restaurant.categories.other") },
]);

const unitOptions = computed(() => [
  { value: "piece", label: t("restaurant.units.piece") },
  { value: "gram", label: t("restaurant.units.gram") },
  { value: "kilogram", label: t("restaurant.units.kilogram") },
  { value: "liter", label: t("restaurant.units.liter") },
  { value: "pack", label: t("restaurant.units.pack") },
]);

const pagination = reactive({ page: 1, limit: 20, total: 0, totalPages: 0 });
const page = computed({
  get: () => pagination.page,
  set: (val) => {
    pagination.page = val;
    fetchItems();
  },
});

const filters = reactive({ search: "", category: "" });
let searchTimeout: ReturnType<typeof setTimeout> | null = null;

const debouncedSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    pagination.page = 1;
    fetchItems();
  }, 300);
};

const clearFilters = () => {
  filters.search = "";
  filters.category = "";
  pagination.page = 1;
  fetchItems();
};

const {
  data: itemsData,
  pending,
  refresh: fetchItems,
} = useAsyncData(
  "inventory-items",
  async () => {
    const params: any = { page: pagination.page, limit: pagination.limit };
    if (filters.search) params.search = filters.search;
    if (filters.category) params.category = filters.category;

    const response = await $fetch<PaginatedResponseModelsInventoryItem>(
      "/api/restaurant/inventory",
      { query: params }
    );
    pagination.total = response.total ?? 0;
    pagination.totalPages = response.totalPages ?? 0;
    return response.data;
  },
  { watch: [() => pagination.page] }
);

const items = computed(() => itemsData.value ?? []);

const modalOpen = ref(false);
const isEditing = ref(false);
const saving = ref(false);
const form = reactive({
  id: 0,
  name: "",
  category: "food",
  unit: "piece",
  quantity: 0,
  unitCost: 0,
  reorderLevel: 0,
  description: "",
  isActive: true,
});

const openCreateModal = () => {
  isEditing.value = false;
  form.id = 0;
  form.name = "";
  form.category = "food";
  form.unit = "piece";
  form.quantity = 0;
  form.unitCost = 0;
  form.reorderLevel = 0;
  form.description = "";
  form.isActive = true;
  modalOpen.value = true;
};

const openEditModal = (item: InventoryItem) => {
  isEditing.value = true;
  Object.assign(form, item);
  modalOpen.value = true;
};

const saveItem = async () => {
  saving.value = true;
  try {
    if (isEditing.value) {
      await $fetch(`/api/restaurant/inventory/${form.id}`, { method: "PUT", body: form });
    } else {
      await $fetch("/api/restaurant/inventory", { method: "POST", body: form });
    }
    modalOpen.value = false;
    await fetchItems();
  } catch (error) {
    console.error("Failed to save item:", error);
  } finally {
    saving.value = false;
  }
};

const deleteModalOpen = ref(false);
const deleting = ref(false);
const selectedItem = ref<InventoryItem | null>(null);

const confirmDelete = (item: InventoryItem) => {
  selectedItem.value = item;
  deleteModalOpen.value = true;
};

const deleteItem = async () => {
  if (!selectedItem.value) return;
  deleting.value = true;
  try {
    await $fetch(`/api/restaurant/inventory/${selectedItem.value.id}`, { method: "DELETE" });
    deleteModalOpen.value = false;
    await fetchItems();
  } catch (error) {
    console.error("Failed to delete item:", error);
  } finally {
    deleting.value = false;
  }
};
</script>
