import { createServer } from 'http';

const listings = [
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
  }
];

let loggedIn = false;

const server = createServer((req, res) => {
  res.writeHead(200, { 'Content-Type': 'application/json' });
  if (req.url.match("/api/cities/.+/listings/?")) {
    console.log(req.url);
    const city = req.url.split("/")[3];
    res.end(JSON.stringify({
      results: listings.filter(listing => listing.city.toLowerCase() === city)
    }));
    return;
  }
  if (req.url.match("/.auth/me/?")) {
    console.log(req.url);
    if (loggedIn === false) {
      res.end(JSON.stringify({
        "clientPrincipal": null
      }));
      return;
    }

    res.end(JSON.stringify({
      "clientPrincipal": {
        "identityProvider": "aad",
        "userId": "a220ff1e00d94af8a8f6a3c7da71e491",
        "userDetails": "someuser@austrianopencloudcommunity.onmicrosoft.com",
        "userRoles": [
          "anonymous",
          "authenticated"
        ]
      }
    }));
    return;
  }

  if (req.url.match("/.auth/login/aad/?")) {
    console.log(req.url);
    loggedIn = true;
    res.end("OK");
    return;
  }

  if (req.url.match("/.auth/logout/?")) {
    console.log(req.url);
    loggedIn = false;
    res.end("OK");
    return;
  }

  res.end(JSON.stringify({results: listings}));
});

const PORT = 3333;
server.listen(PORT, () => {
  console.log(`Server is listening on port ${PORT}`);
});
