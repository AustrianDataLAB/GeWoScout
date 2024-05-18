import { describe, it, expect } from 'vitest'

import { mount } from '@vue/test-utils'
import SearchForm from '../SearchForm.vue'

describe('SearchForm', () => {
  it('renders properly', () => {
    expect(SearchForm).toBeTruthy();

    const wrapper = mount(SearchForm);
    expect(wrapper.text()).toContain('Type of acquisition');
    expect(wrapper.text()).toContain('Price €');
  })
})
