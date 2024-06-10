<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { getLoggedInUser } from '@/common/user-service';
import { useUserStore } from '@/common/store';

const menubarItems = ref([]);

const usermenuItems = ref([
  {
    label: 'Profile',
    items: [
      {
        label: 'Settings',
        icon: 'pi pi-cog',
        command: () => {
          settingsDialogVisible.value = true;
        }
      },
      {
        label: 'Logout',
        icon: 'pi pi-sign-out',
        command: () => {
          logout();
        }
      }
    ]
  }
]);

const menu = ref();

const settingsDialogVisible = ref(false);

const userStore = useUserStore();

onMounted(async () => {
  const userInfo = await getLoggedInUser();
  console.log('user', userInfo);
  if (userInfo !== null) {
    userStore.loggedIn = true;
    userStore.email = userInfo.userDetails;
  }
});

const toggle = (event: any) => {
  menu.value.toggle(event);
};

async function login() {
  window.open('/.auth/login/aad', '_self');
}

function logout() {
  userStore.loggedIn = false;
  userStore.email = null;
  // TODO aad action?
}
</script>

<template>
  <vueMenubar :model="menubarItems">
    <template #start>
      <svg
        width="40"
        height="40"
        viewBox="0 0 35 40"
        fill="yellow"
        xmlns="http://www.w3.org/2000/svg"
        class="h-2rem"
      >
        <text x="0" y="15" fill="var(--primary-color)" stroke="var(--primary-color)" font-size="15">
          GeWoScout
        </text>
        <text x="0" y="40" fill="var(--text-color)" stroke="var(--text-color)" font-size="15">
          Scout
        </text>
      </svg>
    </template>
    <template #item="{ item, props }">
      <a v-ripple class="flex align-items-center" v-bind="props.action">
        <span :class="item.icon" />
        <span class="ml-2">{{ item.label }}</span>
      </a>
    </template>
    <template #end>
      <div class="flex align-items-center gap-2">
        <vueAvatar
          icon="pi pi-user"
          shape="circle"
          size="large"
          class="cursor-pointer"
          v-if="userStore.loggedIn"
          @click="toggle"
        />
        <vueButton label="Login" @click="login()" v-else></vueButton>

        <vueMenu ref="menu" id="overlay_menu" :model="usermenuItems" :popup="true" />

        <vueDialog
          v-model:visible="settingsDialogVisible"
          modal
          header="Edit Preferences"
          :style="{ width: '25rem' }"
        >
          <span class="p-text-secondary block mb-5">Update your information.</span>
          <div class="flex align-items-center gap-3 mb-5">
            <label for="email" class="font-semibold w-6rem">Email</label>
            <vueInputText id="email" class="flex-auto" autocomplete="off" />
          </div>
          <!-- Add other preferences to edit -->
          <div class="flex justify-content-end gap-2">
            <vueButton
              type="button"
              label="Cancel"
              severity="secondary"
              @click="settingsDialogVisible = false"
            ></vueButton>
            <vueButton
              type="button"
              label="Save"
              @click="settingsDialogVisible = false"
            ></vueButton>
          </div>
        </vueDialog>
      </div>
    </template>
  </vueMenubar>
</template>

<style scoped></style>
