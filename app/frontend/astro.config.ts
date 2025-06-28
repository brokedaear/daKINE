// @ts-check
import { defineConfig, envField } from 'astro/config'

import sitemap from '@astrojs/sitemap'
import AutoImport from 'astro-auto-import'
import compressor from 'astro-compressor'
import mdx from '@astrojs/mdx'

import react from '@astrojs/react'
import favicons from 'astro-favicons'

import netlify from '@astrojs/netlify'
import tailwindcss from '@tailwindcss/vite'
// import { sessionCookieName } from './src/lib/session/constants'

// https://astro.build/config
export default defineConfig({
  site: 'https://brokedaear.com',
  output: 'server',
  compressHTML: import.meta.env.PROD,
  // This is enabled by default, however, its better to be
  // clear about it here.
  experimental: {
    // https://www.piyushmehta.com/blog/astro-v5-9-content-security-policy
    csp: {
      directives: [
        "default-src: 'self'",
        'upgrade-insecure-requests',
        "connect-src: 'self' https://api.stripe.com",
        // For some reason a 'font-src' directive is not necessary.
        //
        // TODO: add one for the backend?
      ],
      styleDirective: {
        resources: [
          // Ensure inline styles work.
          'self',
          "'unsafe-line'",
          // Ensure google fonts work.
          'self',
          'https://fonts.googleapis.com',
        ],
      },
      scriptDirective: {
        resources: ['self', 'http://localhost:4321'],
        strictDynamic: true,
      },
    },
  },
  security: {
    checkOrigin: true,
  },
  session: {
    // driver: 'redis',
    // options: {
    //   url: import.meta.env.REDIS_URL,
    //   base: 'session',
    //   ttl: 86400,
    // },
    driver: 'lru-cache',
    cookie: {
      name: 'bde-web-session',
      // Use 'lax' over 'strict', users aren't logged in if visiting
      // the site from an external link.
      sameSite: 'lax',
      secure: import.meta.env.PROD,
      path: '/',
    },
    ttl: 86400, // an entire day.
  },
  env: {
    schema: {
      STRIPE_PUBLIC_API_KEY: envField.string({
        context: 'client',
        access: 'public',
      }),
      BACKEND_SERVICE_URL: envField.string({
        context: 'client',
        access: 'public',
      }),
    },
  },
  adapter: netlify({
    edgeMiddleware: true,
    experimentalStaticHeaders: true,
  }),
  vite: {
    //@ts-ignore TS-2322
    plugins: [tailwindcss()],
  },
  integrations: [
    react(),
    sitemap(),
    favicons(),

    AutoImport({
      imports: [],
    }),

    // mdx MUST come after AutoImport.
    mdx(),

    // Compressor MUST come last.
    compressor(),
  ],
})
