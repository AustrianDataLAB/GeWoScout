import type { ApiListingsResponse, Listing } from '@/types/ApiResponseListings';
import { EnergyClass, Type } from '@/types/Enums';
import type SearchInputs from '@/types/SearchInputs';
import axios from 'axios';

// flatpak run org.chromium.Chromium --disable-site-isolation-trials --disable-web-security --user-data-dir="~/chromiumteest"

interface ListingsParams {
  listingType: string;
  housingCooperative?: string; // geno
  minRentPricePerMonth?: number;
  maxRentPricePerMonth?: number;
  postalCode?: string;
  minRoomCount?: number;
  maxRoomCount?: number;
  minSqm?: number;
  maxSqm?: number;
  availableFrom?: string;
  minYearBuilt?: number;
  maxYearBuilt?: number;
  minHwgEnergyClass?: string;
  minFgeeEnergyClass?: string;

  minCooperativeShare?: number;
  maxCooperativeShare?: number;
  minSalePrice?: number;
  maxSalePrice?: number;
}

export async function getListings(searchInputs: SearchInputs): Promise<Listing[]> {
  const params: ListingsParams = {
    listingType: Type[searchInputs.listingType]
  };

  if (searchInputs.housingCooperative !== '') {
    params.housingCooperative = searchInputs.housingCooperative;
  }
  if (searchInputs.postalCode !== '') {
    params.postalCode = searchInputs.postalCode;
  }
  if (searchInputs.minRoomCount !== null) {
    params.minRoomCount = searchInputs.minRoomCount;
  }
  if (searchInputs.maxRoomCount !== null) {
    params.maxRoomCount = searchInputs.maxRoomCount;
  }
  if (searchInputs.minSqm !== null) {
    params.minSqm = searchInputs.minSqm;
  }
  if (searchInputs.maxSqm !== null) {
    params.maxSqm = searchInputs.maxSqm;
  }
  if (searchInputs.availableFrom !== null) {
    params.availableFrom = searchInputs.availableFrom.toISOString();
  }
  if (searchInputs.minYearBuilt !== null) {
    params.minYearBuilt = searchInputs.minYearBuilt;
  }
  if (searchInputs.maxYearBuilt !== null) {
    params.maxYearBuilt = searchInputs.maxYearBuilt;
  }
  if (searchInputs.minRentPricePerMonth !== null) {
    params.minRentPricePerMonth = searchInputs.minRentPricePerMonth;
  }
  if (searchInputs.maxRentPricePerMonth !== null) {
    params.maxRentPricePerMonth = searchInputs.maxRentPricePerMonth;
  }
  if (searchInputs.minCooperativeShare !== null) {
    params.minCooperativeShare = searchInputs.minCooperativeShare;
  }
  if (searchInputs.maxCooperativeShare !== null) {
    params.maxCooperativeShare = searchInputs.maxCooperativeShare;
  }
  if (searchInputs.minSalePrice !== null) {
    params.minSalePrice = searchInputs.minSalePrice;
  }
  if (searchInputs.maxSalePrice !== null) {
    params.maxSalePrice = searchInputs.maxSalePrice;
  }
  if (searchInputs.minHwgEnergyClass !== null) {
    params.minHwgEnergyClass = EnergyClass[searchInputs.minHwgEnergyClass];
  }
  if (searchInputs.minFgeeEnergyClass !== null) {
    params.minFgeeEnergyClass = EnergyClass[searchInputs.minFgeeEnergyClass];
  }

  console.log(params);

  try {
    const response = await axios.get(`/api/cities/${searchInputs.city}/listings`, {
      params: params
    });

    console.log(response);
    console.log(response.data);

    const listings: ApiListingsResponse = response.data;
    return listings.results;
  } catch (error) {
    console.error(error);

    return [];
  }
}
