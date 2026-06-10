<template>
  <div class="p-6 max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold mb-6">ویرایش مدیر سیستم</h1>
    <UForm :schema="schema" :state="state" @submit="onSubmit">
      <div class="grid grid-cols-2 gap-4">
        <UFormGroup label="نام" name="firstName">
          <UInput v-model="state.firstName" />
        </UFormGroup>
        <UFormGroup label="نام خانوادگی" name="lastName">
          <UInput v-model="state.lastName" />
        </UFormGroup>
        <UFormGroup label="ایمیل" name="email">
          <UInput v-model="state.email" type="email" />
        </UFormGroup>
        <UFormGroup label="نام کاربری" name="username">
          <UInput v-model="state.username" />
        </UFormGroup>
        <UFormGroup label="شماره تماس" name="contactNumber">
          <UInput v-model="state.contactNumber" />
        </UFormGroup>
        <UFormGroup label="نقش" name="role">
          <UInput v-model="state.role" />
        </UFormGroup>
      </div>
      <UFormGroup label="هتل‌ها" name="hotelIds" class="mt-4">
        <USelect v-model="state.hotelIds" multiple :options="hotelOptions" />
      </UFormGroup>
      <UFormGroup label="سوپر ادمین" name="isSuperAdmin" class="mt-4">
        <UToggle v-model="state.isSuperAdmin" />
      </UFormGroup>
      <UFormGroup label="فعال" name="isActive" class="mt-4">
        <UToggle v-model="state.isActive" />
      </UFormGroup>
      <div class="mt-6 flex gap-3">
        <UButton type="submit" :loading="submitting">ذخیره</UButton>
        <UButton variant="outline" to="/admins">انصراف</UButton>
      </div>
    </UForm>
  </div>
</template>

<script setup lang="ts">
import type { Hotel, SanitizedAdmin } from "~/utils/client";
import { putApiAdminsId } from "~/utils/client";

definePageMeta({
  requiresPermission: "users:update",
});

const route = useRoute();
const adminId = route.params.id as string;
const { t } = useI18n();

const schema = {
  firstName: { type: "string", required: true },
  lastName: { type: "string", required: true },
  email: { type: "string", required: true },
  username: { type: "string", required: true },
};

const state = reactive({
  firstName: "",
  lastName: "",
  email: "",
  username: "",
  contactNumber: "",
  role: "",
  isSuperAdmin: false,
  isActive: true,
  hotelIds: [] as string[],
});

const { data: admin } = useFetch(`/api/admins/${adminId}`, {
  key: `admin-${adminId}`,
  transform: (res) => (res as { data?: SanitizedAdmin })?.data,
});

const { data: hotels } = useFetch("/api/hotels", {
  key: "hotels-list",
  transform: (res) => (res as { data?: Hotel[] })?.data ?? [],
});

const hotelOptions = computed(() =>
  (hotels.value ?? []).map((h) => ({ label: h.name ?? "", value: h.id ?? "" }))
);

watchEffect(() => {
  if (admin.value) {
    state.firstName = admin.value.firstName ?? "";
    state.lastName = admin.value.lastName ?? "";
    state.email = admin.value.email ?? "";
    state.username = admin.value.username ?? "";
    state.contactNumber = admin.value.contactNumber ?? "";
    state.role = admin.value.role ?? "";
    state.isSuperAdmin = admin.value.isSuperAdmin ?? false;
    state.isActive = admin.value.isActive ?? true;
    state.hotelIds = (admin.value.adminHotels ?? []).map((h) => h.hotelId ?? "");
  }
});

const submitting = ref(false);

async function onSubmit() {
  submitting.value = true;
  try {
    await putApiAdminsId({
      path: { id: adminId },
      body: state,
      requestValidator: undefined,
    });
    navigateTo("/admins");
  } catch (e) {
    console.error(e);
  } finally {
    submitting.value = false;
  }
}
</script>
