import { fileURLToPath, URL } from 'node:url';

import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import { createServer, type IncomingMessage, RequestListener, Server, ServerResponse } from 'node:http'

export const listings: any[] = [
  {
    title: 'Leo am Teich - Wohnen am Wasser',
    postalCode: 1010,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Leo am Teich - Wohnen am Wasser - Provisionsfrei!',
    postalCode: 1011,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Leo am Teich - Wohnen am Wasser',
    postalCode: 1010,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Leo am Teich - Wohnen am Wasser',
    postalCode: 1015,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: '2 Zimmer mit Küche und riesigem Balkon!',
    postalCode: 1013,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Martha im Grün - gefördertes Eigentum beim Badeteich',
    postalCode: 1012,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: '2 Zimmer mit Küche und riesigem Balkon!',
    postalCode: 1012,
    city: 'Vienna',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: '2 Zimmer mit Küche und riesigem grazer Balkon!',
    postalCode: 1013,
    city: 'Graz',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
  {
    title: 'Martha im Grün - gefördertes grazer Eigentum',
    postalCode: 1012,
    city: 'Graz',
    id: '1',
    _partitionKey: 'test',
    housingCooperative: 'test',
    projectId: 'test',
    listingId: 'test',
    country: 'test',
    address: 'test',
    roomCount: 0,
    squareMeters: 0,
    availabilityDate: 'test',
    yearBuilt: 0,
    hwgEnergyClass: 'test',
    fgeeEnergyClass: 'test',
    listingType: 'test',
    rentPricePerMonth: 0,
    cooperativeShare: 0,
    salePrice: 0,
    additionalFees: 0,
    detailsUrl: 'test',
    previewImageUrl: 'test',
    scraperId: 'test',
    createdAt: 'test',
    lastModifiedAt: 'test'
  },
]

const startDummyServer = () => {
  const requestListener: RequestListener = (req, res) => {
    res.setHeader('Content-Type', 'application/json');
    if (req.url?.match("/api/cities/.+/listings")) {
      const city = req.url.split('/')[3];
      console.log('Returning listings for city:', city)
      const filteredListings = listings.filter(listing => listing.city.toLowerCase() === city);

      res.writeHead(200);
      res.end(JSON.stringify({"results": filteredListings}));
    }
  };

  const serverPort = Math.floor(Math.random() * (3693 - 3333 + 1)) + 3333;

  const apiServer = createServer(requestListener);
  apiServer.listen(serverPort, () => {
    console.log(`Dummy data server is running on port ${serverPort}`);
  });

  return {server: apiServer, port: serverPort};
}

let dummy:  {server: Server<typeof IncomingMessage, typeof ServerResponse>, port: number} = null



export default defineConfig({
  plugins: [
    vue(),
    {
      name: 'integrated-dummy-server',
      configureServer(server) {
        dummy?.server.close();
        dummy = startDummyServer();

        server.middlewares.use((req, res, next) => {
          if (req.url.startsWith('/api')) {
            dummy.server.emit("request", req, res);
          } else {
            next();
          }
        });
      },

      closeBundle() {
        dummy?.server.close();
      }
    }
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
