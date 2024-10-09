
HOMEPESA-MC Beta User API
 
 1.
 
Introduction
The HOMEPESA-MC Beta User API v1 exposes all of the end user functionality currently available in the Beta version of the HOMEPESA-MC API.

By providing an API, developers who want to create their own front-end user interface for HOMEPESA-MC will have significant flexibility to implement it in whatever way they choose. This means developers can present HOMEPESA-MC in any language, setup their own signup flow, optimize it for whatever devices their users prefer, develop a mobile app, integrate other services such as chat, etc.

Importantly, an API enables developers to integrate HOMEPESA-MC functionality directly into their own apps (e.g., import transfer data into an accounting application, instruct mutual credit transfers from an e-wallet application, interact with HOMEPESA-MC data via a chat bot, etc.).

Background
To understand how the HOMEPESA-MC Beta API works, please read the following design documents:

HOMEPESA-MC Beta Data Model
HOMEPESA-MC Beta Functionality
Test Server
These API docs assume you are running the API server on your local machine using Docker & Docker Compose. See the How to Start instructions in the project's GitHub repo for more details.

License
MIT License
Servers

https://mccs.homesako.com/api/v1
Manage Account
Create and manage a user and its linked entity



POST
/signup
Create a new account - a user resource and associated entity resource
Individual users can create an account in HOMEPESA-MC by providing an email address, creating a password and adding some other details about themselves and their "business". A "business" need not be a formally established business; it could simply be a list of their skills that they are willing to offer to other participants in the network, or another type of entity such as an association, not-for-profit, NGO, etc.

We use the term entity to generically identify businesses, non-profits, NGOs, sole proprietorships, limited companies, etc.

The POST /signup API call creates a user resource and an entity (e.g., a business) resource and links the two together. The user resource references its associated entity, and the entity resource references its associated user. These two related resources make up what is referred to as an account.

An account can be created with only a unique email address and a password, enabling a new user to get started quickly. The userID and entityID are returned in the response body, along with a JSON Web Token to authorize the user of the account for further API calls related to the user and entity resources just created.

All other details (entity name, user first/last names, goods/services offered or wanted, etc.) are optional and can be added later. The decision about customer data requirements for signup can be determined by developers implementing a front end to the API, since they will be working with the organizations implementing HOMEPESA-MC. If other information is collected during the initial registration process, these extra details can be passed immediately in this POST /signup API endpoint or later using the PATCH /user and/or PATCH /user/entities/{entityID} API endpoints.

The current (v1) implementation of the API has a one-to-one relationship between the user and the entity. This will be changed to a many-to-many relationship in a future implementation.



Request body

`application/json`

Signup request body

Examples: 
requiredFields
Example Value


     Schema
{
  "userEmail": "test@gmail.com",
  "password": "Test321321$$"
}

Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Example Value


     Schema
{
  "data": {
    "userID": "5eec78f4a880b7c235f66e80",
    "entityID": "5eec78f4a880b7c235f66e7c",
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI2NDIxNjQsInVzZXJJRCI6IjVlZWM3OGY0YTg4MGI3YzIzNWY2NmU4MCIsImFkbWluIjpmYWxzZX0.WhPMfV9S-xgOHcuMT_K-fBhR_K6MpXiJ15GYn4Jz7im1dvhwnV2bEnwuWeFUockl45StxguvIOA5qJ-_3xA14CuP0wJbZa3hVH4jnYXonHlCyHDB8w67RLN9IMFGnSEshhh4D3RjQVpEpBm7jLhQcHKOSQIqUU_RfPkiNxpUkDI6t1RW_-rhY4UsTTuxnC5SOeajzOgiDFM4NwJfjebys8xDGTqYoi4dpCJZEtD_U_X9BuEOovRRJo0TY6m76XxUB9J5U_Hfjm7k_A3aLv3WgDScRv_k-LSsOvviGk1A2ct0nQ2RaVY1udA-76rv1xbvSd26Xds2XtrPb_SzUL-J8A"
  }
}

400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


POST

`/login`

Login and receive an authorization token
A user will need to authenticate with an email/password combination which, if successful, will result in a JSON Web Token (JWT) being passed to the user to use with each API request that requires authorization.

The JWT will expire after 24 hours, or as soon as the user calls the /logout API.



Request body

`application/json`

A JSON object containing an email address and a password

Example Value


     Schema
{
  "email": "test@gmail.com",
  "password": "Test321321$$"
}

