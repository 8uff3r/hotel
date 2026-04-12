<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/rooms" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />{{ t("rooms.back_to_rooms") }}</UButton
      >
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Room {{ room?.roomNumber }}</h1>
    </div>

    <div v-if="pending" class="flex justify-center py-12">
      <UIcon name="i-lucide-loader-2" class="h-8 w-8 animate-spin" />
    </div>

    <div v-else-if="room">
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Room Details -->
        <div class="lg:col-span-2">
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <span class="text-lg font-semibold">{{ t("common.room_details") }}</span>
                <UBadge :color="getStatusColor(room.status)" variant="soft">
                  {{ room.status }}
                </UBadge>
              </div>
            </template>

            <form @submit.prevent="handleSubmit">
              <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
                <!-- Room Number -->
                <UFormField :label="t('rooms.room_number')" name="roomNumber" required>
                  <UInput
                    v-model="form.roomNumber"
                    :placeholder="t('rooms.e_g_101_205a')"
                    :disabled="saving"
                  />
                </UFormField>

                <!-- Room Type -->
                <UFormField :label="t('rooms.room_type')" name="roomType" required>
                  <USelect v-model="form.roomType" :items="roomTypes" :disabled="saving" />
                </UFormField>

                <!-- Floor -->
                <UFormField :label="t('common.floor')" name="floor">
                  <UInput
                    v-model.number="form.floor"
                    type="number"
                    :placeholder="t('rooms.e_g_1_2')"
                    :disabled="saving"
                  />
                </UFormField>

                <!-- Capacity -->
                <UFormField :label="t('rooms.capacity_guests')" name="capacity" required>
                  <UInput v-model.number="form.capacity" type="number" min="1" :disabled="saving" />
                </UFormField>

                <!-- Base Price -->
                <UFormField :label="t('rooms.base_price')" name="basePrice" required>
                  <UInput
                    v-model.number="form.basePrice"
                    type="number"
                    min="0"
                    step="0.01"
                    :disabled="saving"
                  />
                </UFormField>

                <!-- Status -->
                <UFormField :label="t('common.status')" name="status" required>
                  <USelect v-model="form.status" :items="statusOptions" :disabled="saving" />
                </UFormField>

                <!-- Description -->
                <UFormField
                  :label="t('common.description')"
                  name="description"
                  class="md:col-span-2"
                >
                  <UTextarea
                    v-model="form.description"
                    :placeholder="t('rooms.room_description')"
                    :rows="3"
                    :disabled="saving"
                  />
                </UFormField>

                <!-- Amenities -->
                <UFormField :label="t('rooms.amenities')" name="amenities" class="md:col-span-2">
                  <div class="flex flex-wrap gap-2">
                    <UBadge
                      v-for="amenity in availableAmenities"
                      :key="amenity"
                      :variant="form.amenities.includes(amenity) ? 'solid' : 'outline'"
                      :color="form.amenities.includes(amenity) ? 'primary' : 'neutral'"
                      class="cursor-pointer"
                      @click="toggleAmenity(amenity)"
                    >
                      {{ amenity }}
                    </UBadge>
                  </div>
                </UFormField>
              </div>

              <div class="mt-6 flex justify-end gap-3">
                <UButton variant="outline" to="/rooms" :disabled="saving"> Cancel </UButton>
                <UButton type="submit" color="primary" :loading="saving"> Save Changes </UButton>
              </div>
            </form>
          </UCard>
        </div>

        <!-- Sidebar Info -->
        <div class="space-y-6">
          <!-- Room Type Badge -->
          <UCard>
            <template #header>
              <span class="font-semibold">Room Type</span>
            </template>
            <div class="flex items-center gap-3">
              <UIcon name="i-lucide-bed-double" class="h-6 w-6 text-primary" />
              <span class="text-lg capitalize">{{ room.roomType }}</span>
            </div>
          </UCard>

          <!-- Capacity -->
          <UCard>
            <template #header>
              <span class="font-semibold">Capacity</span>
            </template>
            <div class="flex items-center gap-3">
              <UIcon name="i-lucide-users" class="h-6 w-6 text-primary" />
              <span class="text-lg">{{ room.capacity }} guests</span>
            </div>
          </UCard>

          <!-- Price -->
          <UCard>
            <template #header>
              <span class="font-semibold">Base Price</span>
            </template>
            <div class="flex items-center gap-3">
              <UIcon name="i-lucide-dollar-sign" class="h-6 w-6 text-primary" />
              <span class="text-lg">${{ room.basePrice.toFixed(2) }}</span>
            </div>
          </UCard>

          <!-- Quick Actions -->
          <UCard>
            <template #header>
              <span class="font-semibold">Quick Actions</span>
            </template>
            <div class="space-y-2">
              <UButton
                v-if="room.status === 'available'"
                to="/reservations/create"
                variant="outline"
                block
                color="success"
              >
                <UIcon name="i-lucide-plus" class="mr-2" />
                New Reservation
              </UButton>
              <UButton v-if="room.status === 'occupied'" variant="outline" block color="warning">
                <UIcon name="i-lucide-log-out" class="mr-2" />
                Check-out Guest
              </UButton>
            </div>
          </UCard>
        </div>
      </div>
    </div>

    <!-- Not Found -->
    <div v-else class="py-12 text-center">
      <UIcon name="i-lucide-alert-circle" class="mx-auto h-12 w-12 text-gray-400" />
      <p class="mt-4 text-lg text-gray-500">Room not found</p>
      <UButton to="/rooms" class="mt-4"> Back to Rooms </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});
