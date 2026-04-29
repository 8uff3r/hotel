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
      <UForm @submit="handleSubmit" :state="form" :schema>
        <div class="flex flex-col gap-16">
          <div class="flex w-full gap-4 max-md:flex-col">
            <div class="grid w-full grid-cols-1 gap-6 md:grid-cols-2">
              <!-- Room Number -->
              <UFormField :label="t('rooms.roomNumber')" name="roomNumber" required>
                <UInput
                  v-model="form.roomNumber"
                  :placeholder="t('rooms.roomNumberPlaceholder')"
                  :disabled="loading"
                />
              </UFormField>

              <UFormField :label="t('rooms.name')" name="name" required>
                <UInput v-model="form.name" :disabled="loading" />
              </UFormField>

              <!-- Room Type -->
              <UFormField :label="t('rooms.roomType')" name="roomType" required>
                <HSelect v-model="form.roomTypeId" :items="types ?? []" :disabled="loading" />
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
                <HSelect v-model="form.statusId" :items="statuses ?? []" :disabled="loading" />
              </UFormField>
            </div>
            <!-- Description -->
            <UFormField
              :label="t('forms.description')"
              name="description"
              class="w-full"
              :ui="{
                container: 'h-full',
                root: 'flex flex-col',
              }"
            >
              <UTextarea
                v-model="form.description"
                :placeholder="t('rooms.descriptionPlaceholder')"
                :disabled="loading"
                :ui="{
                  root: 'w-full h-full',
                  base: 'h-full resize-none',
                }"
              />
            </UFormField>
          </div>

          <!-- Amenities -->
          <UFormField :label="t('rooms.amenities')" name="amenities" class="md:col-span-2">
            <div class="flex flex-wrap gap-2">
              <UBadge
                v-for="amenity in availableAmenities"
                :key="amenity.id"
                :variant="form.amenities?.find((v) => v.id === amenity.id) ? 'solid' : 'outline'"
                :color="form.amenities?.find((v) => v.id === amenity.id) ? 'primary' : 'neutral'"
                class="cursor-pointer"
                @click="toggleAmenity(amenity.id!)"
              >
                {{ amenity.label }}
              </UBadge>
            </div>
          </UFormField>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton variant="outline" to="/rooms" :disabled="loading">
            {{ t("actions.cancel") }}
          </UButton>
          <UButton type="submit" color="primary" :loading="loading">
            {{ t("rooms.createRoom") }}
          </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import type { FormSubmitEvent } from "@nuxt/ui";
import type z from "zod";
import { zRoom } from "~/utils/client/zod.gen";

definePageMeta({
  requiresPermission: PERMISSIONS.rooms.rooms.create,
});
const { t } = useI18n();

const loading = ref(false);

const schema = zRoom;
type Schema = z.output<typeof schema>;

const form = ref<Schema>({ amenities: [], hotelId: undefined } as any);

const { data: availableAmenities } = useAsyncData("room-amenities", async () => {
  const res = await getApiRoomsAmenities({});
  return res.data;
});
const { data: statuses } = useAsyncData("room-statuses", async () => {
  const res = await getApiRoomsStatuses({});
  return res.data;
});
const { data: types } = useAsyncData("room-types", async () => {
  const res = await getApiRoomsTypes({});
  return res.data;
});

const toggleAmenity = (amenityId: number) => {
  const index = form.value.amenities?.findIndex((a) => a.id === amenityId);
  if (index !== undefined && index !== -1) {
    form.value.amenities?.splice(index, 1);
  } else {
    const amenity = availableAmenities.value?.find((v) => v.id === amenityId);
    console.log(amenity);
    if (amenity) form.value.amenities?.push(amenity);
  }
};

const handleSubmit = async (event: FormSubmitEvent<Schema>) => {
  loading.value = true;
  try {
    await postApiRooms({
      body: event.data,
    });

    navigateTo("/rooms");
  } catch (error) {
    console.error("Failed to create room:", error);
  } finally {
    loading.value = false;
  }
};
</script>
