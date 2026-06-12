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
        <div class="space-y-6 lg:col-span-2">
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <span class="text-lg font-semibold">{{ t("common.room_details") }}</span>
                <div class="flex items-center gap-2">
                  <UBadge
                    color="neutral"
                    :style="{
                      color: `#${room.status?.colorHex}`,
                    }"
                    variant="soft"
                  >
                    {{ room.status?.label }}
                  </UBadge>
                </div>
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
                  <HSelect v-model="form.roomTypeId" :items="types" :disabled="saving" />
                </UFormField>

                <!-- Floor -->
                <UFormField :label="t('common.floor')" name="floorId">
                  <HSelect v-model="form.floorId" :items="floors" :disabled="saving" />
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
                  <HSelect v-model="form.statusId" :items="statuses" :disabled="saving" />
                </UFormField>

                <!-- Description -->
                <UFormField
                  :label="t('common.description')"
                  name="description"
                  class="md:col-span-2"
                >
                  <UTextarea
                    v-model="form.description"
                    class="w-full"
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
                      :key="amenity.id"
                      :variant="
                        form.amenities?.find((v) => v.id === amenity.id) ? 'solid' : 'outline'
                      "
                      :color="
                        form.amenities?.find((v) => v.id === amenity.id) ? 'primary' : 'neutral'
                      "
                      class="cursor-pointer"
                      @click="toggleAmenity(amenity.id!)"
                    >
                      {{ amenity.label }}
                    </UBadge>
                  </div>
                </UFormField>
              </div>

              <div class="mt-6 flex justify-end gap-3">
                <UButton variant="outline" to="/rooms" :disabled="saving">
                  {{ t("actions.cancel") }}
                </UButton>
                <UButton type="submit" color="primary" :loading="saving">
                  {{ t("actions.saveChanges") }}
                </UButton>
              </div>
            </form>
          </UCard>

          <!-- Pictures Section -->
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <span class="text-lg font-semibold">{{ t("rooms.pictures") }}</span>
              </div>
            </template>
            <div class="space-y-4">
              <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
                <div v-for="pic in pictures" :key="pic.id" class="group relative">
                  <img
                    :src="pic.url"
                    :alt="pic.description"
                    class="h-24 w-full rounded-lg border object-cover"
                  />
                  <UButton
                    variant="ghost"
                    size="xs"
                    color="error"
                    class="absolute top-1 right-1 opacity-0 transition-opacity group-hover:opacity-100"
                    @click="removePicture(pic.id)"
                  >
                    <UIcon name="i-lucide-trash" class="h-3 w-3" />
                  </UButton>
                </div>
                <div
                  class="flex h-24 items-center justify-center rounded-lg border border-dashed border-gray-300 dark:border-gray-600"
                >
                  <UButton variant="ghost" size="sm" @click="showAddPicture = true">
                    <UIcon name="i-lucide-plus" class="mr-1 h-4 w-4" />
                    افزودن تصویر
                  </UButton>
                </div>
              </div>
            </div>
          </UCard>

          <!-- Add Picture Modal -->
          <UModal v-model:open="showAddPicture">
            <template #title>افزودن تصویر</template>
            <template #content>
              <div class="space-y-4 p-4">
                <UFormField label="آدرس تصویر (URL)" required>
                  <UInput v-model="newPicture.url" placeholder="https://..." />
                </UFormField>
                <UFormField label="توضیحات">
                  <UInput v-model="newPicture.description" placeholder="توضیحات تصویر" />
                </UFormField>
                <div class="flex justify-end gap-2">
                  <UButton variant="outline" size="sm" @click="showAddPicture = false"
                    >انصراف</UButton
                  >
                  <UButton size="sm" color="primary" :loading="addingPicture" @click="addPicture"
                    >ذخیره</UButton
                  >
                </div>
              </div>
            </template>
          </UModal>
        </div>

        <!-- Sidebar Info -->
        <div class="space-y-6">
          <!-- Room Type Badge -->
          <UCard>
            <template #header>
              <span class="font-semibold">{{ t("rooms.room_type") }}</span>
            </template>
            <div class="flex items-center gap-3">
              <UIcon name="i-lucide-bed-double" class="h-6 w-6 text-primary" />
              <span class="text-lg capitalize">{{ room.roomType?.label }}</span>
            </div>
          </UCard>

          <!-- Capacity -->
          <UCard>
            <template #header>
              <span class="font-semibold">{{ t("rooms.capacity_guests") }}</span>
            </template>
            <div class="flex items-center gap-3">
              <UIcon name="i-lucide-users" class="h-6 w-6 text-primary" />
              <span class="text-lg">{{ room.capacity }}</span>
            </div>
          </UCard>

          <!-- Price -->
          <UCard>
            <template #header>
              <span class="font-semibold">{{ t("rooms.base_price") }}</span>
            </template>
            <div class="flex items-center gap-3">
              <UIcon name="i-lucide-dollar-sign" class="h-6 w-6 text-primary" />
              <span class="text-lg">${{ room.basePrice?.toFixed(2) }}</span>
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
import type z from "zod";
import {
  getApiRoomsIdPictures,
  postApiRoomsIdPictures,
  deleteApiRoomsIdPicturesPictureId,
  putApiRoomsId,
} from "~/utils/client";
import { zRoom } from "~/utils/client/zod.gen";

