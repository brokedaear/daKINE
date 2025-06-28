import { test, expect } from '@playwright/test'

test.describe('Homepage', { tag: '@unit' }, () => {
  test('Title visible', async ({ page }) => {
    await page.goto('/')

    await expect(
      page.getByRole('heading', { name: 'BROKE DA EAR! AUDIO PLUGINS' })
    ).toBeVisible()
  })

  test('Contact info visible', async ({ page }) => {
    await page.goto('/')

    await expect(page.getByText('howzit@brokedaear.com')).toBeVisible()
  })

  test('Logo image visible', async ({ page }) => {
    await page.goto('/')

    await expect(
      page.getByRole('img', { name: 'Broke da EAR! Logo' })
    ).toBeVisible()
  })
})
