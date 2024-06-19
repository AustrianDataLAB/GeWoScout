import { describe, it, expect } from 'vitest';

import { mount } from '@vue/test-utils';
import SearchForm from '../SearchForm.vue';
import PrimeVue from 'primevue/config';
import { createPinia, setActivePinia } from 'pinia';

// Import PrimeVue and its styles
import './assets/main.css';
import 'primevue/resources/themes/lara-light-amber/theme.css';
import 'primeicons/primeicons.css';
import 'primeflex/primeflex.css';

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
              ripple: true,
              locale: {
                firstDayOfWeek: 0, // Setting Sunday as the first day of the week
                dayNames: [
                  'Sunday',
                  'Monday',
                  'Tuesday',
                  'Wednesday',
                  'Thursday',
                  'Friday',
                  'Saturday'
                ],
                dayNamesShort: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
                dayNamesMin: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'],
                monthNames: [
                  'January',
                  'February',
                  'March',
                  'April',
                  'May',
                  'June',
                  'July',
                  'August',
                  'September',
                  'October',
                  'November',
                  'December'
                ],
                monthNamesShort: [
                  'Jan',
                  'Feb',
                  'Mar',
                  'Apr',
                  'May',
                  'Jun',
                  'Jul',
                  'Aug',
                  'Sep',
                  'Oct',
                  'Nov',
                  'Dec'
                ],
                today: 'Today',
                clear: 'Clear'
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
