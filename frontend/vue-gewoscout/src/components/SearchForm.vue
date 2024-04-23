<script setup lang="ts">
import { ref } from 'vue';

const value = ref(null);

</script>

<template>
  <div class="form">
    <div class="flex flex-column gap-2">
        <label for="username">Username</label>
        <vueInputText id="username" v-model="value" />
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
