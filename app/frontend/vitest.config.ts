/// <reference types="vitest" />
import { getViteConfig } from 'astro/config'

export default getViteConfig({
  //@ts-ignore TS-2353
  test: {
    /* for example, use global to avoid globals imports (describe, test, expect): */
    globals: true,
  },
})
