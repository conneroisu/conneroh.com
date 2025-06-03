import { describe, expect, it } from 'vitest'

// Example utility function to test
function formatDate(date: Date): string {
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    timeZone: 'UTC'
  }).format(date)
}

describe('formatDate', () => {
  it('formats date correctly', () => {
    const date = new Date('2025-01-01T00:00:00.000Z')
    const formatted = formatDate(date)
    expect(formatted).toBe('January 1, 2025')
  })
  
  it('handles different dates', () => {
    const date = new Date('2024-12-25T00:00:00.000Z')
    const formatted = formatDate(date)
    expect(formatted).toBe('December 25, 2024')
  })
  
  it('handles date validation', () => {
    const validDate = new Date('2025-01-01T00:00:00.000Z')
    expect(validDate.getTime()).toBeGreaterThan(0)
    
    const invalidDate = new Date('invalid')
    expect(isNaN(invalidDate.getTime())).toBe(true)
  })
})