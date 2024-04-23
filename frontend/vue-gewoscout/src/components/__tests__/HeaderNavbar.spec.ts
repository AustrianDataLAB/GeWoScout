import { describe, it, expect } from 'vitest'

import { mount } from '@vue/test-utils'
import HeaderNavbar from '../HeaderNavbar.vue'

describe('HeaderNavbar', () => {
  it('renders properly', () => {
    const wrapper = mount(HeaderNavbar)
    expect(wrapper.text()).toContain('Search')
  })
})
