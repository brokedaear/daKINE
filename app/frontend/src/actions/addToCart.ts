import { defineAction } from 'astro:actions'
import { z } from 'astro:schema'

export const server = {
  addToCart: defineAction({
    input: z.object({
      productId: z.string(),
      productName: z.string(),
      productPrice: z.number(),
    }),
    handler: async (input, context) => {
      const cart = await context.session?.get('shoppingCart')
      if (cart) {
        if (cart.items.has(input.productId)) {
          const c = cart.items.get(input.productId)
          c!.count += 1
          cart.items.set(input.productId, c!)
        } else {
          const item: CartItem = {
            id: input.productId,
            name: input.productName,
            count: 1,
            price: input.productPrice,
          }
          cart.items.set(input.productId, item)
        }
        context.session?.set('shoppingCart', cart)
        return cart
      }
      // Redirect the user to login if not logged in.
      redirect()
    },
  }),
}
