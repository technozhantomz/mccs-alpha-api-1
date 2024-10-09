package es

import (
	"context"
	"log"

	"github.com/olivere/elastic/v7"
)

func checkIndexes(client *elastic.Client) {
	for _, indexName := range indexes {
		checkIndex(client, indexName)
	}
}

func checkIndex(client *elastic.Client, index string) {
	ctx := context.Background()

	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		panic(err)
	}

	if exists {
		return
	}

	createIndex, err := client.CreateIndex(index).BodyString(indexMappings[index]).Do(ctx)
	if err != nil {
		panic(err)
	}
	if !createIndex.Acknowledged {
		panic("CreateIndex " + index + " was not acknowledged.")
	} else {
		log.Println("Successfully created " + index + " index")
	}
}

var indexes = []string{"entities", "users", "tags", "journals", "user_actions"}

// Notes:
// 1. Using nested fields for arrays of objects.
var indexMappings = map[string]string{
	"entities": `
	{
		"settings": {
			"analysis": {
				"analyzer": {
					"tag_analyzer": {
						"type": "custom",
						"tokenizer": "whitespace",
						"filter": [
							"lowercase",
							"asciifolding"
						]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"entityID": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"name": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"email": {
					"type": "keyword"
				},
				"status": {
					"type": "keyword"
				},
				"offers": {
					"type" : "nested",
					"properties": {
						"createdAt": {
							"type": "date"
						},
						"name": {
							"type": "text",
							"analyzer": "tag_analyzer",
							"fields": {
								"keyword": {
									"type": "keyword",
									"ignore_above": 256
								}
							}
						}
					}
				},
				"wants": {
					"type" : "nested",
					"properties": {
						"createdAt": {
							"type": "date"
						},
						"name": {
							"type": "text",
							"analyzer": "tag_analyzer",
							"fields": {
								"keyword": {
									"type": "keyword",
									"ignore_above": 256
								}
							}
						}
					}
				},
				"categories": {
					"type": "text",
					"analyzer": "tag_analyzer",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"city": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"region": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"country": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"accountNumber": {
					"type": "keyword"
				},
				"balance": {
					"type" : "float"
				},
				"maxNegBal": {
					"type" : "float"
				},
				"maxPosBal": {
					"type" : "float"
				}
			}
		}
	}`,
	"users": `
	{
		"mappings": {
			"properties": {
				"email": {
					"type": "keyword"
				},
				"firstName": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"lastName": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"userID": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				}
			}
		}
	}`,
	"tags": `
	{
		"settings": {
			"analysis": {
				"analyzer": {
					"tag_analyzer": {
						"type": "custom",
						"tokenizer": "whitespace",
						"filter": [
							"lowercase",
							"asciifolding"
						]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"name": {
					"type": "text",
					"analyzer": "tag_analyzer",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"offerAddedAt": {
					"type": "date"
				},
				"tagID": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"wantAddedAt": {
					"type": "date"
				}
			}
		}
	}`,
	"journals": `
	{
		"mappings": {
			"properties": {
				"transferID": {
					"type": "keyword"
				},
				"fromAccountNumber": {
					"type": "keyword"
				},
				"toAccountNumber": {
					"type": "keyword"
				},
				"status": {
					"type": "keyword"
				},
				"createdAt": {
					"type": "date"
				}
			}
		}
	}`,
	"user_actions": `
	{
		"mappings": {
			"properties": {
				"userID": {
					"type": "keyword"
				},
				"email": {
					"type": "keyword"
				},
				"action": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword",
							"ignore_above": 256
						}
					}
				},
				"detail": {
					"type": "text",
					"fields": {
						"keyword": {
							"type": "keyword"
						}
					}
				},
				"category": {
					"type": "keyword"
				},
				"createdAt": {
					"type": "date"
				}
			}
		}
	}`,
}
