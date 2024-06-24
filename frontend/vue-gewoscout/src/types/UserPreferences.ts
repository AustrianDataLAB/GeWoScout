import type { EnergyClass, Type } from './Enums';

export default interface UserPreferences {
  availableFrom: Date | null;
  city: string | null;
  email: string | null;
  housingCooperative: string | null;
  listingType: Type;
  maxCooperativeShare: number | null;
  maxRentPrice: number | null;
  maxRoomCount: number | null;
  maxSalePrice: number | null;
  maxSqm: number | null;
  maxYearBuilt: number | null;
  minCooperativeShare: number | null;
  minFgeeEnergyClass: EnergyClass;
  minHwgEnergyClass: EnergyClass;
  minRentPrice: number | null;
  minRoomCount: number | null;
  minSalePrice: number | null;
  minSqm: number | null;
  minYearBuilt: number | null;
  postalCode: string | null;
  // projectId: string;
  // roomCount: number;
  // title: string;
}
