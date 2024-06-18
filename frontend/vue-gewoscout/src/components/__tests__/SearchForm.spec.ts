import { describe, it, expect } from 'vitest';

import { mount } from '@vue/test-utils';
import SearchForm from '../SearchForm.vue';
import PrimeVue from 'primevue/config';
import { createPinia, setActivePinia } from 'pinia';

describe('SearchForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
  });

  it('renders properly', () => {
    expect(SearchForm).toBeTruthy();

    const wrapper = mount(SearchForm, {
      global: {
        plugins: [PrimeVue],
        mocks: {
          $primevue: {
            config: {
              ripple: true
            }
          }
        }
      }
    });

    expect(wrapper.text()).toContain('Type of acquisition');
    expect(wrapper.text()).toContain('Price €');
  });
});
