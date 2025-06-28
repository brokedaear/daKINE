// pages/login/github/index.ts
import { sessionCookieName } from '@lib/session/constants'
import { generateState } from 'arctic'

import type { APIContext } from 'astro'

export async function GET(context: APIContext): Promise<Response> {
  const state = generateState()

  const cooky = context.cookies.get(sessionCookieName)

  context.session?.regenerate()

  // context.cookies.set(sessionCookieName, state, {
  //   path: '/',
  //   secure: import.meta.env.PROD,
  //   httpOnly: true,
  //   maxAge: 60 * 10, // 10 minutes
  //   sameSite: 'lax',
  // })

  return context.redirect('/account')
}
