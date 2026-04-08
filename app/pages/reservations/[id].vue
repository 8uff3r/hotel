<template>
  <div>
    <div class="mb-6">
      <UButton variant="ghost" to="/reservations" class="mb-4">
        <UIcon name="i-lucide-arrow-left" class="mr-2" />
        Back to Reservations
      </UButton>
      <div class="flex items-center justify-between">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          Reservation #{{ reservation?.id }}
        </h1>
        <div class="flex items-center gap-2">
          <UBadge
            v-if="reservation"
            :color="getStatusColor(reservation.status)"
            variant="solid"
            size="lg"
          >
            {{ reservation.status }}
          </UBadge>
        </div>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <ULoader size="lg" />
    </div>

    <div v-else-if="reservation" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <!-- Main Details -->
      <div class="space-y-6 lg:col-span-2">
        <!-- Guest & Room Info -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">Reservation Details</h3>
          </template>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <!-- Guest Info -->
            <div>
              <p class="mb-1 text-sm text-gray-500">Guest</p>
              <p class="text-lg font-medium">
                {{ reservation.guest?.firstName }} {{ reservation.guest?.lastName }}
              </p>
              <p class="text-gray-600">{{ reservation.guest?.email }}</p>
              <p class="text-gray-600">{{ reservation.guest?.phone }}</p>
              <p v-if="reservation.guest?.address" class="text-gray-600">
                {{ reservation.guest?.address }}, {{ reservation.guest?.city }},
                {{ reservation.guest?.country }}
              </p>
            </div>

            <!-- Room Info -->
            <div>
              <p class="mb-1 text-sm text-gray-500">Room</p>
              <p class="text-lg font-medium">{{ reservation.room?.roomNumber }}</p>
              <p class="text-gray-600 capitalize">{{ reservation.room?.roomType }}</p>
              <p class="text-gray-600">Floor {{ reservation.room?.floor }}</p>
              <p class="text-gray-600">Capacity: {{ reservation.room?.capacity }} guests</p>
            </div>
          </div>
        </UCard>

        <!-- Dates & Status -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">Stay Information</h3>
          </template>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div>
              <p class="text-sm text-gray-500">Check-in Date</p>
              <p class="font-medium">{{ formatDate(reservation.checkInDate) }}</p>
              <p v-if="reservation.actualCheckIn" class="text-sm text-green-600">
                <UIcon name="i-lucide-check-circle" class="mr-1 inline" />
                Checked in: {{ formatDateTime(reservation.actualCheckIn) }}
              </p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Check-out Date</p>
              <p class="font-medium">{{ formatDate(reservation.checkOutDate) }}</p>
              <p v-if="reservation.actualCheckOut" class="text-sm text-green-600">
                <UIcon name="i-lucide-check-circle" class="mr-1 inline" />
                Checked out: {{ formatDateTime(reservation.actualCheckOut) }}
              </p>
            </div>
          </div>

          <div class="mt-4 border-t pt-4">
            <p class="text-sm text-gray-500">Number of Guests</p>
            <p class="font-medium">{{ reservation.numberOfGuests }}</p>
          </div>

          <div v-if="reservation.specialRequests" class="mt-4 border-t pt-4">
            <p class="text-sm text-gray-500">Special Requests</p>
            <p class="mt-1">{{ reservation.specialRequests }}</p>
          </div>
        </UCard>

        <!-- Payment Info -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold">Payment</h3>
              <UBadge :color="getPaymentColor(reservation.paymentStatus)" variant="soft">
                {{ reservation.paymentStatus }}
              </UBadge>
            </div>
          </template>
          <div class="grid grid-cols-2 gap-6">
            <div>
              <p class="text-sm text-gray-500">Total Amount</p>
              <p class="text-2xl font-medium">${{ reservation.totalAmount?.toFixed(2) }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Amount Paid</p>
              <p class="text-2xl font-medium">${{ reservation.paidAmount?.toFixed(2) }}</p>
            </div>
          </div>
          <div
            v-if="reservation.totalAmount - reservation.paidAmount > 0"
            class="mt-4 border-t pt-4"
          >
            <p class="text-sm text-gray-500">Balance Due</p>
            <p class="text-xl font-medium text-red-600">
              ${{ (reservation.totalAmount - reservation.paidAmount).toFixed(2) }}
            </p>
          </div>
        </UCard>
      </div>

      <!-- Sidebar Actions -->
      <div class="space-y-6">
        <!-- Actions -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">Actions</h3>
          </template>
          <div class="space-y-3">
            <UButton
              v-if="reservation.status === 'confirmed'"
              color="success"
              block
              :loading="processing"
              @click="handleCheckIn"
            >
              <UIcon name="i-lucide-log-in" class="mr-2" />
              Check In Guest
            </UButton>
            <UButton
              v-if="reservation.status === 'checked_in'"
              color="warning"
              block
              :loading="processing"
              @click="handleCheckOut"
            >
              <UIcon name="i-lucide-log-out" class="mr-2" />
              Check Out Guest
            </UButton>
            <UButton
              v-if="reservation.status === 'confirmed'"
              color="error"
              variant="outline"
              block
              :loading="processing"
              @click="handleCancel"
            >
              <UIcon name="i-lucide-x" class="mr-2" />
              Cancel Reservation
            </UButton>
            <UButton
              v-if="reservation.status === 'no_show'"
              color="info"
              block
              :loading="processing"
              @click="handleCheckIn"
            >
              <UIcon name="i-lucide-user-check" class="mr-2" />
              Mark as Arrived
            </UButton>
          </div>
        </UCard>

        <!-- Room Preview -->
        <UCard v-if="reservation.room">
          <template #header>
            <h3 class="text-lg font-semibold">Room Details</h3>
          </template>
          <div class="space-y-3">
            <div>
              <p class="text-sm text-gray-500">Room Number</p>
              <p class="font-medium">{{ reservation.room.roomNumber }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Type</p>
              <p class="font-medium capitalize">{{ reservation.room.roomType }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Base Price</p>
              <p class="font-medium">${{ reservation.room.basePrice }}/night</p>
            </div>
            <div v-if="reservation.room.amenities">
              <p class="text-sm text-gray-500">Amenities</p>
              <div class="mt-1 flex flex-wrap gap-1">
                <UBadge
                  v-for="amenity in parseAmenities(reservation.room.amenities)"
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

        <!-- Meta Info -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold">Meta Information</h3>
          </template>
          <div class="space-y-2 text-sm">
            <p class="text-gray-500">Created: {{ formatDateTime(reservation.createdAt) }}</p>
            <p class="text-gray-500">Updated: {{ formatDateTime(reservation.updatedAt) }}</p>
          </div>
        </UCard>
      </div>
    </div>

    <!-- Not Found -->
    <UCard v-else>
      <div class="py-12 text-center">
        <p class="text-gray-500">Reservation not found</p>
        <UButton to="/reservations" class="mt-4">Back to Reservations</UButton>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist"],
});

interface Guest {
  id: number;
  firstName: string | null;
  lastName: string | null;
  email: string | null;
  phone: string | null;
  address: string | null;
  city: string | null;
  country: string | null;
}

interface Room {
  id: number;
  roomNumber: string | null;
  roomType: string | null;
  floor: number | null;
  capacity: number | null;
  basePrice: number | null;
  amenities: string | null;
}

interface Reservation {
  id: number;
  guestId: number;
  roomId: number;
  checkInDate: string;
  checkOutDate: string;
  actualCheckIn: string | null;
  actualCheckOut: string | null;
  status: string;
  totalAmount: number;
  paidAmount: number;
  paymentStatus: string;
  numberOfGuests: number;
  specialRequests: string | null;
  createdAt: string;
  updatedAt: string;
  guest?: Guest;
  room?: Room;
}

const route = useRoute();
const reservationId = Number(route.params.id);

const loading = ref(true);
const processing = ref(false);
const reservation = ref<Reservation | null>(null);

const fetchReservation = async () => {
  loading.value = true;
  try {
    const response = await $fetch(`/api/reservations/${reservationId}`);
    reservation.value = response.data;
  } catch (error) {
    console.error("Failed to fetch reservation:", error);
    reservation.value = null;
  } finally {
    loading.value = false;
  }
};

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
};

