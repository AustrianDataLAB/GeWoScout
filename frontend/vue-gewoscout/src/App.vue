<script setup lang="ts">
import HeaderNavbar from './components/HeaderNavbar.vue';
import SearchForm from './components/SearchForm.vue';
import SearchResultList from './components/SearchResultList.vue';
</script>

<template>
  <header>
    <HeaderNavbar />
  </header>

  <main>
    <SearchForm />
      <SearchResultList searchCity="vienna" />
  </main>
</template>

<style scoped>

</style>
