export interface Search {
  id: string;
  name: string;
  url: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Property {
  id: string;
  searchId: string;
  externalId: string;
  name: string;
  address: string;
  city: string;
  state: string;
  zip: string;
  propertyType: string;
  price: number | null;
  sizeSqFt: number | null;
  description: string;
  url: string;
  imageUrl: string;
  listedDate: string | null;
  scrapedAt: string;
  createdAt: string;
}

export interface SearchFormData {
  name: string;
  url: string;
  active: boolean;
}
