import type { Listing } from '@/types/ApiResponseListings';
import { defineStore } from 'pinia';
import { ref, type Ref } from 'vue';

export const useListingsStore = defineStore('listings', () => {
  const listings: Ref<Listing[] | null> = ref(null);
  const continuationToken: Ref<string | null> = ref(null);

  return { listings, continuationToken };
});

export const useUserStore = defineStore('user', () => {
  const loggedIn: Ref<boolean> = ref(false);

  return { loggedIn };
});
