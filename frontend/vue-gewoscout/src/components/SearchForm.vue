<script lang="ts" setup>
import type SearchInputs from '@/types/SearchInputs';
import { onMounted, ref, type Ref } from 'vue';

import Button from 'primevue/button';
import Dropdown from 'primevue/dropdown';
import InputNumber from 'primevue/inputnumber';
import InputText from 'primevue/inputtext';
import SelectButton from 'primevue/selectbutton';
import Calendar from 'primevue/calendar';
import Accordion from 'primevue/accordion';
import AccordionTab from 'primevue/accordiontab';

import { getListings } from '@/common/api-service';
import { useListingsStore } from '@/common/store';
import { EnergyClass, Type } from '@/types/Enums';

const searchInputs: Ref<SearchInputs> = ref({
  listingType: Type.both,
  city: 'vienna',
  housingCooperative: '',
  postalCode: '',
  minRoomCount: null,
  maxRoomCount: null,
  minSqm: null,
  maxSqm: null,
  availableFrom: null,
  minYearBuilt: null,
  maxYearBuilt: null,
  minHwgEnergyClass: null,
  minFgeeEnergyClass: null,
  minRentPrice: null,
  maxRentPrice: null,
  minCooperativeShare: null,
  maxCooperativeShare: null,
  minSalePrice: null,
  maxSalePrice: null
});

const cities = ref([
  { name: 'Vienna', code: 'vienna' },
  { name: 'Graz', code: 'graz' }
]);

