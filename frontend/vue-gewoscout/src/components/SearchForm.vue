<script lang="ts" setup>
import type SearchInputs from '@/types/SearchInputs';
import { onMounted, ref, type Ref } from 'vue';

import Button from 'primevue/button';
import Dropdown from 'primevue/dropdown';
import InputNumber from 'primevue/inputnumber';
import SelectButton from 'primevue/selectbutton';
import { getListings } from '@/common/api-service';
import { useListingsStore } from '@/common/store';

const selectedPriceFrom = ref(null);
const selectedPriceTo = ref();

const selectedRoomsFrom = ref();
const selectedRoomsTo = ref();

const selectedAreaFrom = ref();
const selectedAreaTo = ref();

const searchInputs: Ref<SearchInputs> = ref({
  city: 'vienna',
  geno: ''
});

const cities = ref([
  { name: 'Vienna', code: 'vienna' },
  { name: 'Graz', code: 'graz' }
]);

const genos = ref([
  { name: 'bwsg', code: 'bwsg' },
  { name: 'ÖVW', code: 'oevw' }
]);

const selectedTypes = ref(['All']);
const types = ref(['All', 'Rent', 'Rent + Option to buy']);

const listingsStore = useListingsStore();

onMounted(async () => {
  listingsStore.listings = await getListings(searchInputs.value);
});

async function search() {
  listingsStore.listings = await getListings(searchInputs.value);
}
</script>

<template>
  <div class="form mt-3">
    <div class="formgrid grid">
      <div class="field col-4">
        <label for="type">Type of acquisition</label>
        <SelectButton
          id="type"
          v-model="selectedTypes"
          :options="types"
          multiple
          aria-labelledby="multiple"
          class="w-full"
        />
      </div>
      <div class="field col-4">
        <label for="city">City</label>
        <Dropdown
          id="city"
          v-model="searchInputs.city"
          :options="cities"
          showClear
          optionLabel="name"
          optionValue="code"
          placeholder="Select a City"
          class="w-full"
        />
      </div>
      <div class="field col-4">
        <label for="geno">Genossenschaft</label>
        <Dropdown
          id="geno"
          v-model="searchInputs.geno"
          :options="genos"
          showClear
          optionLabel="name"
          optionValue="code"
          placeholder="Select a Genossenschaft"
          class="w-full"
        />
      </div>
      <div class="field col-4">
        <label for="priceFrom">Price €</label>
        <div class="flex flex-row gap-2">
          <InputNumber
            inputId="priceFrom"
            v-model="selectedPriceFrom"
            placeholder="from"
            inputClass="w-10rem"
          />
          <p>-</p>
          <InputNumber
            inputId="priceTo"
            v-model="selectedPriceTo"
            placeholder="to"
            inputClass="w-10rem"
          />
        </div>
      </div>
      <div class="field col-2">
        <label for="roomsFrom">Rooms</label>
        <div class="flex flex-row gap-2">
          <InputNumber
            inputId="roomsFrom"
            v-model="selectedRoomsFrom"
            placeholder="from"
            inputClass="w-full"
          />
          <p>-</p>
          <InputNumber
            inputId="roomsTo"
            v-model="selectedRoomsTo"
            placeholder="to"
            inputClass="w-full"
          />
        </div>
      </div>
      <div class="field col-2">
        <label for="areaFrom">Area m²</label>
        <div class="flex flex-row gap-2">
          <InputNumber
            inputId="areaFrom"
            v-model="selectedAreaFrom"
            placeholder="from"
            inputClass="w-full"
          />
          <p>-</p>
          <InputNumber
            inputId="areaTo"
            v-model="selectedAreaTo"
            placeholder="to"
            inputClass="w-full"
          />
        </div>
      </div>
      <div class="field col-4 text-right align-self-end">
        <Button class="mr-3" label="Reset Filters" icon="pi pi-undo" severity="secondary" />
        <Button label="Search" icon="pi pi-search" @click="search()" />
      </div>
    </div>
  </div>
</template>

<style scoped></style>
