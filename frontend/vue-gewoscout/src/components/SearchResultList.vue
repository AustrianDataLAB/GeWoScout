<script lang="ts" setup>
import Button from 'primevue/button';
import Card from 'primevue/card';
import ScrollTop from 'primevue/scrolltop';
import ProgressSpinner from 'primevue/progressspinner';
import Divider from 'primevue/divider';
import Dialog from 'primevue/dialog';
import { useListingsStore } from '@/common/store';
import { ref, type Ref } from 'vue';
import type { Listing } from '@/types/ApiResponseListings';

const listingsStore = useListingsStore();

const detailsDialogVisible = ref(false);

const selectedListing: Ref<Listing | null> = ref(null);

function redirectToApartment(index: number) {
  if (listingsStore.listings !== null) {
    window.open(listingsStore.listings[index].detailsUrl, '_blank');
  }
}

function redirectToApartmentFromDetails() {
  if (selectedListing.value !== null) {
    window.open(selectedListing.value.detailsUrl, '_blank');
  }
}

function openDetails(index: number) {
  detailsDialogVisible.value = true;

  if (listingsStore.listings !== null) {
    selectedListing.value = listingsStore.listings[index];
  }
}

window.onscroll = () => {
  const bottomOfWindow =
    document.documentElement.scrollTop + window.innerHeight ===
    document.documentElement.offsetHeight;

  if (bottomOfWindow) {
    // TODO load new listings and add to end of list
    console.log('End of list');
  }
};
</script>

<template>
  <div class="flex justify-content-center">
    <ProgressSpinner v-show="listingsStore.listings == null" />
  </div>

  <div v-if="listingsStore.listings !== null">
    <h3 v-if="listingsStore.listings.length == 0" class="flex justify-content-center">
      No apartments found, Sorry :(
    </h3>
    <div v-else class="cards mt-3 grid">
      <div class="col-12 lg:col-4" v-for="(item, index) in listingsStore.listings" :key="index">
        <Card style="overflow: hidden">
          <template #header>
            <img :src="item.previewImageUrl" width="420" height="180" />
          </template>
          <template #title>{{ item.title }}</template>
          <template #subtitle
            ><span class="pi pi-map-marker"></span> {{ item.postalCode }} {{ item.city }}</template
          >
          <template #content>
            <div class="card flex justify-content-around">
              <div class="flex flex-column m-0">
                <p>Rooms</p>
                <p class="text-center m-0">{{ item.roomCount }}</p>
              </div>
              <div class="flex flex-column m-0">
                <p>Area</p>
                <p class="text-center m-0">{{ item.squareMeters }} m²</p>
              </div>
            </div>
          </template>
          <template #footer>
            <div class="flex gap-3 mt-1">
              <Button
                label="Details"
                severity="secondary"
                outlined
                class="w-full"
                @click="openDetails(index)"
              />
              <Button
                label="Request"
                icon="pi pi-external-link"
                class="w-full"
                @click="redirectToApartment(index)"
              />
            </div>
          </template>
        </Card>
      </div>
      <ScrollTop />

      <Dialog
        v-model:visible="detailsDialogVisible"
        modal
        header="Details"
        :style="{ width: '30rem' }"
      >
        <div v-if="selectedListing !== null">
          <Card style="overflow: hidden">
            <template #header>
              <img :src="selectedListing.previewImageUrl" width="420" height="180" />
            </template>
            <template #title>{{ selectedListing.title }}</template>
            <template #subtitle
              ><span class="pi pi-map-marker"></span> {{ selectedListing.postalCode }}
              {{ selectedListing.city }}</template
            >
            <template #content>
              <div class="card flex justify-content-around">
                <div class="flex flex-column m-0">
                  <p>Rooms</p>
                  <p class="text-center m-0">{{ selectedListing.roomCount }}</p>
                </div>
                <div class="flex flex-column m-0">
                  <p>Area</p>
                  <p class="text-center m-0">{{ selectedListing.squareMeters }} m²</p>
                </div>
              </div>
              <Divider align="center" type="solid" />
              <div class="card flex justify-content-around">
                <div class="flex flex-column m-0">
                  <p><strong>Genossenschaft</strong></p>
                  <p><strong>Availability Date</strong></p>
                  <p><strong>Year Built</strong></p>
                  <p><strong>HWG Energy Class</strong></p>
                  <p><strong>FGEE Energy Class</strong></p>
                  <p><strong>Rent Price Per Month</strong></p>
                  <p><strong>Cooperative Share</strong></p>
                  <p><strong>Sale Price</strong></p>
                  <p><strong>Additional Fees</strong></p>
                </div>
                <div class="flex flex-column m-0 text-right">
                  <p>{{ selectedListing.housingCooperative }}</p>
                  <p>{{ selectedListing.availabilityDate }}</p>
                  <p>{{ selectedListing.yearBuilt }}</p>
                  <p>{{ selectedListing.hwgEnergyClass }}</p>
                  <p>{{ selectedListing.fgeeEnergyClass }}</p>
                  <p>€ {{ selectedListing.rentPricePerMonth }}</p>
                  <p>€ {{ selectedListing.cooperativeShare }}</p>
                  <p>€ {{ selectedListing.salePrice }}</p>
                  <p>€ {{ selectedListing.additionalFees }}</p>
                </div>
              </div>
            </template>
            <template #footer>
              <div class="flex gap-3 mt-1">
                <Button label="Request" class="w-full" @click="redirectToApartmentFromDetails()" />
              </div>
            </template>
          </Card>
        </div>
      </Dialog>
    </div>
  </div>
</template>

<style scoped>
img {
  background-image: url('/src/assets/temp.jpg');
  background-size: 418px 180px;
  background-repeat: no-repeat;
}
</style>
