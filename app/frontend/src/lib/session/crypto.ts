import { generateRandomString, type RandomReader } from '@oslojs/crypto/random'
import { encodeBase32 } from '@oslojs/encoding'

const customCharSet = '0123456789ABCDEFGHJKMNPQRSTVWXYZ'

/**
 * Creates a base 32 encoded string. Feed a random string of 15 bytes length
 * into this function to generate a 24-character token.
 * @param s string to encode
 * @returns base 32 encoded string
 */
export function encodedString(s: string): string {
  const data: Uint8Array = new TextEncoder().encode(s)
  return encodeBase32(data)
}

/**
 * Creates a random string using `@oslojs` and a random reader. This function
 * can be used to as a first step to generate JWTs, tokens, those of that ilk.
 * It is recommended to use a `len` of 15, and so the default value is 15.
 * @param len length of bytes to generate
 * @returns string
 */
export function randomString(len: number = 15): string {
  return generateRandomString(random, customCharSet, len)
}

const random: RandomReader = {
  read(bytes: Uint8Array): void {
    // Not available in Node versions < 20
    crypto.getRandomValues(bytes)
  },
}
