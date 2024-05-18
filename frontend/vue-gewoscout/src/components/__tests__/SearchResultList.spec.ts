import { describe, it, expect } from 'vitest'

import { mount } from '@vue/test-utils'
import SearchResultList from '../SearchResultList.vue'

describe('SearchResultList', () => {
  it('renders properly', () => {
    expect(SearchResultList).toBeTruthy();

    const wrapper = mount(SearchResultList);
    expect(wrapper.text()).toContain('Modern 3-Bedroom Apartment in Central Vienna');
    expect(wrapper.text()).toContain('Cozy Studio in the Heart of the City');
  })
})
