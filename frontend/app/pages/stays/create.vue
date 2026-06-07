<template>
  <div class="p-6 max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold mb-6">ثبت پذیرش جدید</h1>
    <UForm :state="state" @submit="onSubmit">
      <div class="grid grid-cols-2 gap-4">
        <UFormGroup label="مهمان" name="guestId">
          <USelect v-model="state.guestId" :options="guestOptions" />
        </UFormGroup>
        <UFormGroup label="اتاق" name="roomId">
          <USelect v-model="state.roomId" :options="roomOptions" />
        </UFormGroup>
        <UFormGroup label="تاریخ ورود" name="entryDate">
          <UInput v-model="state.entryDate" type="datetime-local" />
        </UFormGroup>
        <UFormGroup label="تاریخ خروج" name="departureDate">
          <UInput v-model="state.departureDate" type="datetime-local" />
        </UFormGroup>
        <UFormGroup label="تعداد نفرات" name="numberOfPeople">
          <UInput v-model.number="state.numberOfPeople" type="number" />
        </UFormGroup>
        <UFormGroup label="قیمت اتاق" name="roomPrice">
          <UInput v-model.number="state.roomPrice" type="number" />
        </UFormGroup>
      </div>
      <div class="grid grid-cols-4 gap-4 mt-4">
        <UFormGroup label="صبحانه" name="breakfast">
          <UToggle v-model="state.breakfast" />
        </UFormGroup>
        <UFormGroup label="نیم‌برد" name="halfBoard">
          <UToggle v-model="state.halfBoard" />
        </UFormGroup>
        <UFormGroup label="فول‌برد" name="fullBoard">
          <UToggle v-model="state.fullBoard" />
        </UFormGroup>
        <UFormGroup label="پارکینگ" name="parking">
          <UToggle v-model="state.parking" />
        </UFormGroup>
      </div>
      <UFormGroup label="یادداشت" name="notes" class="mt-4">
        <UTextarea v-model="state.notes" />
      </UFormGroup>
      <div class="mt-6 flex gap-3">
        <UButton type="submit" :loading="submitting">ثبت پذیرش</UButton>
        <UButton variant="outline" to="/stays">انصراف</UButton>
      </div>
    </UForm>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresPermission: "guests:create",
});

const state = reactive({
  hotelId: "",
  guestId: undefined as number | undefined,
  roomId: undefined as number | undefined,
  entryDate: "",
  departureDate: "",
  numberOfPeople: 1,
  roomPrice: 0,
  breakfast: false,
  halfBoard: false,
  fullBoard: false,
  parking: false,
  notes: "",
});

const { data: guests } = useFetch("/api/guests", { transform: (res: any) => res?.data ?? [] });
const { data: rooms } = useFetch("/api/rooms", { transform: (res: any) => res?.data ?? [] });

const guestOptions = computed(() =>
  (guests.value ?? []).map((g: any) => ({ label: `${g.firstName} ${g.lastName}`, value: g.id }))
);
const roomOptions = computed(() =>
  (rooms.value ?? []).map((r: any) => ({ label: r.roomNumber, value: r.id }))
);

const authStore = useAuthStore();
watchEffect(() => {
  if (authStore.currentHotelId) {
    state.hotelId = authStore.currentHotelId;
  }
});

const submitting = ref(false);
async function onSubmit() {
  submitting.value = true;
  try {
    await $fetch("/api/stays", { method: "POST", body: state });
    navigateTo("/stays");
  } catch (e) {
    console.error(e);
  } finally {
    submitting.value = false;
  }
}
</script>
