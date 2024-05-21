export interface ApiListingsResponse {
  results: Listing[];
}

export interface Listing {
  id: string;
  _partitionKey: string;
  title: string;
  housingCooperative: string;
  projectId: string;
  listingId: string;
  country: string;
  city: string;
  postalCode: number;
  address: string;
  roomCount: number;
  squareMeters: number;
  availabilityDate: string;
  yearBuilt: number;
  hwgEnergyClass: string;
  fgeeEnergyClass: string;
  listingType: string;
  rentPricePerMonth: number;
  cooperativeShare: number;
  salePrice: number;
  additionalFees: number;
  detailsUrl: string;
  previewImageUrl: string;
  scraperId: string;
  createdAt: string;
  lastModifiedAt: string;
}