<<<--------------------------------->>>

Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Example Value


     Schema
{
  "data": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI5ODU0MzcsInVzZXJJRCI6IjVlZjFiNWQ5NTA2ODQ3ZWE5MTRkYjA2NiIsImFkbWluIjpmYWxzZX0.E3YJWpUN0h2LS3IfE2niYclTnhKUGXDay7SR_VeVzaq4a_lZPs2w4yxFUWbbeAcPeMgtcd3bEoX6PtDhW3hjJ8XVbqmkJQRmh3uw5ULfgzkIzrQw3twaG8TO6ARM4UvKJWz2Wqb6czd5SesJmM2htIuY3DBJ_u4r9x3hshM5_0kHalfEZtQvae0KrJ2_eBjPmFOdON62QungzStkjTKTfqsvystwFSdfOwAltg0Nri1Z6q-E9AZMTnAxzKNqjp4Ja3hX1IoZXPiV8F0-yl1PhrKI_YzP57pOt84T72_WK3z1_hofHAOwby5-2qvtoKWSmxqzKsIYUYeG89W7h4TvBQ"
  }
}

400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


POST
/logout
Logout and expire the authorization token

Logging out will immediately expire the JWT currently associated to the user's account.



Responses
Code  Description 200 OK


401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


POST
/password-reset
Request a password reset token
When a password is lost by the user, a new one can be requested by providing the email address associated with the user's account. A password reset token will then be sent to that email address.



Request body

`application/json`

A JSON object containing an email address

Example Value


     Schema
{
  "email": "dennis-kip@dev.null"
}
Responses
Code  Description 200 OK


400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


POST
/password-reset/{token}
Reset a password
The new password can be set by providing it along with the password reset token that was received at the user's email address.

Parameters
Name  Description
token *
string
(path)
The password reset token.

token
Request body

`application/json`

A JSON object containing a password

Example Value


     Schema
{
  "password": "1EvenM00rTrulySecurePassword!@?!"
}
Responses
Code  Description 200 OK


400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


POST
/password-change
Change the password

A logged-in user can change the password by sending the new password along with the JWT.



Request body

`application/json`

A JSON object containing a password

Example Value


     Schema
{
  "password": "1EvenM00rTrulySecurePassword!@?!"
}
Responses
Code  Description 200 OK

400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


GET
/user
View a user's own details


PATCH
/user
Modify a user's own details


GET
/user/entities
View an entity's details

Users can request the details of their linked entity as recorded in the HOMEPESA-MC database.



Responses
Code  Description 200 OK

`application/json`

Controls Accept header.
Example Value


     Schema
{
  "data": [
    {
      "id": "5eec78f4a880b7c235f66e7c",
      "accountNumber": "6838115832533278",
      "name": "New World Consulting Limited",
      "email": "nwcltd@dev.null",
      "telephone": "+2542090807060",
      "incType": "ltd",
      "companyNumber": "A12345",
      "website": "https://neworco.null",
      "declaredTurnover": 20000,
      "description": "We show you how good things can be and what you need to do to make them happen.",
      "address": "123 Yellow Brick Road",
      "city": "Nairobi",
      "region": "LOWER-Nairobi",
      "postalCode": "00508",
      "country": "Nairobi-kenya",
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
          "dateProposed": "2023-07-23T12:42:57.786628Z"
        }
      ]
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


PATCH
/user/entities/{entityID}
Modify an entity's details

Users can change their entity's details in the HOMEPESA-MC database, except for the status field which can only be changed by an administrator. The id and accountNumber fields are system-generated and therefore are not changeable by either users or administrators.

Changes to the offers and wants array fields must include all relevant tags because they will overwrite the arrays already stored in the database.

The email for an entity is separate from the linked user's email, although the user email address is set for the entity's email as well when the account is first created, but only if an entity email is not specified when creating an account at the POST /signup endpoint. The entity email, which receives notifications, can be changed by the user.

Parameters
Name  Description
entityID *
string
(path)
The unique entity ID

5eec78f4a880b7c235f66e7c
Request body

`application/json`

The entity fields a user wants to update

Example Value


     Schema
{
  "name": "New World Pizza PLC",
  "email": "nwpplc@dev.null",
  "telephone": "+2542098765432",
  "incType": "plc",
  "companyNumber": "B67890",
  "website": "https://nwpizza.null",
  "declaredTurnover": 10000,
  "description": "We show you how good things can taste and where you need to go to eat them!",
  "address": "456 Yellow Brick Road",
  "city": "Nairobi",
  "region": "LOWER-Nairobi",
  "postalCode": "00508",
  "country": "Nairobi-kenya",
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
Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Example Value


     Schema
{
  "data": {
    "id": "5eec78f4a880b7c235f66e7c",
    "accountNumber": "6838115832533278",
    "name": "New World Pizza PLC",
    "email": "nwpplc@dev.null",
    "telephone": "+2542098765432",
    "incType": "plc",
    "companyNumber": "B67890",
    "website": "https://nwpizza.null",
    "declaredTurnover": 10000,
    "description": "We show you how good things can taste and where you need to go to eat them!",
    "address": "456 Yellow Brick Road",
    "city": "Nairobi",
    "region": "LOWER-Nairobi",
    "postalCode": "00508",
    "country": "Nairobi-kenya",
    "status": "pending",
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
        "dateProposed": "2023-07-23T12:42:57.786628Z"
      }
    ]
  }
}

