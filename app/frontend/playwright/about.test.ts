import { test, expect } from '@playwright/test'

test.describe('About page', { tag: '@unit' }, () => {
  test('test', async ({ page }) => {
    await page.goto('/about')
    await page.getByRole('heading', { name: 'ABOUT US' }).isVisible()
  })
})

test.describe('Privacy Policy page', { tag: '@unit' }, () => {
  test('Privacy Policy heading visible', async ({ page }) => {
    await page.goto('/about/privacy-policy')
    await expect(
      page.getByRole('heading', { name: 'Broke da EAR! Privacy Policy' })
    ).toBeVisible()
  })
  test('Effective Date heading visible', async ({ page }) => {
    await page.goto('/about/privacy-policy')
    await expect(page.getByText('Effective Date: Thu Jan 23')).toBeVisible()
  })
})
