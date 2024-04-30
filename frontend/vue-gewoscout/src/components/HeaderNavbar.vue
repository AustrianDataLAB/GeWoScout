<script setup lang="ts">
import { ref } from "vue";

const items = ref([
  {
    label: 'Search',
    icon: 'pi pi-search',
  }
]);
</script>

<template>
  <vueMenubar :model="items">
    <template #start>
      <svg width="40" height="40" viewBox="0 0 35 40" fill="yellow" xmlns="http://www.w3.org/2000/svg" class="h-2rem">
        <text x="0" y="15" fill="var(--primary-color)" stroke="var(--primary-color)" font-size="15">GeWoScout</text>
        <text x="0" y="40" fill="var(--text-color)" stroke="var(--text-color)" font-size="15">Scout</text>
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
            <vueButton label="Login"></vueButton>
            <!-- <vueAvatar image="/images/avatar/amyelsner.png" shape="circle" /> -->
        </div>
    </template>
  </vueMenubar>
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