400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

403 
User does not have permission to access the resource.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Permission denied."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}

Find Entities
Search and filter entities



GET
/categories
Get a full or partial list of entity categories
The directory of entities is split up into categories that are manually assigned to each entity by an administrator. The entire list or just a subset of it by first letter(s) (prefix) and/or partial match (fragment) can be requested.

Searching for prefix "t" or "tr" returns all categories that begin with "t" or "tr" (e.g., "Teas" & "Transport" for "t", or only "Transport" for "tr").

Searching for fragment "tr" returns all categories that have "tr" anywhere in their name, including the beginning (e.g., "Transport" & "Carpentry").

Parameters
Name  Description
prefix
string
(query)
first letter of a category

Example : c

c
fragment
string
(query)
partial match of word/characters

Example : sport

sport
page
integer
(query)
The page number

Default value : 1

1
page_size
integer
(query)
minimum: 1
maximum: 100
The number of results per page

Default value : 10

10
Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Examples

fullList
Example Value


     Schema
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

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


GET
/tags
Get a full or partial list of tags used for offers and wants
Tags are words or short phrases that describe the goods or services an entity can provide (offers) to other entities or needs (wants) from other entities, in order to facilitate trades with them. For example, a Chinese restaurant might use offer tags such as "restaurant", "chinese", "take-out", "dim-sum", "delivery" to describe its service.

Up to 10 offer tags and 10 want tags can be specified per entity. A list of tags that are fuzzy, partial or exact matches to the input (fragment) provided can be requested. If no fragment is provided all tags are returned.

Parameters
Name  Description
fragment
string
(query)
partial match of word/characters

Example : be

be
page
integer
(query)
The page number

Default value : 1

1
page_size
integer
(query)
minimum: 1
maximum: 100
The number of results per page

Default value : 10

10
Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Examples

fullList
Example Value


     Schema
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

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


GET
/entities
Get a list of entities

Entities can be searched based on (1) their offers and wants tags, (2) the category they are listed under, (3) the date the offers and wants tags were added to their profile, (4) their entity name, and/or (5) an entity's selection of favorites.

If the request does not contain a JWT along with a querying_entity_id specified, the search for favorites functionality will not work (e.g., all results will show the isFavorite flag as false).

If a querying_entity_id is specified and both the requesting entity and the entity returned in the search are tradingAccepted status, the email address of the searched entity will also be included.

Parameters
Name  Description
offers
string
(query)
A list of goods/services offered by an entity

Example : pizza,pasta

pizza,pasta
wants
string
(query)
A list of good/services wanted by an entity

Example : vegetables

vegetables
category
string
(query)
A list of entities by category can be retrieved in the search functionality

restaurant
tagged_since
string
(query)
Get a list of entities that have had specified offers or wants tags added since a specific date and time

2019-12-25T12:12:12.001Z
name
string
(query)
A full or partial name of an entity can be searched

Alice's Restau
favorites_only
boolean
(query)
Show Favorites Only

Default value : false


false
querying_entity_id
string
(query)
The entity ID to which the filter is applied (requires user to be logged in)

5e561916ca06e1c8596eee9e
page
integer
(query)
The page number

Default value : 1

1
page_size
integer
(query)
minimum: 1
maximum: 100
The number of results per page

Default value : 10

10
Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Examples

