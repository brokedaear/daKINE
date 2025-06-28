import { test, expect } from '@playwright/test'
import { fakePage } from './testVariables'

test('404 has error code title', async ({ page }) => {
  await page.goto(fakePage)

  await expect(page.getByRole('heading', { name: '404' })).toBeVisible()
})