const genos = ref([
  { name: 'bwsg', code: 'bwsg' },
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

const listingsStore = useListingsStore();

onMounted(async () => {
  listingsStore.listings = await getListings(searchInputs.value);
});

async function search() {
  listingsStore.listings = await getListings(searchInputs.value);
}

function reset() {
  searchInputs.value = {
    listingType: Type.both,
    city: 'vienna',
    housingCooperative: '',
    postalCode: '',
    minRoomCount: null,
    maxRoomCount: null,
    minSqm: null,
    maxSqm: null,
    availableFrom: null,
    minYearBuilt: null,
    maxYearBuilt: null,
    minHwgEnergyClass: null,
    minFgeeEnergyClass: null,
    minRentPrice: null,
    maxRentPrice: null,
    minCooperativeShare: null,
    maxCooperativeShare: null,
    minSalePrice: null,
    maxSalePrice: null
  };
}
</script>

<template>
  <div class="form mt-3">
    <div class="formgrid grid">
      <div class="field col-3">
        <label for="type">Type of acquisition</label>
        <SelectButton
          id="type"
          v-model="searchInputs.listingType"
          :options="types"
          optionLabel="name"
          optionValue="type"
          aria-labelledby="basic"
          class="w-full"
        />
      </div>
      <div class="field col-3">
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
      <div class="field col-3">
        <label for="geno">Genossenschaft</label>
        <Dropdown
          id="geno"
          v-model="searchInputs.housingCooperative"
          :options="genos"
          showClear
          optionLabel="name"
          optionValue="code"
          placeholder="Select a Genossenschaft"
          class="w-full"
        />
      </div>
      <div class="field col-3"></div>
      <div class="field col-3">
        <label for="priceFrom">Price € per Month</label>
        <div class="flex flex-row gap-2">
          <InputNumber
            inputId="priceFrom"
            v-model="searchInputs.minRentPrice"
            placeholder="from"
            inputClass="w-full"
            locale="de-DE"
          />
          <p>-</p>
          <InputNumber
            inputId="priceTo"
            v-model="searchInputs.maxRentPrice"
            placeholder="to"
            inputClass="w-full"
            locale="de-DE"
          />
        </div>
      </div>
      <div class="field col-3">
        <label for="roomsFrom">Rooms</label>
        <div class="flex flex-row gap-2">
          <InputNumber
            inputId="roomsFrom"
            v-model="searchInputs.minRoomCount"
            placeholder="from"
            inputClass="w-full"
            :useGrouping="false"
          />
          <p>-</p>
          <InputNumber
            inputId="roomsTo"
            v-model="searchInputs.maxRoomCount"
            placeholder="to"
            inputClass="w-full"
            :useGrouping="false"
          />
        </div>
      </div>
      <div class="field col-3">
        <label for="areaFrom">Area m²</label>
        <div class="flex flex-row gap-2">
          <InputNumber
            inputId="areaFrom"
            v-model="searchInputs.minSqm"
            placeholder="from"
            inputClass="w-full"
            :useGrouping="false"
          />
          <p>-</p>
          <InputNumber
            inputId="areaTo"
            v-model="searchInputs.maxSqm"
            placeholder="to"
            inputClass="w-full"
            :useGrouping="false"
          />
        </div>
      </div>
      <div class="field col-3 text-right align-self-end">
        <Button
          class="mr-3"
          label="Reset Filters"
          icon="pi pi-undo"
          severity="secondary"
          @click="reset()"
        />
        <Button label="Search" icon="pi pi-search" @click="search()" />
      </div>
    </div>
    <Accordion>
      <AccordionTab header="Detailed Search">
        <div class="formgrid grid detailedSearch">
          <div class="field col-3">
            <label for="postalCode">Postal Code</label>
            <InputText inputId="postalCode" v-model="searchInputs.postalCode" class="w-full" />
          </div>
          <div class="field col-3">
            <label for="dateAvaialble">Available From</label>
            <Calendar
              inputId="dateAvaialble"
              v-model="searchInputs.availableFrom"
              dateFormat="dd.mm.yy"
              showIcon
              iconDisplay="input"
              placeholder="dd.mm.yyyy"
              class="w-full"
            />
          </div>
          <div class="field col-3">
            <label for="hwgClass">Hwg Energy Class</label>
            <Dropdown
              id="hwgClass"
              v-model="searchInputs.minHwgEnergyClass"
              :options="energyClasses"
              showClear
              optionLabel="name"
              optionValue="value"
              placeholder="Worst acceptable"
              class="w-full"
            />
          </div>
          <div class="field col-3">
            <label for="fgeeClass">Fgee Energy Class</label>
            <Dropdown
              id="fgeeClass"
              v-model="searchInputs.minFgeeEnergyClass"
              :options="energyClasses"
              showClear
              optionLabel="name"
              optionValue="value"
              placeholder="Worst acceptable"
              class="w-full"
            />
          </div>
          <div class="field col-3">
            <label for="yearBuild">Year Built</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="yearBuild"
                v-model="searchInputs.minYearBuilt"
                placeholder="min"
                inputClass="w-full"
                :useGrouping="false"
              />
              <p>-</p>
              <InputNumber
                v-model="searchInputs.maxYearBuilt"
                placeholder="max"
                inputClass="w-full"
                :useGrouping="false"
              />
            </div>
          </div>
          <div class="field col-3">
            <label for="cooperativeShare">Cooperative Share</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="cooperativeShare"
                v-model="searchInputs.minCooperativeShare"
                placeholder="min"
                inputClass="w-full"
                locale="de-DE"
              />
              <p>-</p>
              <InputNumber
                v-model="searchInputs.maxCooperativeShare"
                placeholder="max"
                inputClass="w-full"
                locale="de-DE"
              />
            </div>
          </div>
          <div class="field col-3">
            <label for="salePrice">Sale Price</label>
            <div class="flex flex-row gap-2">
              <InputNumber
                inputId="salePrice"
                v-model="searchInputs.minSalePrice"
                placeholder="min"
                inputClass="w-full"
                locale="de-DE"
              />
              <p>-</p>
              <InputNumber
                v-model="searchInputs.maxSalePrice"
                placeholder="max"
                inputClass="w-full"
                locale="de-DE"
              />
            </div>
          </div>
        </div>
      </AccordionTab>
    </Accordion>
  </div>
</template>

<style scoped>
.detailedSearch {
  color: black !important;
}
</style>
