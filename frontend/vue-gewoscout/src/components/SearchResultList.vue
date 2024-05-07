<script lang="ts" setup>
import { ref } from "vue";
import Button from 'primevue/button';
import Card from 'primevue/card';

const results = ref([
  {
    titel: 'Leo am Teich - Wohnen am Wasser',
    plz: '1010',
    city: 'Vienna'
  },
  {
    titel: 'Leo am Teich - Wohnen am Wasser - Provisionsfrei!',
    plz: '1011',
    city: 'Vienna'
  },
  {
    titel: 'Leo am Teich - Wohnen am Wasser',
    plz: '1010',
    city: 'Vienna'
  },
  {
    titel: 'Leo am Teich - Wohnen am Wasser',
    plz: '1015',
    city: 'Vienna'
  },
  {
    titel: '2 Zimmer mit Küche und riesigem Balkon!',
    plz: '1013',
    city: 'Vienna'
  },
  {
    titel: 'Martha im Grün - gefördertes Eigentum beim Badeteich',
    plz: '1012',
    city: 'Vienna'
  },
  {
    titel: '2 Zimmer mit Küche und riesigem Balkon!',
    plz: '1012',
    city: 'Vienna'
  }
])

</script>

<template>
  <div class="cards mt-3 grid">
    <div class="col-12 lg:col-4" v-for="(item, index) in results" :key="index">
      <Card style="overflow: hidden">
        <template #header>
          <img alt="user header" src="../assets/tempFlat.jpg" width="450" height="180" />
        </template>
        <template #title>{{ item.titel }}</template>
        <template #subtitle><span class="pi pi-map-marker"></span> {{ item.plz }} {{ item.city }}</template>
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
            <Button label="Request" icon="pi pi-external-link" class="w-full" />
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  position: relative;
  top: -10px;
}

h3 {
  font-size: 1.2rem;
}

.greetings h1,
.greetings h3 {
  text-align: center;
}

@media (min-width: 1024px) {
  .greetings h1,
  .greetings h3 {
    text-align: left;
  }
}
</style>