withoutAuth
Example Value


     Schema
{
  "data": [
    {
      "id": "5eec78f4a880b7c235f66e7c",
      "accountNumber": "6838115832533278",
      "name": "New World Pizza PLC",
      "telephone": "+2542098765432",
      "incType": "plc",
      "companyNumber": "B67890",
      "website": "https://nwpizza.null",
      "declaredTurnover": 10000,
      "description": "We show you how good things can taste and where you need to go to eat them!",
      "address": "456 Yellow Brick Road",
      "city": "Nairobi",
      "region": "LOWER-Nairobi",
      "postalCode": "00508",
      "country": "Nairobi-kenya",
      "status": "pending",
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
      ],
      "categories": [
        "restaurant"
      ],
      "isFavorite": false
    }
  ],
  "meta": {
    "numberOfResults": 1,
    "totalPages": 1
  }
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

403 
User does not have permission to access the resource.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Permission denied."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


GET
/entities/{entityID}
Get a single entity

Returns a single entity's details using its ID. Requests can be made without a JWT from an authenticated user (see GET /entities above for more information).

Parameters
Name  Description
entityID *
string
(path)
The unique entity ID

5eec78f4a880b7c235f66e7c
querying_entity_id
string
(query)
The entity ID to which the filter is applied (requires user to be logged in)

5e561916ca06e1c8596eee9e
Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Examples

withoutAuth
Example Value


     Schema
{
  "data": {
    "id": "5eec78f4a880b7c235f66e7c",
    "accountNumber": "6838115832533278",
    "name": "New World Pizza PLC",
    "telephone": "+2542098765432",
    "incType": "plc",
    "companyNumber": "B67890",
    "website": "https://nwpizza.null",
    "declaredTurnover": 10000,
    "description": "We show you how good things can taste and where you need to go to eat them!",
    "address": "456 Yellow Brick Road",
    "city": "Nairobi",
    "region": "LOWER-Nairobi",
    "postalCode": "00508",
    "country": "Nairobi-kenya",
    "status": "pending",
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
    ],
    "categories": [
      "restaurant"
    ],
    "isFavorite": false
  }
}

400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


POST
/favorites
Create and manage a list of favorite entities

A user can toggle any entity as a favorite, which makes it easy to retrieve using the favorites_only search parameter described above (see GET /entities).



Request body

`application/json`

Set or unset a favorite entity

Example Value


     Schema
{
  "addToEntityID": "5de8dea0bdb7911205c0a6d7",
  "favoriteEntityID": "5e53e2e34b7e2bd4030e72ce",
  "isFavorite": true
}
Responses
Code  Description 200 OK


400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

403 
User does not have permission to access the resource.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Permission denied."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


POST
/send-email
Send a message to an entity

A user can send a message by email to an entity without seeing the email address of the receiving entity. This enables the receiving entity to keep its email private until the entity operator decides to reply to the sender and continue the conversation by email.

Entities with status accepted, tradingPending, tradingAccepted and tradingRejected can send an email to other entities with these same 4 statuses.

Entities with status pending or rejected cannot send an email to any other entity. They also cannot receive emails because there is no way to find their ID since they do not show up in search results because of their pending/rejected status.



Request body

`application/json`

Send an email to an entity

Example Value


     Schema
{
  "senderEntityID": "5de8dea0bdb7911205c0a6d7",
  "receiverEntityID": "5de8dea0bdb7911205c48d42",
  "body": "This is the email message body."
}
Responses
Code  Description 200 OK


400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

403 
User does not have permission to access the resource.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Permission denied."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}

Transfer Credits
Initiate and authorize mutual credit transfers



POST
/transfers
Initiate a transfer

A user can initiate a transfer out of or into the account of its entity, which must then be approved or rejected by the user operating the receiving entity, whose account will be credited or debited accordingly. Both entities must have tradingAccepted status in order to set up a transfer between them.

If the transfer parameter is set to out, the initiator will create a transfer that will debit funds from the initiator's entity's account. If transfer is in, the initiator will create a transfer that results in funds being credited to the initiator's entity's account. Either way, the transfer must be approved by the receiver (see PATCH /transfers/{transferID}) in order for the inbound or outbound transfer to move to or from the receiver's entity's account.



Request body

`application/json`

Initiate a transfer to/from entity's account from/to another entity

