<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <UButton to="/users" variant="ghost" size="sm" class="mb-2">
          <UIcon name="i-lucide-arrow-left" class="mr-1" />
          Back to Users
        </UButton>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Create User</h1>
      </div>
    </div>

    <UCard>
      <form @submit.prevent="createUser">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-sm font-medium">Email *</label>
            <UInput v-model="form.email" type="email" placeholder="user@example.com" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Password *</label>
            <UInput
              v-model="form.password"
              type="password"
              placeholder="Minimum 6 characters"
              required
            />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">First Name *</label>
            <UInput v-model="form.firstName" placeholder="John" required />
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Last Name *</label>
            <UInput v-model="form.lastName" placeholder="Doe" required />
          </div>

          <div class="md:col-span-2">
            <label class="mb-1 block text-sm font-medium">Roles *</label>
            <div class="flex flex-wrap gap-4">
              <UCheckbox
                v-for="role in roleOptions"
                :key="role.value"
                :model-value="form.roles.includes(role.value)"
                :label="role.label"
                @update:model-value="toggleRole(role.value, $event)"
              />
            </div>
            <p v-if="form.roles.length === 0" class="mt-1 text-sm text-red-500">
              At least one role is required
            </p>
          </div>

          <div>
            <label class="mb-1 block text-sm font-medium">Active</label>
            <UCheckbox v-model="form.isActive" label="User account is active" />
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <UButton type="button" variant="outline" to="/users">Cancel</UButton>
          <UButton
            type="submit"
            color="primary"
            :loading="loading"
            :disabled="form.roles.length === 0"
          >
            Create User
          </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  requiresRole: ["admin"],
});

const form = reactive({
  email: "",
  password: "",
  firstName: "",
  lastName: "",
  roles: [] as string[],
  isActive: true,
});

const roleOptions = [
  { value: "admin", label: "Admin" },
  { value: "manager", label: "Manager" },
  { value: "receptionist", label: "Receptionist" },
  { value: "staff", label: "Staff" },
];

const toggleRole = (role: string, checked: boolean | string) => {
  const isChecked = !!checked;
  if (isChecked) {
    if (!form.roles.includes(role)) {
      form.roles.push(role);
    }
  } else {
    const index = form.roles.indexOf(role);
    if (index > -1) {
      form.roles.splice(index, 1);
    }
  }
};

const loading = ref(false);
const router = useRouter();

const createUser = async () => {
  if (form.roles.length === 0) return;

  loading.value = true;
  try {
    await $fetch("/api/users", {
      method: "POST",
      body: {
        email: form.email,
        password: form.password,
        firstName: form.firstName,
        lastName: form.lastName,
        roles: form.roles,
        isActive: form.isActive,
      },
    });
    router.push("/users");
  } catch (error: any) {
    console.error("Failed to create user:", error);
    alert(error.data?.message || "Failed to create user");
  } finally {
    loading.value = false;
  }
};
</script>
