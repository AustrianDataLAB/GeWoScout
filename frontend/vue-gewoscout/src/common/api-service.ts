import type { ApiListingsResponse, Listing } from '@/types/ApiResponseListings';
import type SearchInputs from '@/types/SearchInputs';
import axios from 'axios';

// flatpak run org.chromium.Chromium --disable-site-isolation-trials --disable-web-security --user-data-dir="~/chromiumteest"

export async function getListings(searchInputs: SearchInputs): Promise<Listing[]> {
  try {
    const response = await axios.get(`/api/cities/${searchInputs.city}/listings`);

    console.log(response);
    console.log(response.data);

    const listings: ApiListingsResponse = response.data;
    return listings.results;
  } catch (error) {
    console.error(error);

    return [];
  }
}
