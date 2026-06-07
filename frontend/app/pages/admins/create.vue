<template>
  <div class="p-6 max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold mb-6">افزودن مدیر سیستم</h1>
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
        <UFormGroup label="رمز عبور" name="password">
          <UInput v-model="state.password" type="password" />
        </UFormGroup>
        <UFormGroup label="شماره تماس" name="contactNumber">
          <UInput v-model="state.contactNumber" />
        </UFormGroup>
        <UFormGroup label="نقش" name="role">
          <UInput v-model="state.role" />
        </UFormGroup>
      </div>
      <UFormGroup label="سوپر ادمین" name="isSuperAdmin" class="mt-4">
        <UToggle v-model="state.isSuperAdmin" />
      </UFormGroup>
      <div class="mt-6 flex gap-3">
        <UButton type="submit" :loading="submitting">ذخیره</UButton>
        <UButton variant="outline" to="/admins">انصراف</UButton>
      </div>
    </UForm>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresPermission: "users:create",
});

const schema = {
  firstName: { type: "string", required: true },
  lastName: { type: "string", required: true },
  email: { type: "string", required: true },
  username: { type: "string", required: true },
  password: { type: "string", required: true },
};

const state = reactive({
  firstName: "",
  lastName: "",
  email: "",
  username: "",
  password: "",
  contactNumber: "",
  role: "",
  isSuperAdmin: false,
  hotelIds: [] as string[],
});

const submitting = ref(false);

async function onSubmit() {
  submitting.value = true;
  try {
    await $fetch("/api/admins", { method: "POST", body: state });
    navigateTo("/admins");
  } catch (e) {
    console.error(e);
  } finally {
    submitting.value = false;
  }
}
</script>
