export default interface ApiResponsePreferences {
  _partitionKey: string;
  id: string;
  title: string;
  projectId: string;
  availableFrom: string;
  city: string;
  email: string;
  housingCooperative: string;
  listingType: string;
  maxCooperativeShare?: number;
  maxRoomCount?: number;
  maxYearBuilt?: number;
  minCooperativeShare?: number;
  minFgeeEnergyClass: string;
  minHwgEnergyClass: string;
  minRoomCount?: number;
  minSalePrice?: number;
  minYearBuilt?: number;
  postalCode?: string;
  roomCount?: number;
  minRentPrice?: number;
  maxRentPrice?: number;
  minSqm?: number;
  maxSqm?: number;
  maxSalePrice?: number;
}
