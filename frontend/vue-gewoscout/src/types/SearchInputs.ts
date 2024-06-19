import type { EnergyClass, Type } from './Enums';

export default interface SearchInputs {
  city: string;
  housingCooperative: string; // geno
  postalCode: string;
  minRoomCount: number | null;
  maxRoomCount: number | null;
  minSqm: number | null;
  maxSqm: number | null;
  availableFrom: Date | null;
  minYearBuilt: number | null;
  maxYearBuilt: number | null;
  minHwgEnergyClass: EnergyClass | null;
  minFgeeEnergyClass: EnergyClass | null;
  listingType: Type;
  minRentPricePerMonth: number | null;
  maxRentPricePerMonth: number | null;
  minCooperativeShare: number | null;
  maxCooperativeShare: number | null;
  minSalePrice: number | null;
  maxSalePrice: number | null;
}
