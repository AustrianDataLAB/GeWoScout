import type { ApiListingsResponse, Listing } from "@/types/ApiResponseListings";
import axios from "axios";

// flatpak run org.chromium.Chromium --disable-site-isolation-trials --disable-web-security --user-data-dir="~/chromiumteest"


export async function getListings(searchCity: string): Promise<Listing[]> {
    try {
        const response = await axios.get(`/api/cities/${searchCity}/listings`);
        // const response = await axios.get('https://nice-pebble-01f706603-development.westeurope.5.azurestaticapps.net/api/health');
        console.log(response);
        console.log(response.data);

        const listings: ApiListingsResponse = response.data;

        return listings.results;
    } catch (error) {
        console.error(error);

        return [];
    }
}