const { t } = useI18n();

const route = useRoute();
const roomId = route.params.id as string;

const saving = ref(false);

const form = reactive({
  roomNumber: "",
  roomType: "single" as "single" | "double" | "suite" | "deluxe",
  floor: null as number | null,
  capacity: 2,
  basePrice: 0,
  status: "available" as "available" | "occupied" | "maintenance" | "out_of_order",
  description: "",
  amenities: [] as string[],
});

const roomTypes = [
  { value: "single", label: "Single" },
  { value: "double", label: "Double" },
  { value: "suite", label: "Suite" },
  { value: "deluxe", label: "Deluxe" },
];

const statusOptions = [
  { value: "available", label: "Available" },
  { value: "occupied", label: "Occupied" },
  { value: "maintenance", label: "Maintenance" },
  { value: "out_of_order", label: "Out of Order" },
];

const availableAmenities = [
  "WiFi",
  "TV",
  "Air Conditioning",
  "Mini Bar",
  "Safe",
  "Ocean View",
  "City View",
  "Balcony",
  "Jacuzzi",
  "Room Service",
];

const toggleAmenity = (amenity: string) => {
  const index = form.amenities.indexOf(amenity);
  if (index === -1) {
    form.amenities.push(amenity);
  } else {
    form.amenities.splice(index, 1);
  }
};

const { data: room, pending } = useAsyncData(async () => {
  const response = await $fetch<Room>(`/api/rooms/${roomId}`);

  if (response) {
    form.roomNumber = response.roomNumber;
    form.roomType = response.roomType;
    form.floor = response.floor;
    form.capacity = response.capacity;
    form.basePrice = response.basePrice;
    form.status = response.status;
    form.description = response.description || "";
    form.amenities = Array.isArray(response.amenities) ? response.amenities : [];
  }
  return response;
});

const handleSubmit = async () => {
  saving.value = true;
  try {
    await $fetch(`/api/rooms/${roomId}`, {
      method: "PUT",
      body: form,
    });

    await navigateTo("/rooms");
  } catch (error) {
    console.error("Failed to update room:", error);
  } finally {
    saving.value = false;
  }
};

const getStatusColor = (status: string): "success" | "warning" | "info" | "error" | "neutral" => {
  const colors: Record<string, "success" | "warning" | "info" | "error" | "neutral"> = {
    available: "success",
    occupied: "warning",
    maintenance: "info",
    out_of_order: "error",
  };
  return colors[status] || "neutral";
};
</script>
