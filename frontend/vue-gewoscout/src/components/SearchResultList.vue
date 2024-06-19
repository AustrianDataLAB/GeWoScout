<script lang="ts" setup>
import Button from 'primevue/button';
import Card from 'primevue/card';
import ScrollTop from 'primevue/scrolltop';
import { useListingsStore } from '@/common/store';

const listingsStore = useListingsStore();

function redirectToApartment(index: number) {
  window.open(listingsStore.listings[index].detailsUrl, '_blank');
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
  <h1 v-if="listingsStore.listings.length == 0">Teeest</h1>
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
  </div>
</template>

<style scoped>
img {
  background-image: url('/src/assets/temp.jpg');
  background-size: 418px 180px;
  background-repeat: no-repeat;
}
</style>
