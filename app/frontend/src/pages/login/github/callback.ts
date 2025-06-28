import {
  generateSessionToken,
  createSession,
  setSessionTokenCookie,
} from '@lib/session'
import { github } from '@lib/oauth'

import type { APIContext } from 'astro'
import type { OAuth2Tokens } from 'arctic'

export async function GET(context: APIContext): Promise<Response> {
  const code = context.url.searchParams.get('code')
  const state = context.url.searchParams.get('state')
  const storedState = context.cookies.get('github_oauth_state')?.value ?? null
  if (code === null || state === null || storedState === null) {
    return new Response(null, {
      status: 400,
    })
  }
  if (state !== storedState) {
    return new Response(null, {
      status: 400,
    })
  }

  let tokens: OAuth2Tokens
  try {
    tokens = await github.validateAuthorizationCode(code)
  } catch (e) {
    // Invalid code or client credentials
    return new Response(null, {
      status: 400,
    })
  }
  const githubUserResponse = await fetch('https://api.github.com/user', {
    headers: {
      Authorization: `Bearer ${tokens.accessToken()}`,
    },
  })
  const githubUser = await githubUserResponse.json()
  const githubUserId = githubUser.id
  const githubUsername = githubUser.login

  // TODO: Replace this with your own DB query.
  const existingUser = await getUserFromGitHubId(githubUserId)

  if (existingUser !== null) {
    const sessionToken = generateSessionToken()
    const session = await createSession(sessionToken, existingUser.id)
    setSessionTokenCookie(context, sessionToken, session.expiresAt)
    return context.redirect('/')
  }

  // TODO: Replace this with your own DB query.
  const user = await createUser(githubUserId, githubUsername)

  const sessionToken = generateSessionToken()
  const session = await createSession(sessionToken, user.id)
  setSessionTokenCookie(context, sessionToken, session.expiresAt)
  return context.redirect('/account')
}
