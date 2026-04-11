hotel/app/pages/rooms/create.vue#L1-120
<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/rooms" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        {{ t("actions.backToRooms") }}
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("rooms.addNew") }}</h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- Room Number -->
          <UFormField :label="t('rooms.roomNumber')" name="roomNumber" required>
            <UInput v-model="form.roomNumber" :placeholder="t('rooms.roomNumberPlaceholder')" :disabled="loading" />
          </UFormField>

          <!-- Room Type -->
          <UFormField :label="t('rooms.roomType')" name="roomType" required>
            <USelect v-model="form.roomType" :items="roomTypes" :disabled="loading" />
          </UFormField>

          <!-- Floor -->
          <UFormField :label="t('rooms.floor')" name="floor">
            <UInput
              v-model.number="form.floor"
              type="number"
              :placeholder="t('rooms.floorPlaceholder')"
              :disabled="loading"
            />
          </UFormField>

          <!-- Capacity -->
          <UFormField :label="t('rooms.capacityGuests')" name="capacity" required>
            <UInput v-model.number="form.capacity" type="number" min="1" :disabled="loading" />
          </UFormField>

          <!-- Base Price -->
          <UFormField :label="t('rooms.basePrice')" name="basePrice" required>
            <UInput
              v-model.number="form.basePrice"
              type="number"
              min="0"
              step="0.01"
              :disabled="loading"
            />
          </UFormField>

          <!-- Status -->
          <UFormField :label="t('rooms.status')" name="status" required>
            <USelect v-model="form.status" :items="statusOptions" :disabled="loading" />
          </UFormField>

          <!-- Description -->
          <UFormField :label="t('forms.description')" name="description" class="md:col-span-2">
            <UTextarea
              v-model="form.description"
              :placeholder="t('rooms.descriptionPlaceholder')"
              :rows="3"
              :disabled="loading"
            />
          </UFormField>

          <!-- Amenities -->
          <UFormField :label="t('rooms.amenities')" name="amenities" class="md:col-span-2">
            <div class="flex flex-wrap gap-2">
              <UBadge
                v-for="amenity in availableAmenities"
                :key="amenity.id"
                :variant="form.amenities.includes(amenity) ? 'solid' : 'outline'"
                :color="form.amenities.includes(amenity) ? 'primary' : 'neutral'"
                class="cursor-pointer"
                @click="toggleAmenity(amenity)"
              >
                {{ amenity.name }}
              </UBadge>
            </div>
          </UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/rooms" :disabled="loading"> {{ t("actions.cancel") }} </UButton>
          <UButton type="submit" color="primary" :loading="loading"> {{ t("rooms.createRoom") }} </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager"],
});
const { t } = useI18n();

const loading = ref(false);

const form = reactive({
  roomNumber: "",
  roomType: "single" as "single" | "double" | "suite" | "deluxe",
  floor: null as number | null,
  capacity: 2,
  basePrice: 0,
  status: "available" as "available" | "occupied" | "maintenance" | "out_of_order",
  description: "",
  amenities: [] as ListItem[],
});

const roomTypes = computed(() => [
  { value: "single", label: t("rooms.types.single") },
  { value: "double", label: t("rooms.types.double") },
  { value: "suite", label: t("rooms.types.suite") },
  { value: "deluxe", label: t("rooms.types.deluxe") },
]);

const statusOptions = computed(() => [
  { value: "available", label: t("rooms.statuses.available") },
  { value: "occupied", label: t("rooms.statuses.occupied") },
  { value: "maintenance", label: t("rooms.statuses.maintenance") },
  { value: "out_of_order", label: t("rooms.statuses.outOfOrder") },
]);

const { data: availableAmenities } = useAsyncData(async () => {
  const res = await $fetch<{ data: ListItem[] }>("/api/rooms/amenities");
  return res.data;
});
const toggleAmenity = (amenity: ListItem) => {
  const index = form.amenities.indexOf(amenity);
  if (index === -1) {
    form.amenities.push(amenity);
  } else {
    form.amenities.splice(index, 1);
  }
};

const handleSubmit = async () => {
  loading.value = true;
  try {
    await $fetch("/api/rooms", {
      method: "POST",
      body: form,
    });

    await navigateTo("/rooms");
  } catch (error) {
    console.error("Failed to create room:", error);
  } finally {
    loading.value = false;
  }
};
</script>