definePageMeta({
  requiresPermission: PERMISSIONS.rooms.rooms.read,
});
const { t } = useI18n();
const toast = useToast();

const route = useRoute();
const roomId = route.params.id as string;

const saving = ref(false);

const schema = zRoom;
type Schema = z.output<typeof schema>;
const form = ref<Schema>({
  amenities: [] as string[],
} as any);

const { data: availableAmenities } = useAsyncData("room-amenities", async () => {
  const res = await getApiRoomsAmenities({});
  return res.data?.data;
});
const { data: statuses } = useAsyncData("room-statuses", async () => {
  const res = await getApiRoomsStatuses({});
  return res.data?.data;
});
const { data: types } = useAsyncData("room-types", async () => {
  const res = await getApiRoomsTypes({});
  return res.data?.data;
});
const { data: floors } = useAsyncData("room-floors", async () => {
  const res = await getApiRoomsFloors({});
  return (
    res.data?.data?.map((f: any) => ({
      ...f,
      label: `Floor ${f.number}${f.description ? ` - ${f.description}` : ""}`,
    })) ?? []
  );
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

const { data: room, pending } = useAsyncData(async () => {
  const response = await getApiRoomsId({
    path: {
      id: roomId,
    },
  });

  form.value = response.data ?? ({ amenities: [] } as any);
  return response.data;
});

const { data: pictures, refresh: refreshPictures } = useAsyncData(
  `room-pictures-${roomId}`,
  async () => {
    const res = await getApiRoomsIdPictures({ path: { id: roomId } });
    return res.data?.data ?? [];
  }
);

const showAddPicture = ref(false);
const addingPicture = ref(false);
const newPicture = reactive({ url: "", description: "" });

const addPicture = async () => {
  if (!newPicture.url) return;
  addingPicture.value = true;
  try {
    await postApiRoomsIdPictures({
      path: { id: roomId },
      body: { url: newPicture.url, description: newPicture.description },
    });
    toast.add({ title: "تصویر اضافه شد", color: "success" });
    newPicture.url = "";
    newPicture.description = "";
    showAddPicture.value = false;
    await refreshPictures();
  } catch (e) {
    toast.add({ title: "خطا در افزودن تصویر", color: "error" });
  } finally {
    addingPicture.value = false;
  }
};

const removePicture = async (pictureId: number | undefined) => {
  if (!pictureId) return;
  try {
    await deleteApiRoomsIdPicturesPictureId({
      path: { id: String(roomId), pictureId: String(pictureId) },
    });
    toast.add({ title: "تصویر حذف شد", color: "success" });
    await refreshPictures();
  } catch (e) {
    toast.add({ title: "خطا در حذف تصویر", color: "error" });
  }
};

const handleSubmit = async () => {
  saving.value = true;
  try {
    await putApiRoomsId({
      path: {
        id: roomId,
      },
      body: form.value,
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
    reserved: "info",
    under_repair: "error",
    cleaning: "neutral",
  };
  return colors[status] || "neutral";
};
</script>
