<template>
  <div class="flex h-full flex-1 grow flex-col">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ t("sana.title") }}</h1>

      <UModal v-model:open="syncAllModalOpen">
        <UButton color="primary">
          <UIcon name="i-lucide-refresh-cw" class="mr-2" />
          {{ t("sana.syncAll") }}
        </UButton>
        <template #header>
          <h2 class="text-lg font-semibold">{{ t("sana.syncAllTitle") }}</h2>
        </template>
        <template #body>
          <p>{{ t("sana.syncAllDescription") }}</p>
        </template>
        <template #footer>
          <div class="flex justify-end gap-3">
            <UButton variant="outline" @click="syncAllModalOpen = false">
              {{ t("actions.cancel") }}
            </UButton>
            <UButton color="primary" :loading="syncingAll" @click="syncAll">
              {{ t("sana.syncAll") }}
            </UButton>
          </div>
        </template>
      </UModal>
    </div>

    <div class="flex flex-1 grow flex-col gap-4">
      <!-- Guests Table -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-lg font-semibold">{{ t("sana.guests") }}</span>
            <UBadge color="primary">{{ guests?.length || 0 }}</UBadge>
          </div>
        </template>

        <UTable :data="guests" :columns="guestColumns" :loading="guestsLoading" striped>
          <template #status-cell="{ row }">
            <UBadge :color="row.original.syncTime ? 'success' : 'warning'" variant="soft">
              {{ row.original.syncTime ? t("sana.synced") : t("sana.notSynced") }}
            </UBadge>
          </template>

          <template #syncTime-cell="{ row }">
            {{ row.original.syncTime ? formatDate(row.original.syncTime) : "-" }}
          </template>

          <template #recordMosafer-cell="{ row }">
            {{ row.original.recordMosafer || "-" }}
          </template>

          <template #shomarePaziresh-cell="{ row }">
            {{ row.original.shomarePaziresh || "-" }}
          </template>

          <template #actions-cell="{ row }">
            <UButton
              v-if="!row.original.syncTime"
              variant="ghost"
              size="sm"
              color="primary"
              :loading="syncingGuestId === row.original.id"
              @click="syncGuest(row.original)"
            >
              <UIcon name="i-lucide-refresh-cw" class="h-4 w-4" />
            </UButton>
          </template>
        </UTable>
      </UCard>

      <!-- Rooms Table -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-lg font-semibold">{{ t("sana.rooms") }}</span>
            <div class="flex items-center gap-3">
              <UBadge color="primary">{{ rooms?.length || 0 }}</UBadge>
              <UButton
                color="primary"
                variant="soft"
                size="sm"
                :loading="syncingRooms"
                @click="syncRooms"
              >
                <UIcon name="i-lucide-refresh-cw" class="mr-1 h-4 w-4" />
                {{ t("sana.syncRooms") }}
              </UButton>
            </div>
          </div>
        </template>

        <UTable :data="rooms" :columns="roomColumns" :loading="roomsLoading" striped>
          <template #roomNumber-cell="{ row }">
            <NuxtLink
              :to="`/rooms/${row.original.room?.id}`"
              class="font-medium text-primary hover:underline"
            >
              {{ row.original.room?.roomNumber || "-" }}
            </NuxtLink>
          </template>

          <template #status-cell="{ row }">
            <UBadge :color="row.original.isSynced ? 'success' : 'warning'" variant="soft">
              {{ row.original.isSynced ? t("sana.synced") : t("sana.notSynced") }}
            </UBadge>
          </template>

          <template #lastSyncTime-cell="{ row }">
            {{ row.original.lastSyncTime ? formatDate(row.original.lastSyncTime) : "-" }}
          </template>

          <template #lastError-cell="{ row }">
            <span v-if="row.original.lastError" class="text-red-500">{{
              row.original.lastError
            }}</span>
            <span v-else class="text-gray-400">-</span>
          </template>

          <template #actions-cell="{ row }">
            <span></span>
          </template>
        </UTable>
      </UCard>
    </div>

    <!-- Error Modal -->
    <UModal v-model:open="errorModalOpen">
      <template #header>
        <h2 class="text-lg font-semibold text-red-600">{{ errorModalTitle }}</h2>
      </template>
      <template #body>
        <div class="text-sm whitespace-pre-wrap text-gray-700 dark:text-gray-300">
          {{ errorModalMessage }}
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end">
          <UButton color="primary" @click="errorModalOpen = false">
            {{ t("actions.ok") }}
          </UButton>
        </div>
      </template>
    </UModal>
  </div>
</template>

<script setup lang="ts">
import type { TableColumn } from "@nuxt/ui";
import { useI18n } from "vue-i18n";
import { postApiSanaGuestsIdSync, postApiSanaSyncAll } from "~/utils/client";

const { t } = useI18n();

