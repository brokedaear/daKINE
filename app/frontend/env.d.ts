interface ImportMetaEnv {
  STRIPE_PUBLIC_API_KEY: string;
  BACKEND_SERVICE_URL: string;
  GITHUB_CLIENT_ID: string;
  GITHUB_CLIENT_SECRET: string;
  REDIS_URL: string;
}

interface ImportMeta {
  env: ImportMetaEnv;
}

declare namespace App {
  interface SessionData {
    user: User;
    shoppingCart: ShoppingCart;
  }
}
