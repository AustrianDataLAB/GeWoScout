import type { ApiListingsResponse, Listing } from '@/types/ApiResponseListings';
import { EnergyClass, Type } from '@/types/Enums';
import type SearchInputs from '@/types/SearchInputs';
import type UserPreferences from '@/types/UserPreferences';
import axios from 'axios';

// flatpak run org.chromium.Chromium --disable-site-isolation-trials --disable-web-security --user-data-dir="~/chromiumteest"

interface ListingsParams {
  listingType: string;
  housingCooperative?: string; // geno
  minRentPrice?: number;
  maxRentPrice?: number;
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
  try {
    const response = await axios.get(`/api/cities/${searchInputs.city}/listings`, {
      params: convertInputToOptionalParamsObj(searchInputs)
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

export async function getUserPreferences(): Promise<UserPreferences[]> {
  try {
    const response = await axios.get('/api/users/preferences');

    const preferences: UserPreferences[] = response.data;
    // TODO probably type fix needed
    return preferences;
  } catch (error) {
    console.error(error);
    return [];
  }
}

export async function setUserPreferences(preferences: UserPreferences): Promise<boolean> {
  try {
    const requestUserPrefs = {
      email: '',
      title: '',
      housingCooperative: '',
      projectId: '',
      postalCode: '',
      roomCount: 1,
      minRoomCount: 1,
      maxRoomCount: 1,
      minSquareMeters: 1,
      maxSquareMeters: 1,
      availableFrom: '2010-01-01',
      minYearBuilt: 1901,
      maxYearBuilt: 1901,
      minHwgEnergyClass: 'F',
      minFgeeEnergyClass: 'F',
      listingType: 'both',
      minRentPricePerMonth: 1,
      maxRentPricePerMonth: 1,
      minCooperativeShare: 1,
      maxCooperativeShare: 1,
      minSalePrice: 1,
      maxSalePrice: null
    };
    const response = await axios.put(
      `/api/users/preferences/${preferences.city}`,
      requestUserPrefs
    );

    console.log(response);
    console.log(response.data);

    return true;
  } catch (error) {
    console.error(error);

    return false;
  }
}

function convertInputToOptionalParamsObj(input: SearchInputs): ListingsParams {
  const params: ListingsParams = {
    listingType: Type[input.listingType]
  };

  if (input.housingCooperative !== '') {
    params.housingCooperative = input.housingCooperative;
  }
  if (input.postalCode !== '') {
    params.postalCode = input.postalCode;
  }
  if (input.minRoomCount !== null) {
    params.minRoomCount = input.minRoomCount;
  }
  if (input.maxRoomCount !== null) {
    params.maxRoomCount = input.maxRoomCount;
  }
  if (input.minSqm !== null) {
    params.minSqm = input.minSqm;
  }
  if (input.maxSqm !== null) {
    params.maxSqm = input.maxSqm;
  }
  if (input.availableFrom !== null) {
    params.availableFrom = input.availableFrom.toISOString();
  }
  if (input.minYearBuilt !== null) {
    params.minYearBuilt = input.minYearBuilt;
  }
  if (input.maxYearBuilt !== null) {
    params.maxYearBuilt = input.maxYearBuilt;
  }
  if (input.minRentPrice !== null) {
    params.minRentPrice = input.minRentPrice;
  }
  if (input.maxRentPrice !== null) {
    params.maxRentPrice = input.maxRentPrice;
  }
  if (input.minCooperativeShare !== null) {
    params.minCooperativeShare = input.minCooperativeShare;
  }
  if (input.maxCooperativeShare !== null) {
    params.maxCooperativeShare = input.maxCooperativeShare;
  }
  if (input.minSalePrice !== null) {
    params.minSalePrice = input.minSalePrice;
  }
  if (input.maxSalePrice !== null) {
    params.maxSalePrice = input.maxSalePrice;
  }
  if (input.minHwgEnergyClass !== null) {
    params.minHwgEnergyClass = EnergyClass[input.minHwgEnergyClass];
  }
  if (input.minFgeeEnergyClass !== null) {
    params.minFgeeEnergyClass = EnergyClass[input.minFgeeEnergyClass];
  }

  console.log(params);
  return params;
}
