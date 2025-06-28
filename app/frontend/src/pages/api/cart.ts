import type { APIContext } from 'astro'

export async function GET({ session }: APIContext) {
  const cart = await session?.get('cart')
  return Response.json({ cart })
}
