<script lang="ts" setup>

import type SearchInputs from '@/types/SearchInputs';
import { ref } from 'vue';

import SelectButton from 'primevue/selectbutton';
import Dropdown from 'primevue/dropdown';
import InputNumber from 'primevue/inputnumber';
import Button from 'primevue/button';

const selectedPriceFrom = ref(null);
const selectedPriceTo = ref();

const selectedRoomsFrom = ref();
const selectedRoomsTo = ref();

const selectedAreaFrom = ref();
const selectedAreaTo = ref();

const searchInputs: SearchInputs = {
  city: '',
}

const selectedCity = ref('vienna');
const cities = ref([
    { name: 'Vienna', code: 'vienna' },
    { name: 'Graz', code: 'graz' }
]);

const selectedGeno = ref();
const genos = ref([
    { name: 'bwsg', code: 'bwsg' },
    { name: 'ÖVW', code: 'oevw' }
]);

const selectedTypes = ref(['All']);
const types = ref(['All', 'Rent', 'Rent + Option to buy']);
</script>

<template>
  <div class="form mt-3">
    <div class="formgrid grid">
      <div class="field col-4">
        <label for="type">Type of acquisition</label>
        <SelectButton id="type" v-model="selectedTypes" :options="types" multiple aria-labelledby="multiple" class="w-full" />
      </div>
      <div class="field col-4">
        <label for="city">City</label>
        <Dropdown id="city" v-model="selectedCity" :options="cities" showClear optionLabel="name" optionValue="code" placeholder="Select a City" class="w-full" />
      </div>
      <div class="field col-4">
        <label for="geno">Genossenschaft</label>
        <vueDroDropdownpdown id="geno" v-model="selectedGeno" :options="genos" showClear optionLabel="name" optionValue="code" placeholder="Select a Genossenschaft" class="w-full" />
      </div>
      <div class="field col-4">
        <label for="priceFrom">Price €</label>
        <div class="flex flex-row gap-2">
          <InputNumber inputId="priceFrom" v-model="selectedPriceFrom" placeholder="from" inputClass="w-10rem" />
          <p>-</p>
          <InputNumber inputId="priceTo" v-model="selectedPriceTo" placeholder="to" inputClass="w-10rem" />
        </div>
      </div>
      <div class="field col-2">
        <label for="roomsFrom">Rooms</label>
        <div class="flex flex-row gap-2">
          <InputNumber inputId="roomsFrom" v-model="selectedRoomsFrom" placeholder="from" inputClass="w-full" />
          <p>-</p>
          <InputNumber inputId="roomsTo" v-model="selectedRoomsTo" placeholder="to" inputClass="w-full" />
        </div>
      </div>
      <div class="field col-2">
        <label for="areaFrom">Area m²</label>
        <div class="flex flex-row gap-2">
          <InputNumber inputId="areaFrom" v-model="selectedAreaFrom" placeholder="from" inputClass="w-full" />
          <p>-</p>
          <InputNumber inputId="areaTo" v-model="selectedAreaTo" placeholder="to" inputClass="w-full" />
        </div>
      </div>
      <div class="field col-4 text-right align-self-end">
        <Button label="Reset Filters" icon="pi pi-undo" severity="secondary" />
      </div>
    </div>
  </div>
</template>


<style scoped>

</style>
