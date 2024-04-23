import { describe, it, expect } from 'vitest'

import { mount } from '@vue/test-utils'
import SearchResultList from '../SearchResultList.vue'

describe('SearchResultList', () => {
  it('renders properly', () => {
    expect(SearchResultList).toBeTruthy();

    const wrapper = mount(SearchResultList)
    expect(wrapper.text()).toContain('Lorem ipsum dolor sit amet')
  })
})