const formatDateTime = (date: string) => {
  return new Date(date).toLocaleString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const getStatusColor = (status: string): "success" | "warning" | "info" | "error" | "neutral" => {
  const colors: Record<string, "success" | "warning" | "info" | "error" | "neutral"> = {
    pending: "warning",
    confirmed: "info",
    checked_in: "success",
    checked_out: "neutral",
    cancelled: "error",
    no_show: "error",
  };
  return colors[status] || "neutral";
};

const getPaymentColor = (status: string): "success" | "warning" | "info" | "error" | "neutral" => {
  const colors: Record<string, "success" | "warning" | "info" | "error" | "neutral"> = {
    pending: "warning",
    partial: "info",
    paid: "success",
    refunded: "error",
  };
  return colors[status] || "neutral";
};

const parseAmenities = (amenities: string | undefined | null) => {
  if (!amenities) return [];
  try {
    return JSON.parse(amenities);
  } catch {
    return amenities.split(",").map((a: string) => a.trim());
  }
};

const handleCheckIn = async () => {
  processing.value = true;
  try {
    await $fetch(`/api/reservations/${reservationId}/check-in`, { method: "post" });
    await fetchReservation();
  } catch (error) {
    console.error("Failed to check in:", error);
  } finally {
    processing.value = false;
  }
};

const handleCheckOut = async () => {
  processing.value = true;
  try {
    await $fetch(`/api/reservations/${reservationId}/check-out`, { method: "post" });
    await fetchReservation();
  } catch (error) {
    console.error("Failed to check out:", error);
  } finally {
    processing.value = false;
  }
};

const handleCancel = async () => {
  processing.value = true;
  try {
    await $fetch(`/api/reservations/${reservationId}`, {
      method: "put",
      body: { status: "cancelled" },
    });
    await fetchReservation();
  } catch (error) {
    console.error("Failed to cancel:", error);
  } finally {
    processing.value = false;
  }
};

onMounted(fetchReservation);
</script>
