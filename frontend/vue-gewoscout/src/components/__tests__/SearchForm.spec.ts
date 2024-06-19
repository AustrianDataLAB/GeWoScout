import { describe, it, expect } from 'vitest';

import { mount } from '@vue/test-utils';
import SearchForm from '../SearchForm.vue';
import PrimeVue from 'primevue/config';
import { createPinia, setActivePinia } from 'pinia';

// Import PrimeVue and its styles
import 'primevue/resources/themes/lara-light-amber/theme.css';
import 'primeicons/primeicons.css';
import 'primeflex/primeflex.css';

describe('SearchForm', () => {
  beforeEach(() => {
    setActivePinia(createPinia());

    window.matchMedia =
      window.matchMedia ||
      function () {
        return {
          addEventListener: function () {}
        };
      };
  });

  it('renders properly', () => {
    expect(SearchForm).toBeTruthy();

    const wrapper = mount(SearchForm, {
      global: {
        plugins: [PrimeVue],
        mocks: {
          $primevue: {
            config: {
              ripple: true,
              locale: {
                firstDayOfWeek: 0 // Setting Sunday as the first day of the week
              }
            }
          }
        }
      }
    });

    expect(wrapper.text()).toContain('Type of acquisition');
    expect(wrapper.text()).toContain('Price €');
  });
});
