<script setup lang="ts">
import { onMounted, ref, type Ref } from 'vue';

import Divider from 'primevue/divider';
import Dropdown from 'primevue/dropdown';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import SelectButton from 'primevue/selectbutton';
import Calendar from 'primevue/calendar';

import { getLoggedInUser } from '@/common/user-service';
import { useUserStore } from '@/common/store';
import { getUserPreferences, setUserPreferences } from '@/common/api-service';
import { EnergyClass, Type } from '@/types/Enums';
import type UserPreferences from '@/types/UserPreferences';

const userStore = useUserStore();

const menubarItems = ref([]);

const usermenuItems = ref([
  {
    label: 'Profile',
    items: [
      {
        label: 'Settings',
        icon: 'pi pi-cog',
        command: () => {
          openSettingsDialog();
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

const usermenu = ref();

const settingsDialogVisible = ref(false);

const cities = ref([
  { name: 'Vienna', code: 'vienna' },
  { name: 'Graz', code: 'graz' }
]);

const genos = ref([
  { name: 'BWS-Gruppe', code: 'BWS-Gruppe' },
  { name: 'ÖVW', code: 'oevw' },
  { name: 'wbv-gpa', code: 'wbv-gpa' }
]);

const energyClasses = ref([
  { name: 'A++', value: EnergyClass['A++'] },
  { name: 'A+', value: EnergyClass['A+'] },
  { name: 'A', value: EnergyClass.A },
  { name: 'B', value: EnergyClass.B },
  { name: 'C', value: EnergyClass.C },
  { name: 'D', value: EnergyClass.D },
  { name: 'E', value: EnergyClass.E },
  { name: 'F', value: EnergyClass.F }
]);

const types = ref([
  { name: 'All', type: Type.both },
  { name: 'Rent', type: Type.rent },
  { name: 'Sale', type: Type.sale }
]);

const selectedCity = ref('vienna');

const userPreferences: Ref<UserPreferences> = ref({
  email: null,
  listingType: Type.both,
  city: 'vienna',
  housingCooperative: null,
  postalCode: null,
  minRoomCount: null,
  maxRoomCount: null,
  minSqm: null,
  maxSqm: null,
  availableFrom: null,
  minYearBuilt: null,
  maxYearBuilt: null,
  minHwgEnergyClass: EnergyClass.F,
  minFgeeEnergyClass: EnergyClass.F,
  minRentPrice: null,
  maxRentPrice: null,
  minCooperativeShare: null,
  maxCooperativeShare: null,
  minSalePrice: null,
  maxSalePrice: null
});

onMounted(async () => {
  const userInfo = await getLoggedInUser();
  console.log('user', userInfo);
  if (userInfo !== null) {
    userStore.loggedIn = true;
  }
});

const toggle = (event: any) => {
  usermenu.value.toggle(event);
};

async function openSettingsDialog() {
  settingsDialogVisible.value = true;

  const apiCallResponse = await getUserPreferences(selectedCity.value);

  if (apiCallResponse !== null) {
    userPreferences.value = apiCallResponse;
  } else {
    resetUserPreferences(selectedCity.value);
  }
}

async function getCityUserPreferences() {
  const apiCallResponse = await getUserPreferences(selectedCity.value);

  if (apiCallResponse !== null) {
    userPreferences.value = apiCallResponse;
  } else [resetUserPreferences(selectedCity.value)];
}

async function saveUserPreferences() {
  settingsDialogVisible.value = false;
  const response = await setUserPreferences(userPreferences.value);

  if (response == false) {
    console.log('ERROOOR need notification');
  }
}

function resetUserPreferences(city: string) {
  userPreferences.value = {
    email: null,
    listingType: Type.both,
    city: city,
    housingCooperative: null,
    postalCode: null,
    minRoomCount: null,
    maxRoomCount: null,
    minSqm: null,
    maxSqm: null,
    availableFrom: null,
    minYearBuilt: null,
    maxYearBuilt: null,
    minHwgEnergyClass: EnergyClass.F,
    minFgeeEnergyClass: EnergyClass.F,
    minRentPrice: null,
    maxRentPrice: null,
    minCooperativeShare: null,
    maxCooperativeShare: null,
    minSalePrice: null,
    maxSalePrice: null
  };
}

async function login() {
  window.open('/.auth/login/aad', '_self');
}

function logout() {
  userStore.loggedIn = false;
  window.open('/.auth/logout', '_self');
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
          v-if="userStore.loggedIn"
          icon="pi pi-user"
          shape="circle"
          size="large"
          class="cursor-pointer"
          @click="toggle"
        />
        <vueButton v-else label="Login" @click="login()"></vueButton>

        <vueMenu ref="usermenu" id="overlay_menu" :model="usermenuItems" :popup="true" />

        <vueDialog
          v-model:visible="settingsDialogVisible"
          modal
          header="Edit Notification Preferences"
          :style="{ width: '30rem' }"
        >
          <div class="flex align-items-center gap-3 mb-5">
            <label for="city" class="font-semibold w-6rem">City</label>
            <Dropdown
              id="city"
              v-model="selectedCity"
              :options="cities"
              optionLabel="name"
              optionValue="code"
              placeholder="Select a City"
              class="w-full"
              v-on:change="getCityUserPreferences()"
            />
          </div>
          <Divider align="center" type="solid">
            <b>Preferences for selected City</b>
          </Divider>
          <div class="field">
            <label for="email">E-Mail</label>
            <InputText id="email" type="email" v-model="userPreferences.email" class="w-full" />
          </div>
          <div class="field">
            <label for="type">Type of acquisition</label>
            <SelectButton
              id="type"
              v-model="userPreferences.listingType"
              :options="types"
              optionLabel="name"
              optionValue="type"
              aria-labelledby="basic"
              class="w-full"
            />
          </div>
          <div class="field">
            <label for="geno">Genossenschaft</label>
            <Dropdown
              id="geno"
              v-model="userPreferences.housingCooperative"
              :options="genos"
              showClear
              optionLabel="name"
              optionValue="code"
              placeholder="Select a Genossenschaft"
              class="w-full"
            />
          </div>
          <div class="field">
            <label for="priceFrom">Price € per Month</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="priceFrom"
                v-model="userPreferences.minRentPrice"
                placeholder="from"
                inputClass="w-full"
                locale="de-DE"
              />
              <p>-</p>
              <InputNumber
                inputId="priceTo"
                v-model="userPreferences.maxRentPrice"
                placeholder="to"
                inputClass="w-full"
                locale="de-DE"
              />
            </div>
          </div>
          <div class="field">
            <label for="roomsFrom">Rooms</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="roomsFrom"
                v-model="userPreferences.minRoomCount"
                placeholder="from"
                inputClass="w-full"
                :useGrouping="false"
              />
              <p>-</p>
              <InputNumber
                inputId="roomsTo"
                v-model="userPreferences.maxRoomCount"
                placeholder="to"
                inputClass="w-full"
                :useGrouping="false"
              />
            </div>
          </div>
          <div class="field">
            <label for="areaFrom">Area m²</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="areaFrom"
                v-model="userPreferences.minSqm"
                placeholder="from"
                inputClass="w-full"
                :useGrouping="false"
              />
              <p>-</p>
              <InputNumber
                inputId="areaTo"
                v-model="userPreferences.maxSqm"
                placeholder="to"
                inputClass="w-full"
                :useGrouping="false"
              />
            </div>
          </div>
          <div class="field">
            <label for="postalCode">Postal Code</label>
            <InputText inputId="postalCode" v-model="userPreferences.postalCode" class="w-full" />
          </div>
          <div class="field">
            <label for="dateAvaialble">Available From</label>
            <Calendar
              inputId="dateAvaialble"
              v-model="userPreferences.availableFrom"
              dateFormat="dd.mm.yy"
              showIcon
              iconDisplay="input"
              placeholder="dd.mm.yyyy"
              class="w-full"
            />
          </div>
          <div class="field">
            <label for="hwgClass">Hwg Energy Class (worst acceptable)</label>
            <Dropdown
              id="hwgClass"
              v-model="userPreferences.minHwgEnergyClass"
              :options="energyClasses"
              showClear
              optionLabel="name"
              optionValue="value"
              placeholder="Worst acceptable"
              class="w-full"
            />
          </div>
          <div class="field">
            <label for="fgeeClass">Fgee Energy Class (worst acceptable)</label>
            <Dropdown
              id="fgeeClass"
              v-model="userPreferences.minFgeeEnergyClass"
              :options="energyClasses"
              showClear
              optionLabel="name"
              optionValue="value"
              placeholder="Worst acceptable"
              class="w-full"
            />
          </div>
          <div class="field">
            <label for="yearBuild">Year Built</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="yearBuild"
                v-model="userPreferences.minYearBuilt"
                placeholder="min"
                inputClass="w-full"
                :useGrouping="false"
              />
              <p>-</p>
              <InputNumber
                v-model="userPreferences.maxYearBuilt"
                placeholder="max"
                inputClass="w-full"
                :useGrouping="false"
              />
            </div>
          </div>
          <div class="field">
            <label for="cooperativeShare">Cooperative Share</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="cooperativeShare"
                v-model="userPreferences.minCooperativeShare"
                placeholder="min"
                inputClass="w-full"
                locale="de-DE"
              />
              <p>-</p>
              <InputNumber
                v-model="userPreferences.maxCooperativeShare"
                placeholder="max"
                inputClass="w-full"
                locale="de-DE"
              />
            </div>
          </div>
          <div class="field">
            <label for="salePrice">Sale Price</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="salePrice"
                v-model="userPreferences.minSalePrice"
                placeholder="min"
                inputClass="w-full"
                locale="de-DE"
              />
              <p>-</p>
              <InputNumber
                v-model="userPreferences.maxSalePrice"
                placeholder="max"
                inputClass="w-full"
                locale="de-DE"
              />
            </div>
          </div>

          <div class="flex justify-content-end gap-2">
            <vueButton
              type="button"
              label="Cancel"
              severity="secondary"
              @click="settingsDialogVisible = false"
            ></vueButton>
            <vueButton type="button" label="Save" @click="saveUserPreferences()"></vueButton>
          </div>
        </vueDialog>
      </div>
    </template>
  </vueMenubar>
</template>

<style scoped></style>
