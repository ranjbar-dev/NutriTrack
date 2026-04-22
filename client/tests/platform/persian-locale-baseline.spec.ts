import { describe, expect, it } from 'vitest'

describe('persian locale baseline scaffolding', () => {
  it('can format a sample number in Persian locale', () => {
    const formatted = new Intl.NumberFormat('fa-IR').format(1405)
    expect(formatted).not.toBe('1405')
  })
})
