const DEV_MODE = import.meta.env.DEV
const PUBLIC_STRIPE_API_KEY = import.meta.env.PUBLIC_STRIPE_API_KEY
const PRIVATE_STRIPE_API_KEY = import.meta.env.PRIVATE_STRIPE_KEY

// Found under "Product Catalog".
const PLUGIN_1_PRICE_ID = import.meta.env.PLUGIN_1_PRICE_ID

const BACKEND_URL = import.meta.env.BACKEND_URL

const GithubClientID = import.meta.env.GITHUB_CLIENT_ID
const GithubClientSecret = import.meta.env.GITHUB_CLIENT_SECRET

if (!PUBLIC_STRIPE_API_KEY) {
  throw new Error('PUBLIC Stripe API Key missing!')
}

if (!PRIVATE_STRIPE_API_KEY) {
  throw new Error('PRIVATE Stripe API Key missing!')
}

if (!PLUGIN_1_PRICE_ID) {
  throw new Error('PLUGIN_1 Price ID missing!')
}

if (!BACKEND_URL) {
  throw new Error('BACKEND_URL is missing!')
}

if (!GithubClientID) {
  throw new Error('Github Client ID missing!')
}
if (!GithubClientSecret) {
  throw new Error('Github Client Secret missing!')
}

export {
  DEV_MODE,
  PUBLIC_STRIPE_API_KEY,
  PRIVATE_STRIPE_API_KEY,
  PLUGIN_1_PRICE_ID,
  GithubClientID,
  GithubClientSecret,
}
