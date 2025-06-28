import { defineMiddleware, sequence } from 'astro:middleware'

const auth = defineMiddleware((context, next) => {
  // Check if the user is logged in.
  return next()
})

const setLastVisit = defineMiddleware(async (context, next) => {
  context.session?.set('lastVisit', new Date())
  return next()
})

export const onRequest = sequence(auth, setLastVisit)
