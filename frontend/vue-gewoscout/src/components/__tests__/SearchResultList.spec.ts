import { describe, it, expect } from 'vitest';

import { mount } from '@vue/test-utils';
import SearchResultList from '../SearchResultList.vue';
import { createPinia, setActivePinia } from 'pinia';
import { createTestingPinia } from '@pinia/testing';
import { useListingsStore } from '@/common/store';

describe('SearchResultList', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('renders properly', () => {
    const wrapper = mount(SearchResultList, {
      global: {
        plugins: [
          createTestingPinia({
            initialState: {
              listings: {
                listings: [
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
                  }
                ]
              }
            }
          })
        ]
      }
    });

    const store = useListingsStore();
    store.listings;

    expect(SearchResultList).toBeTruthy();
    expect(wrapper.text()).toContain('Leo am Teich - Wohnen am Wasser');
  });
});
