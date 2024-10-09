1. signup

{
  "userEmail": "dennis.k@dev.null",
  "password": "1TrulySecurePassword!!"
}

2. login

{
  "email": "dennis.k@dev.null",
  "password": "1TrulySecurePassword!!"
}

.
.
.
.

4. passreset token

{
  "password": "1EvenM00rTrulySecurePassword!@?!"
}

5.pass-change
{
  "password": "1EvenM00rTrulySecurePassword!@?!"
}

.
.
.
.

10. Get-User-Entities (businesses)----- Financious <---

{
  "data": [
    {
      "id": "5eec78f4a880b7c235f66e7c",
      "accountNumber": "6838115832533278",
      "name": "New World Consulting Limited",
      "email": "nwcltd@dev.null",
      "telephone": "+442090807060",
      "incType": "ltd",
      "companyNumber": "A12345",
      "website": "https://neworco.null",
      "declaredTurnover": 20000,
      "description": "We show you how good things can be and what you need to do to make them happen.",
      "address": "123 Yellow Brick Road",
      "city": "Nairobi",
      "region": "Greater Nairobi",
      "postalCode": "76412 2ENG",
      "country": "England",
      "status": "pending",
      "showTagsMatchedSinceLastLogin": true,
      "receiveDailyMatchNotificationEmail": true,
      "offers": [
        "consulting",
        "system-design"
      ],
      "wants": [
        "it-services",
        "accounting"
      ],
      "categories": [
        "consulting"
      ],
      "balance": 0,
      "maxPositiveBalance": 500,
      "maxNegativeBalance": 0,
      "pendingTransfers": [
        {
          "id": "1dimEdDxOJYjDeaP6HEy4cIhLdD",
          "transfer": "in",
          "isInitiator": false,
          "accountNumber": "5211115451222517",
          "entityName": "Hipster Brews",
          "amount": 157.8,
          "description": "Payment of inv. 1234",
          "status": "transferInitiated",
          "dateProposed": "2020-06-23T12:42:57.786628Z"
        }
      ]
    }
  ]
}

11. result-for above

{
  "name": "New World Pizza PLC",
  "email": "nwpplc@dev.null",
  "telephone": "+442098765432",
  "incType": "plc",
  "companyNumber": "B67890",
  "website": "https://nwpizza.null",
  "declaredTurnover": 10000,
  "description": "We show you how good things can taste and where you need to go to eat them!",
  "address": "456 Yellow Brick Road",
  "city": "Nairobi",
  "region": "Greater Nairobi",
  "postalCode": "76412 00100",
  "country": "England",
  "showTagsMatchedSinceLastLogin": false,
  "receiveDailyMatchNotificationEmail": false,
  "offers": [
    "pizza",
    "wine"
  ],
  "wants": [
    "flour",
    "mozarella",
    "tomato"
  ]
}

12. Get-entity-cagetories

{
  "data": [
    {
      "id": "5ef1b5ccd35791b2dc384425",
      "name": "agriculture"
    },
    {
      "id": "5ef1b5ccd35791b2dc384426",
      "name": "cafes"
    },
    {
      "id": "5ef1b5ccd35791b2dc384427",
      "name": "cleaning-services"
    },
    {
      "id": "5ef1b5ccd35791b2dc384428",
      "name": "professional-services"
    },
    {
      "id": "5ef1b5ccd35791b2dc384429",
      "name": "restaurant"
    },
    {
      "id": "5ef1b5ccd35791b2dc384430",
      "name": "restaurant-supplies"
    },
    {
      "id": "5ef1b5ccd35791b2dc384431",
      "name": "sports-equipment"
    },
    {
      "id": "5ef1b5ccd35791b2dc384432",
      "name": "transport"
    }
  ],
  "meta": {
    "numberOfResults": 8,
    "totalPages": 1
  }
}

13. get-entity-TAGS (what a business offer)

{
  "data": [
    {
      "id": "5ef1b5ccd35791b2dc384410",
      "name": "art-supplies"
    },
    {
      "id": "5ef1b5ccd35791b2dc384411",
      "name": "catering"
    },
    {
      "id": "5ef1b5ccd35791b2dc384412",
      "name": "computer-repair"
    },
    {
      "id": "5ef1b5ccd35791b2dc384413",
      "name": "dog-sitting"
    },
    {
      "id": "5ef1b5ccd35791b2dc384414",
      "name": "micro-brewed-beer"
    },
    {
      "id": "5ef1b5ccd35791b2dc384415",
      "name": "organic-vegetables"
    },
    {
      "id": "5ef1b5ccd35791b2dc384416",
      "name": "piano-tuning"
    },
    {
      "id": "5ef1b5ccd35791b2dc384417",
      "name": "pizza"
    },
    {
      "id": "5ef1b5ccd35791b2dc384418",
      "name": "van-rental"
    },
    {
      "id": "5ef1b5ccd35791b2dc384419",
      "name": "well-being"
    }
  ],
  "meta": {
    "numberOfResults": 10,
    "totalPages": 1
  }
}