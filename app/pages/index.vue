<template>
  <div>
    <h1 class="mb-6 text-3xl font-bold text-gray-900 dark:text-white">Dashboard</h1>

    <div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Total Rooms</span>
            <UIcon name="i-lucide-bed" class="h-5 w-5 text-gray-400" />
          </div>
        </template>
        <div class="text-3xl font-bold">120</div>
        <p class="mt-1 text-sm text-gray-500">Available: 85</p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400">Occupancy</span>
            <UIcon name="i-lucide-users" class="h-5 w-5 text-gray-400" />
          </div>
        </template>
        <div class="text-3xl font-bold">72%</div>
        <p class="mt-1 text-sm text-gray-500">86 rooms occupied</p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400"
              >Today's Revenue</span
            >
            <UIcon name="i-lucide-dollar-sign" class="h-5 w-5 text-gray-400" />
          </div>
        </template>
        <div class="text-3xl font-bold">$12,450</div>
        <p class="mt-1 text-sm text-green-500">+8% from yesterday</p>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400"
              >Check-ins Today</span
            >
            <UIcon name="i-lucide-log-in" class="h-5 w-5 text-gray-400" />
          </div>
        </template>
        <div class="text-3xl font-bold">18</div>
        <p class="mt-1 text-sm text-gray-500">Check-outs: 12</p>
      </UCard>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold">Recent Reservations</h2>
            <UButton variant="ghost" size="sm" to="/reservations">View All</UButton>
          </div>
        </template>
        <div class="space-y-4">
          <div
            v-for="i in 5"
            :key="i"
            class="flex items-center justify-between border-b border-gray-100 py-2 last:border-0 dark:border-gray-800"
          >
            <div>
              <p class="font-medium">John Doe</p>
              <p class="text-sm text-gray-500">Room {{ 100 + i }}</p>
            </div>
            <div class="text-right">
              <UBadge color="success" variant="soft">Confirmed</UBadge>
              <p class="mt-1 text-xs text-gray-500">{{ new Date().toLocaleDateString() }}</p>
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold">Quick Actions</h2>
          </div>
        </template>
        <div class="grid grid-cols-2 gap-4">
          <UButton to="/reservations/create" block size="lg" color="primary">
            <UIcon name="i-lucide-plus" class="mr-2" />
            New Reservation
          </UButton>
          <UButton to="/reservations/check-in" block size="lg" color="neutral">
            <UIcon name="i-lucide-log-in" class="mr-2" />
            Check-in Guest
          </UButton>
          <UButton to="/rooms" block size="lg" color="neutral">
            <UIcon name="i-lucide-bed" class="mr-2" />
            Manage Rooms
          </UButton>
          <UButton to="/parking/check-in" block size="lg" color="neutral">
            <UIcon name="i-lucide-car" class="mr-2" />
            Parking Check-in
          </UButton>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin", "manager", "receptionist", "staff"],
});
</script>