Example Value


     Schema
{
  "transfer": "out",
  "initiator": "7132460355005184",
  "receiver": "1234567887654321",
  "amount": 1.1,
  "description": "Here is the payment"
}
Responses


     Schema
{
  "id": "1ZceiVuQyGqeUYlC6UIKgEnaBkD",
  "from": "7132460355005184",
  "to": "0382855564717143",
  "amount": 177.5,
  "description": "Payment of your invoice number 12345",
  "status": "transferInitiated",
  "dateProposed": "2023-07-05T14:09:17.446965528Z"
}

400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

403 
User does not have permission to access the resource.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Permission denied."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


PATCH
/transfers/{transferID}
Confirm or cancel a transfer

The receiver can either accept or reject the transfer by specifying it in the action parameter.

The initiator of the transfer can cancel the transfer before the receiver has accepted or rejected it.

If a transfer is rejected or cancelled, a cancellationReason can be provided so that the other party understands why the initiator or receiver cancelled or rejected the transfer.

Parameters
Name  Description
transferID *
string
(path)
The unique transfer ID

1UZ7G7qJrIlwpVK9iSPXgx0A2xN
Request body

`application/json`

Examples: 
accept
Example Value


     Schema
{
  "action": "accept"
}
Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Examples

acceptResponse
Example Value


     Schema
{
  "id": "1UZ7G7qJrIlwpVK9iSPXgx0A2xN",
  "transfer": "out",
  "isInitiator": true,
  "accountNumber": "1234567887654321",
  "entityName": "Rhynyx",
  "amount": "1.1,",
  "description": "Payment of invoice number 12345",
  "status": "transferCompleted",
  "dateProposed": "2019-12-25T12:12:12.123Z",
  "dateCompleted": "2019-12-26T13:13:13.456Z"
}

400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}

Review Transfer Activity
View pending and completed mutual credit transfers



GET
/transfers
Get a list of transfers

A user can request a list of mutual credit transfers for the account of the entity. Transfers can be filtered by status (all, initiated, completed or cancelled).

The querying_entity_id is the ID of the entity whose account the information is being requested for. The user requesting must be associated with that entity or no information will be returned.

Parameters
Name  Description
status *
string
(query)
The status of the transfer

Available values : all, initiated, completed, cancelled


all
querying_entity_id *
string
(query)
The entity ID to which the account is linked

5e561916ca06e1c8596eee9e
page
integer
(query)
The page number

Default value : 1

1
page_size
integer
(query)
minimum: 1
maximum: 100
The number of results per page

Default value : 10

10
Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Examples

completed
Example Value


     Schema
{
  "data": [
    {
      "id": "1UZ7G7qJrIlwpVK9iSPXgx0A2xN",
      "transfer": "out",
      "isInitiator": true,
      "accountNumber": "1234567887654321",
      "entityName": "Rhynyx",
      "amount": "1.1,",
      "description": "Payment of invoice number 12345",
      "status": "transferCompleted",
      "dateProposed": "2019-12-25T12:12:12.123Z",
      "dateCompleted": "2019-12-26T13:13:13.456Z"
    }
  ],
  "meta": {
    "numberOfResults": 1,
    "totalPages": 1
  }
}

400 The request is missing a required parameter.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "<named> parameter is missing."
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


GET
/balance
Get the account balance

The current balance for the account of the entity is returned from this request. Currently there is only one credit unit ("currency") implemented in HOMEPESA-MC, but multiple units may be supported in the future.

Parameters
Name  Description
querying_entity_id *
string
(query)
The entity ID to which the account is linked

5e561916ca06e1c8596eee9e
Responses
Code  Description 200 OK

Media type

`application/json`

Controls Accept header.
Example Value


     Schema
{
  "data": [
    {
      "unit": "ocn-uk",
      "balance": -1.23
    }
  ]
}

401 
There was an issue with the authentication data for the request.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Could not authenticate you."
    }
  ]
}

403 
User does not have permission to access the resource.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Permission denied."
    }
  ]
}

429 The request limit for this resource has been reached for the current rate limit window.

Media type

`application/json`

Example Value


     Schema
{
  "errors": [
    {
      "message": "Rate limit exceeded."
    }
  ]
}

500 
An unknown internal error occurred.

Media type

`application/json`


Example Value


     Schema
{
  "errors": [
    {
      "message": "Internal server error triggered."
    }
  ]
}


Schemas
SignupRequiredFields
SignupAllFields
User
Entity
Category
Tag
TransferInitiated
TransferView
Balance
Error
Meta