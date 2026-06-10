<template>
  <div class="p-6 max-w-4xl mx-auto">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">تنظیمات هتل</h1>
      <UButton variant="outline" to="/hotels/profile">بازگشت به پروفایل</UButton>
    </div>

    <UCard>
      <template #header>
        <h2 class="font-semibold">تنظیمات عمومی</h2>
      </template>
      <UForm :state="form" @submit="saveSettings">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <UFormField label="ساعت ورود استاندارد" name="standardCheckInTime" required>
            <UInput v-model="form.standardCheckInTime" placeholder="14:00" />
          </UFormField>
          <UFormField label="ساعت خروج استاندارد" name="standardCheckOutTime" required>
            <UInput v-model="form.standardCheckOutTime" placeholder="12:00" />
          </UFormField>
          <UFormField label="ساعت حسابرسی شب" name="nightAuditHour" required>
            <UInput v-model="form.nightAuditHour" placeholder="03:00" />
          </UFormField>
        </div>
        <div class="mt-4 flex gap-3">
          <UButton type="submit" :loading="saving">ذخیره</UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template>

<script setup lang="ts">
import { useToast } from "@nuxt/ui/composables";
import { useI18n } from "vue-i18n";
import { getApiHotelsIdSettings, putApiHotelsIdSettings } from "~/utils/client";
import { useAuthStore } from "~/stores/auth";

const authStore = useAuthStore();
const { t } = useI18n();
const hotelId = computed(() => authStore.currentHotelId);

const toast = useToast();
const saving = ref(false);

const form = reactive({
  standardCheckInTime: "14:00",
  standardCheckOutTime: "12:00",
  nightAuditHour: "03:00",
});

const { data: settings } = useAsyncData(
  "hotel-settings-edit",
  async () => {
    if (!hotelId.value) return null;
    const res = await getApiHotelsIdSettings({ path: { id: hotelId.value } });
    return res.data?.setting;
  },
  { watch: [hotelId] }
);

watchEffect(() => {
  if (settings.value) {
    form.standardCheckInTime = settings.value.standardCheckInTime ?? "14:00";
    form.standardCheckOutTime = settings.value.standardCheckOutTime ?? "12:00";
    form.nightAuditHour = settings.value.nightAuditHour ?? "03:00";
  }
});

async function saveSettings() {
  saving.value = true;
  try {
    await putApiHotelsIdSettings({
      path: { id: hotelId.value },
      body: {
        hotelId: hotelId.value,
        standardCheckInTime: form.standardCheckInTime,
        standardCheckOutTime: form.standardCheckOutTime,
        nightAuditHour: form.nightAuditHour,
      },
    });
    toast.add({ title: "تنظیمات ذخیره شد", color: "success" });
  } catch (e: any) {
    console.error(e);
    toast.add({ title: "خطا در ذخیره تنظیمات", color: "error" });
  } finally {
    saving.value = false;
  }
}
</script>
