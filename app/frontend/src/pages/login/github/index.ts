import { github } from '@lib/oauth'
import { generateState } from 'arctic'

import type { APIContext } from 'astro'

export async function GET(context: APIContext): Promise<Response> {
  const state = generateState()
  const url = github.createAuthorizationURL(state, [])

  context.cookies.set('github_oauth_state', state, {
    path: '/',
    secure: import.meta.env.PROD,
    httpOnly: true,
    maxAge: 60 * 10, // 10 minutes
    sameSite: 'lax',
  })

  return context.redirect(url.toString())
}
