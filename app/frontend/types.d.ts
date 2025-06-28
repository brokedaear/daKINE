interface User {
  id: number;
  githubId: number;
  email: string;
}

interface CartItem {
  id: string;
  name: string;
  count: number;
  price: number;
}

type ItemID = string

interface ShoppingCart {
  items: Map<ItemID, CartItem>;
}
