<script lang="ts" setup>
import { ref } from "vue";
import Button from 'primevue/button';
import Card from 'primevue/card';

const results = ref([
  {
    titel: 'Advanced Card',
  },
  {
    titel: 'Advanced Card',
  },
  {
    titel: 'Advanced Card',
  },
  {
    titel: 'Advanced Card',
  },
  {
    titel: 'Advanced Card',
  },
  {
    titel: 'Advanced Card',
  },
  {
    titel: 'Advanced Card',
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
        <template #subtitle>Card subtitle</template>
        <template #content>
          <p class="m-0">
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Lorem ipsum dolor sit amet, consectetur adipisicing elit. Lorem ipsum dolor sit amet, consectetur adipisicing elit. 
          </p>
        </template>
        <template #footer>
          <div class="flex gap-3 mt-1">
            <Button label="Cancel" severity="secondary" outlined class="w-full" />
            <Button label="Save" class="w-full" />
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
