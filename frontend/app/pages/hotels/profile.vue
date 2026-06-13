<template>
  <div class="mx-auto max-w-4xl p-6">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold">پروفایل هتل</h1>
      <div class="flex gap-2">
        <UButton v-if="!editing" @click="editing = true">ویرایش</UButton>
        <UButton v-if="editing" variant="outline" @click="editing = false">انصراف</UButton>
      </div>
    </div>

    <div v-if="hotel" class="space-y-6">
      <UCard>
        <template #header>
          <h2 class="font-semibold">اطلاعات اصلی</h2>
        </template>
        <div v-if="!editing" class="grid grid-cols-2 gap-4">
          <p><strong>نام:</strong> {{ hotel.name }}</p>
          <p><strong>کد:</strong> {{ hotel.code }}</p>
          <p><strong>آدرس:</strong> {{ hotel.address }}</p>
          <p><strong>تلفن:</strong> {{ hotel.phone || "-" }}</p>
          <p><strong>ایمیل:</strong> {{ hotel.email || "-" }}</p>
          <p><strong>ظرفیت:</strong> {{ hotel.totalCapacity || "-" }}</p>
          <p><strong>تعداد طبقات:</strong> {{ hotel.numberOfFloors || "-" }}</p>
          <p><strong>نام مدیرعامل:</strong> {{ hotel.ceoName || "-" }}</p>
          <p><strong>امکانات نزدیک:</strong> {{ hotel.nearbyFacilities || "-" }}</p>
        </div>
        <UForm v-else :state="editState" @submit="saveHotel">
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="نام" name="name">
              <UInput v-model="editState.name" />
            </UFormGroup>
            <UFormGroup label="کد" name="code">
              <UInput v-model="editState.code" disabled />
            </UFormGroup>
            <UFormGroup label="آدرس" name="address">
              <UInput v-model="editState.address" />
            </UFormGroup>
            <UFormGroup label="تلفن" name="phone">
              <UInput v-model="editState.phone" />
            </UFormGroup>
            <UFormGroup label="ایمیل" name="email">
              <UInput v-model="editState.email" />
            </UFormGroup>
            <UFormGroup label="ظرفیت" name="totalCapacity">
              <UInput v-model.number="editState.totalCapacity" type="number" />
            </UFormGroup>
            <UFormGroup label="تعداد طبقات" name="numberOfFloors">
              <UInput v-model.number="editState.numberOfFloors" type="number" />
            </UFormGroup>
            <UFormGroup label="نام مدیرعامل" name="ceoName">
              <UInput v-model="editState.ceoName" />
            </UFormGroup>
            <UFormGroup label="امکانات نزدیک" name="nearbyFacilities" class="col-span-2">
              <UTextarea v-model="editState.nearbyFacilities" />
            </UFormGroup>
          </div>
          <div class="mt-4 flex gap-3">
            <UButton type="submit" :loading="saving">ذخیره</UButton>
          </div>
        </UForm>
      </UCard>

      <UCard>
        <template #header>
          <h2 class="font-semibold">تنظیمات هتل</h2>
        </template>
        <div v-if="settings" class="grid grid-cols-3 gap-4">
          <p><strong>ساعت ورود استاندارد:</strong> {{ settings.standardCheckInTime }}</p>
          <p><strong>ساعت خروج استاندارد:</strong> {{ settings.standardCheckOutTime }}</p>
          <p><strong>ساعت حسابرسی شب:</strong> {{ settings.nightAuditHour }}</p>
        </div>
        <div v-else>تنظیماتی ثبت نشده</div>
        <div class="mt-4">
          <UButton size="xs" to="/settings">ویرایش تنظیمات</UButton>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h2 class="font-semibold">تصاویر هتل</h2>
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
                @click="pic.id && removePicture(pic.id)"
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

      <UModal v-model:open="showAddPicture">
        <template #title>افزودن تصویر هتل</template>
        <template #content>
          <div class="space-y-4 p-4">
            <UFormField label="آدرس تصویر (URL)" required>
              <UInput v-model="newPicture.url" placeholder="https://..." />
            </UFormField>
            <UFormField label="توضیحات">
              <UInput v-model="newPicture.description" placeholder="توضیحات تصویر" />
            </UFormField>
            <div class="flex justify-end gap-2">
              <UButton variant="outline" size="sm" @click="showAddPicture = false">انصراف</UButton>
              <UButton size="sm" color="primary" :loading="addingPicture" @click="addPicture"
                >ذخیره</UButton
              >
            </div>
          </div>
        </template>
      </UModal>
    </div>
    <div v-else>هتل یافت نشد</div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";

const authStore = useAuthStore();
const hotelId = computed(() => authStore.currentHotelId);
const toast = useToast();

const editing = ref(false);
const saving = ref(false);

const { data: hotel } = useAsyncData(
  "hotel-profile",
  async () => {
    if (!hotelId.value) return null;
    const res = await getApiHotelsId({ path: { id: hotelId.value } });
    return res.data;
  },
  { watch: [hotelId] }
);

const { data: settings } = useAsyncData(
  "hotel-settings",
  async () => {
    if (!hotelId.value) return null;
    const res = await getApiHotelsIdSettings({ path: { id: hotelId.value } });
    return res.data?.setting;
  },
  { watch: [hotelId] }
);

const { data: pictures, refresh: refreshPictures } = useAsyncData(
  "hotel-pictures",
  async () => {
    if (!hotelId.value) return null;
    const res = await getApiHotelsIdPictures({ path: { id: hotelId.value } });
    return res.data?.data ?? [];
  },
  { watch: [hotelId] }
);

const showAddPicture = ref(false);
const addingPicture = ref(false);
const newPicture = reactive({ url: "", description: "" });

const addPicture = async () => {
  if (!newPicture.url) return;
  addingPicture.value = true;
  try {
    await postApiHotelsIdPictures({
      path: { id: hotelId.value },
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
    await deleteApiHotelsIdPicturesPictureId({
      path: { id: hotelId.value, pictureId: String(pictureId) },
    });
    toast.add({ title: "تصویر حذف شد", color: "success" });
    await refreshPictures();
  } catch (e) {
    toast.add({ title: "خطا در حذف تصویر", color: "error" });
  }
};

const editState = reactive({
  name: "",
  code: "",
  address: "",
  phone: "",
  email: "",
  totalCapacity: 0,
  numberOfFloors: 0,
  ceoName: "",
  nearbyFacilities: "",
});

watchEffect(() => {
  if (hotel.value) {
    editState.name = hotel.value.name ?? "";
    editState.code = hotel.value.code ?? "";
    editState.address = hotel.value.address ?? "";
    editState.phone = hotel.value.phone ?? "";
    editState.email = hotel.value.email ?? "";
    editState.totalCapacity = hotel.value.totalCapacity ?? 0;
    editState.numberOfFloors = hotel.value.numberOfFloors ?? 0;
    editState.ceoName = hotel.value.ceoName ?? "";
    editState.nearbyFacilities = hotel.value.nearbyFacilities ?? "";
  }
});

async function saveHotel() {
  saving.value = true;
  try {
    await putApiHotelsId({ path: { id: hotelId.value }, body: editState });
    editing.value = false;
    navigateTo("/hotels/profile", { replace: true });
  } catch (e) {
    console.error(e);
  } finally {
    saving.value = false;
  }
}
</script>
