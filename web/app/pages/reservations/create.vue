<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/reservations" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        Back to Reservations
      </UButton>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">New Reservation</h1>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- Main Form -->
      <div class="lg:col-span-2">
        <UCard>
          <form @submit.prevent="handleSubmit">
            <div class="space-y-6">
              <!-- Guest Selection -->
              <div>
                <UFormField label="Guest" name="guestId" required>
                  <USelect
                    v-model="form.guestId"
                    :items="guests"
                    :loading="loadingGuests"
                    placeholder="Select a guest"
                    option-label="fullName"
                    option-value="id"
                    searchable
                    @update:searching="fetchGuests"
                  />
                </UFormField>
                <UButton variant="link" size="sm" class="mt-1 px-0" @click="showGuestModal = true">
                  <UIcon name="i-lucide-plus" class="mr-1" />
                  Add New Guest
                </UButton>
              </div>

              <!-- Room Selection -->
              <div>
                <UFormField label="Room" name="roomId" required>
                  <USelect
                    v-model="form.roomId"
                    :items="availableRooms"
                    :loading="loadingRooms"
                    placeholder="Select a room"
                    option-label="label"
                    option-value="id"
                    @change="updatePricing"
                  />
                </UFormField>
              </div>

              <!-- Dates -->
              <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <UFormField label="Check-in Date" name="checkInDate" required>
                  <UInput
                    v-model="form.checkInDate"
                    type="date"
                    :min="minDate"
                    @change="updatePricing"
                  />
                </UFormField>
                <UFormField label="Check-out Date" name="checkOutDate" required>
                  <UInput
                    v-model="form.checkOutDate"
                    type="date"
                    :min="form.checkInDate || minDate"
                    @change="updatePricing"
                  />
                </UFormField>
              </div>

              <!-- Number of Guests -->
              <UFormField label="Number of Guests" name="numberOfGuests" required>
                <UInput
                  v-model.number="form.numberOfGuests"
                  type="number"
                  min="1"
                  :max="selectedRoom?.capacity || 1"
                />
              </UFormField>

              <!-- Pricing -->
              <div class="rounded-lg bg-gray-50 p-4 dark:bg-gray-800">
                <div class="flex items-center justify-between">
                  <span class="text-gray-600 dark:text-gray-400">Estimated Total</span>
                  <span class="text-2xl font-bold">${{ estimatedTotal.toFixed(2) }}</span>
                </div>
                <p v-if="numberOfNights > 0" class="mt-1 text-sm text-gray-500">
                  {{ numberOfNights }} night{{ numberOfNights > 1 ? "s" : "" }} × ${{
                    selectedRoom?.basePrice?.toFixed(2)
                  }}/night
                </p>
              </div>

              <!-- Payment -->
              <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <UFormField label="Amount Paid" name="paidAmount">
                  <UInput
                    v-model.number="form.paidAmount"
                    type="number"
                    min="0"
                    step="0.01"
                    @change="updatePaymentStatus"
                  />
                </UFormField>
                <UFormField label="Payment Status" name="paymentStatus">
                  <USelect v-model="form.paymentStatus" :items="paymentStatusOptions" disabled />
                </UFormField>
              </div>

              <!-- Special Requests -->
              <UFormField label="Special Requests" name="specialRequests">
                <UTextarea
                  v-model="form.specialRequests"
                  placeholder="Any special requests..."
                  :rows="3"
                />
              </UFormField>

              <!-- Actions -->
              <div class="flex justify-end gap-3">
                <UButton variant="outline" to="/reservations" :disabled="loading"> Cancel </UButton>
                <UButton type="submit" color="primary" :loading="loading">
                  Create Reservation
                </UButton>
              </div>
            </div>
          </form>
        </UCard>
      </div>

      <!-- Sidebar - Room Preview -->
      <div>
        <UCard v-if="selectedRoom">
          <template #header>
            <h3 class="text-lg font-semibold">Room Details</h3>
          </template>
          <div class="space-y-3">
            <div>
              <p class="text-sm text-gray-500">Room Number</p>
              <p class="font-medium">{{ selectedRoom.roomNumber }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Type</p>
              <p class="font-medium capitalize">{{ selectedRoom.roomType }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Floor</p>
              <p class="font-medium">{{ selectedRoom.floor }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Capacity</p>
              <p class="font-medium">{{ selectedRoom.capacity }} guests</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Price per Night</p>
              <p class="font-medium">${{ selectedRoom.basePrice?.toFixed(2) }}</p>
            </div>
            <div v-if="selectedRoom.amenities">
              <p class="text-sm text-gray-500">Amenities</p>
              <div class="mt-1 flex flex-wrap gap-1">
                <UBadge
                  v-for="amenity in parseAmenities(selectedRoom.amenities)"
                  :key="amenity"
                  variant="soft"
                  size="sm"
                >
                  {{ amenity }}
                </UBadge>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <!-- New Guest Modal -->
    <UModal v-model:open="showGuestModal" title="Add New Guest">
      <template #content>
        <form @submit.prevent="createGuest" class="p-4">
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <UFormField label="First Name" name="newGuest.firstName" required>
                <UInput v-model="newGuest.firstName" placeholder="John" />
              </UFormField>
              <UFormField label="Last Name" name="newGuest.lastName" required>
                <UInput v-model="newGuest.lastName" placeholder="Doe" />
              </UFormField>
            </div>
            <UFormField label="Email" name="newGuest.email">
              <UInput v-model="newGuest.email" type="email" placeholder="john@example.com" />
            </UFormField>
            <UFormField label="Phone" name="newGuest.phone">
              <UInput v-model="newGuest.phone" placeholder="+1 234 567 8900" />
            </UFormField>
            <div class="grid grid-cols-2 gap-4">
              <UFormField label="ID Type" name="newGuest.idType">
                <USelect v-model="newGuest.idType" :items="idTypeOptions" placeholder="Select" />
              </UFormField>
              <UFormField label="ID Number" name="newGuest.idNumber">
                <UInput v-model="newGuest.idNumber" placeholder="ID number" />
              </UFormField>
            </div>
            <UFormField label="Address" name="newGuest.address">
              <UInput v-model="newGuest.address" placeholder="Street address" />
            </UFormField>
            <div class="grid grid-cols-2 gap-4">
              <UFormField label="City" name="newGuest.city">
                <UInput v-model="newGuest.city" placeholder="City" />
              </UFormField>
              <UFormField label="Country" name="newGuest.country">
                <UInput v-model="newGuest.country" placeholder="Country" />
              </UFormField>
            </div>
          </div>
          <div class="mt-6 flex justify-end gap-3">
            <UButton variant="outline" type="button" @click="showGuestModal = false">
              Cancel
            </UButton>
            <UButton type="submit" color="primary" :loading="creatingGuest"> Create Guest </UButton>
          </div>
        </form>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

interface Guest {
  id: number;
  firstName: string;
  lastName: string;
  fullName: string;
  email: string | null;
  phone: string | null;
}

interface Room {
  id: number;
  roomNumber: string;
  roomType: string;
  floor: number | null;
  capacity: number;
  basePrice: number;
  amenities: string | null;
}

const loading = ref(false);
const loadingGuests = ref(false);
const loadingRooms = ref(false);
const showGuestModal = ref(false);
const creatingGuest = ref(false);

const guests = ref<Guest[]>([]);
const rooms = ref<Room[]>([]);

const form = reactive({
  guestId: undefined as number | undefined,
  roomId: undefined as number | undefined,
  checkInDate: "",
  checkOutDate: "",
  numberOfGuests: 1,
  paidAmount: 0,
  paymentStatus: "pending" as "pending" | "partial" | "paid",
  specialRequests: "",
});

const newGuest = reactive({
  firstName: "",
  lastName: "",
  email: "",
  phone: "",
  idType: undefined as string | undefined,
  idNumber: "",
  address: "",
  city: "",
  country: "",
});

const idTypeOptions = [
  { value: "passport", label: "Passport" },
  { value: "national_id", label: "National ID" },
  { value: "driver_license", label: "Driver License" },
  { value: "other", label: "Other" },
];

const paymentStatusOptions = [
  { value: "pending", label: "Pending" },
  { value: "partial", label: "Partial" },
  { value: "paid", label: "Paid" },
];

const minDate = computed(() => {
  const today = new Date();
  return today.toISOString().split("T")[0];
});

const selectedRoom = computed(() => {
  return rooms.value.find((r) => r.id === form.roomId);
});

const numberOfNights = computed(() => {
  if (!form.checkInDate || !form.checkOutDate) return 0;
  const checkIn = new Date(form.checkInDate);
  const checkOut = new Date(form.checkOutDate);
  const diff = checkOut.getTime() - checkIn.getTime();
  return Math.max(0, Math.ceil(diff / (1000 * 60 * 60 * 24)));
});

const estimatedTotal = computed(() => {
  if (!selectedRoom.value || numberOfNights.value <= 0) return 0;
  return selectedRoom.value.basePrice * numberOfNights.value;
});

const availableRooms = computed(() => {
  return rooms.value.map((room) => ({
    ...room,
    label: `${room.roomNumber} - ${room.roomType} ($${room.basePrice}/night)`,
  }));
});

const fetchGuests = async (search?: string) => {
  loadingGuests.value = true;
  try {
    const params = new URLSearchParams();
    if (search) params.append("search", search);
    params.append("limit", "50");

    const response = await $fetch<{ data: any[] }>(`/api/guests?${params.toString()}`);
    if (response.data) {
      guests.value = response.data.map((g: any) => ({
        id: g.id,
        firstName: g.firstName || "",
        lastName: g.lastName || "",
        fullName: `${g.firstName || ""} ${g.lastName || ""}`,
        email: g.email,
        phone: g.phone,
      }));
    }
  } catch (error) {
    console.error("Failed to fetch guests:", error);
  } finally {
    loadingGuests.value = false;
  }
};

const fetchRooms = async () => {
  loadingRooms.value = true;
  try {
    const response = await $fetch("/api/rooms?status=available&limit=100");
    if (response.data) {
      rooms.value = response.data;
    }
  } catch (error) {
    console.error("Failed to fetch rooms:", error);
  } finally {
    loadingRooms.value = false;
  }
};

const updatePricing = () => {
  updatePaymentStatus();
};

const updatePaymentStatus = () => {
  if (form.paidAmount >= estimatedTotal.value && estimatedTotal.value > 0) {
    form.paymentStatus = "paid";
  } else if (form.paidAmount > 0) {
    form.paymentStatus = "partial";
  } else {
    form.paymentStatus = "pending";
  }
};

const parseAmenities = (amenities: string | undefined) => {
  if (!amenities) return [];
  try {
    return JSON.parse(amenities);
  } catch {
    return amenities.split(",").map((a) => a.trim());
  }
};

const handleSubmit = async () => {
  if (!form.guestId || !form.roomId || !form.checkInDate || !form.checkOutDate) {
    return;
  }

  loading.value = true;
  try {
    await $fetch("/api/reservations", {
      method: "POST",
      body: {
        guestId: form.guestId,
        roomId: form.roomId,
        checkInDate: form.checkInDate,
        checkOutDate: form.checkOutDate,
        numberOfGuests: form.numberOfGuests,
        specialRequests: form.specialRequests,
        totalAmount: estimatedTotal.value,
        paidAmount: form.paidAmount,
      },
    });

    await navigateTo("/reservations");
  } catch (error) {
    console.error("Failed to create reservation:", error);
  } finally {
    loading.value = false;
  }
};

const createGuest = async () => {
  creatingGuest.value = true;
  try {
    const response = await $fetch("/api/guests", {
      method: "POST",
      body: newGuest,
    });

    if (response.data?.id) {
      form.guestId = response.data.id;
    }
    await fetchGuests();
    showGuestModal.value = false;

    // Reset form
    newGuest.firstName = "";
    newGuest.lastName = "";
    newGuest.email = "";
    newGuest.phone = "";
    newGuest.idType = undefined;
    newGuest.idNumber = "";
    newGuest.address = "";
    newGuest.city = "";
    newGuest.country = "";
  } catch (error) {
    console.error("Failed to create guest:", error);
  } finally {
    creatingGuest.value = false;
  }
};

onMounted(() => {
  fetchGuests();
  fetchRooms();
});
</script>
