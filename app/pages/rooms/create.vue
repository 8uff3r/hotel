hotel/app/pages/rooms/create.vue#L1-120
<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/rooms" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        Back to Rooms
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Add New Room</h1>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <!-- Room Number -->
          <UFormField label="Room Number" name="roomNumber" required>
            <UInput v-model="form.roomNumber" placeholder="e.g., 101, 205A" :disabled="loading" />
          </UFormField>

          <!-- Room Type -->
          <UFormField label="Room Type" name="roomType" required>
            <USelect v-model="form.roomType" :items="roomTypes" :disabled="loading" />
          </UFormField>

          <!-- Floor -->
          <UFormField label="Floor" name="floor">
            <UInput
              v-model.number="form.floor"
              type="number"
              placeholder="e.g., 1, 2"
              :disabled="loading"
            />
          </UFormField>

          <!-- Capacity -->
          <UFormField label="Capacity (guests)" name="capacity" required>
            <UInput v-model.number="form.capacity" type="number" min="1" :disabled="loading" />
          </UFormField>

          <!-- Base Price -->
          <UFormField label="Base Price ($)" name="basePrice" required>
            <UInput
              v-model.number="form.basePrice"
              type="number"
              min="0"
              step="0.01"
              :disabled="loading"
            />
          </UFormField>

          <!-- Status -->
          <UFormField label="Status" name="status" required>
            <USelect v-model="form.status" :items="statusOptions" :disabled="loading" />
          </UFormField>

          <!-- Description -->
          <UFormField label="Description" name="description" class="md:col-span-2">
            <UTextarea
              v-model="form.description"
              placeholder="Room description..."
              :rows="3"
              :disabled="loading"
            />
          </UFormField>

          <!-- Amenities -->
          <UFormField label="Amenities" name="amenities" class="md:col-span-2">
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
          <UButton variant="outline" to="/rooms" :disabled="loading"> Cancel </UButton>
          <UButton type="submit" color="primary" :loading="loading"> Create Room </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager"],
});

const loading = ref(false);

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
