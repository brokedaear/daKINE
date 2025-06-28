import { test, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import { fakePage } from './testVariables'

test.describe('Page accessibility', { tag: '@unit' }, () => {
  test('Homepage has no accessibility issues', async ({ page }) => {
    await page.goto('/')

    const accessibilityScanResults = await new AxeBuilder({ page }).analyze()

    expect(accessibilityScanResults.violations).toHaveLength(0)
  })

  test('About page has no accessibility issues', async ({ page }) => {
    await page.goto('/about')

    const accessibilityScanResults = await new AxeBuilder({ page }).analyze()

    expect(accessibilityScanResults.violations).toHaveLength(0)
  })

  test('Privacy Policy page has no accessibility issues', async ({ page }) => {
    await page.goto('/about/privacy-policy')

    const accessibilityScanResults = await new AxeBuilder({ page }).analyze()

    expect(accessibilityScanResults.violations).toHaveLength(0)
  })

  test('404 page has no accessibility issues', async ({ page }) => {
    await page.goto(fakePage)

    const accessibilityScanResults = await new AxeBuilder({ page }).analyze()

    expect(accessibilityScanResults.violations).toHaveLength(0)
  })
})
