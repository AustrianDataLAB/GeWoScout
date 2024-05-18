<script lang="ts" setup>

import { onMounted, ref, type Ref } from "vue";
import Button from 'primevue/button';
import Card from 'primevue/card';
import { getListings } from "@/common/api-service";
import type { Listing } from "@/types/ApiResponseListings";

const props = withDefaults(defineProps<{ searchCity?: string }>(), {
  searchCity: 'vienna'
})

/* 
const results: Ref<Listing[]> = ref([
  {
    title: 'Leo am Teich - Wohnen am Wasser',
    postalCode: 1010,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Leo am Teich - Wohnen am Wasser - Provisionsfrei!',
    postalCode: 1011,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Leo am Teich - Wohnen am Wasser',
    postalCode: 1010,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Leo am Teich - Wohnen am Wasser',
    postalCode: 1015,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: '2 Zimmer mit Küche und riesigem Balkon!',
    postalCode: 1013,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Martha im Grün - gefördertes Eigentum beim Badeteich',
    postalCode: 1012,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: '2 Zimmer mit Küche und riesigem Balkon!',
    postalCode: 1012,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  }
]); 
*/

// const realResults : Ref<Listing[]> = ref(await getListings(props.searchCity));

const realResults: Ref<Listing[]> = ref([]);
onMounted(async () => {realResults.value = await getListings(props.searchCity);});

function redirectToAppartment(index: number) {
  window.open(
    realResults.value[index].detailsUrl,
    '_blank'
  );
}

</script>


<template>
  <div class="cards mt-3 grid">
    <div class="col-12 lg:col-4" v-for="(item, index) in realResults" :key="index">
      <Card style="overflow: hidden">
        <template #header>
           <img 
            alt="appartment" 
            :src="item.previewImageUrl" 
            width="450" height="180" 
            onerror="this.onerror=null;this.src='/src/assets/temp.jpg';" />
        </template>
        <template #title>{{ item.title }}</template>
        <template #subtitle><span class="pi pi-map-marker"></span> {{ item.postalCode }} {{ item.city }}</template>
        <template #content>
          <div class="card flex justify-content-around">
            <div class="flex flex-column m-0">
              <p>Rooms</p>
              <p class="text-center m-0">2</p>
            </div>
            <div class="flex flex-column m-0">
              <p>Area</p>
              <p class="text-center m-0">75 m²</p>
            </div>
          </div>
        </template>
        <template #footer>
          <div class="flex gap-3 mt-1">
            <Button label="Details" severity="secondary" outlined class="w-full" />
            <Button label="Request" icon="pi pi-external-link" class="w-full" @click="redirectToAppartment(index)" />
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>


<style scoped>

</style>