interface SanaGuest {
  id: number;
  guestId: number;
  guest?: {
    firstName: string;
    lastName: string;
    nationalId?: string;
  };
  recordMosafer: number;
  shomarePaziresh: string;
  shomareOtagh: string;
  syncTime?: string;
}

interface SanaRoomRack {
  id: number;
  hotelId: string;
  room?: {
    id: number;
    roomNumber: string;
  };
  rac: string;
  lastSyncTime?: string;
  isSynced: boolean;
  lastError?: string;
}

interface SanaSyncAllResult {
  status: string;
  guestsSynced: number;
  guestsFailed: number;
  roomsSynced: number;
  roomsFailed: number;
  errors?: string[];
}

interface SanaSyncRoomsResult {
  status: string;
  roomsSynced: number;
  roomsFailed: number;
  error?: string;
}

const guestColumns: TableColumn<SanaGuest>[] = [
  { accessorKey: "guest.firstName", header: t("sana.columns.name") },
  { accessorKey: "guest.lastName", header: t("sana.columns.lastName") },
  { accessorKey: "guest.nationalId", header: t("sana.columns.nationalId") },
  { accessorKey: "recordMosafer", header: t("sana.columns.recordMosafer") },
  { accessorKey: "shomarePaziresh", header: t("sana.columns.shomarePaziresh") },
  { accessorKey: "status", header: t("sana.columns.status") },
  { accessorKey: "syncTime", header: t("sana.columns.syncTime") },
  { accessorKey: "actions", header: t("sana.columns.actions") },
];

const roomColumns: TableColumn<SanaRoomRack>[] = [
  { accessorKey: "roomNumber", header: t("sana.columns.roomNumber") },
  { accessorKey: "rac", header: t("sana.columns.rac") },
  { accessorKey: "status", header: t("sana.columns.status") },
  { accessorKey: "lastSyncTime", header: t("sana.columns.lastSyncTime") },
  { accessorKey: "lastError", header: t("sana.columns.lastError") },
  { accessorKey: "actions", header: t("sana.columns.actions") },
];

const syncAllModalOpen = ref(false);
const syncingAll = ref(false);
const syncingRooms = ref(false);
const syncingGuestId = ref<number | null>(null);

const errorModalOpen = ref(false);
const errorModalTitle = ref("");
const errorModalMessage = ref("");

const showError = (title: string, message: string) => {
  errorModalTitle.value = title;
  errorModalMessage.value = message;
  errorModalOpen.value = true;
};

const extractErrorMessage = (error: any): string => {
  if (typeof error === "string") return error;
  if (error?.data?.error) return error.data.error;
  if (error?.data?.detail) return error.data.detail;
  if (error?.data?.message) return error.data.message;
  if (error?.message) return error.message;
  return t("sana.unknownError");
};

const {
  data: guests,
  pending: guestsLoading,
  refresh: refreshGuests,
} = useAsyncData("sana-guests", () => $fetch<SanaGuest[]>("/api/sana/guests"));

const {
  data: rooms,
  pending: roomsLoading,
  refresh: refreshRooms,
} = useAsyncData("sana-rooms", () => $fetch<SanaRoomRack[]>("/api/sana/rooms"));

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("fa-IR", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
};

const syncGuest = async (guest: SanaGuest) => {
  syncingGuestId.value = guest.id;
  try {
    await postApiSanaGuestsIdSync({ path: { id: guest.id.toString() } });
    await refreshGuests();
  } catch (error: any) {
    const msg = extractErrorMessage(error);
    showError(t("sana.syncGuestError"), msg);
  } finally {
    syncingGuestId.value = null;
  }
};


const syncRooms = async () => {
  syncingRooms.value = true;
  try {
    const result = await $fetch<SanaSyncRoomsResult>("/api/sana/sync-rooms", {
      method: "POST",
    });
    await refreshRooms();
    if (result.error) {
      showError(t("sana.syncRoomsError"), result.error);
    }
  } catch (error: any) {
    const msg = extractErrorMessage(error);
    showError(t("sana.syncRoomsError"), msg);
  } finally {
    syncingRooms.value = false;
  }
};

const syncAll = async () => {
  syncingAll.value = true;
  try {
    const result = (await postApiSanaSyncAll({
      requestValidator: undefined,
    })) as unknown as SanaSyncAllResult;

    await refreshGuests();
    await refreshRooms();
    syncAllModalOpen.value = false;

    if (result.errors && result.errors.length > 0) {
      const summary = `${t("sana.syncAllPartial")}\n\n${result.errors.join("\n")}`;
      showError(t("sana.syncAllError"), summary);
    }
  } catch (error: any) {
    const msg = extractErrorMessage(error);
    showError(t("sana.syncAllError"), msg);
  } finally {
    syncingAll.value = false;
  }
};
</script>
