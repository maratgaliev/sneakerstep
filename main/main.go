package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/graphql-go/graphql"
)

func _check(err error) {
	if err != nil {
		panic(err)
	}
}

type Sneaker struct {
	ID       int
	Title    string
	Price    string
	Date     string
	Image    string
	Provider string
}

var SneakerList []Sneaker

func parseUrl(url string) []Sneaker {
	fmt.Println("request: " + url)
	doc, err := goquery.NewDocument(url)
	_check(err)

	doc.Find(".release-group__container").Each(func(i int, item *goquery.Selection) {
		date1 := item.Find(".clg-releases__date__day").Text()
		date2 := item.Find(".clg-releases__date__month").Text()
		date := date1 + "/" + date2 + "/2019"
		item.Find(".sneaker-release-item").Each(func(i int, sneaker_block *goquery.Selection) {
			id := i + 1
			title := sneaker_block.Find(".sneaker-release__title").Text()
			price := strings.TrimSpace(sneaker_block.Find(".sneaker-release__option--price").Text())
			image, _ := sneaker_block.Find(".sneaker-release__img-16x9 a img").Attr("src")
			sneaker := Sneaker{id, title, price, date, image, "SOLECOLLECTOR"}
			SneakerList = append(SneakerList, sneaker)
		})
	})

	return SneakerList
}

func init() {
	parseUrl("https://solecollector.com/sneaker-release-dates/all-release-dates")
}

// define custom GraphQL ObjectType `sneakerType` for our Golang struct `Sneaker`
// Note that
// - the fields in our sneakerType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
var sneakerType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Sneaker",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.String,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"image": &graphql.Field{
			Type: graphql.String,
		},
		"provider": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// root query
// we just define a trivial example here, since root query is required.
// Test with curl
// curl -g 'http://localhost:8080/graphql?query={lastSneaker{id,text,done}}'
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{

		/*
		   curl -g 'http://localhost:8080/graphql?query={sneaker(id:"b"){id,title}}'
		*/
		"sneaker": &graphql.Field{
			Type:        sneakerType,
			Description: "Get single sneaker",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				idQuery, isOK := params.Args["id"].(int)
				if isOK {
					// Search for el with id
					for _, sneaker := range SneakerList {
						if sneaker.ID == idQuery {
							return sneaker, nil
						}
					}
				}

				return Sneaker{}, nil
			},
		},

		/*
		   curl -g 'http://localhost:8080/graphql?query={sneakerList{id,text,done}}'
		*/
		"sneakerList": &graphql.Field{
			Type:        graphql.NewList(sneakerType),
			Description: "List of sneakers",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return SneakerList, nil
			},
		},
	},
})

// define schema, with our rootQuery and rootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
	//Mutation: rootMutation,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	// Display some basic instructions
	fmt.Println("Now server is running on port 8080")
	fmt.Println("Get single sneaker: curl -g 'http://localhost:8080/graphql?query={sneaker(id:\"b\"){id,text,done}}'")
	fmt.Println("Load sneaker list: curl -g 'http://localhost:8080/graphql?query={sneakerList{id,text,done}}'")
	fmt.Println("Access the web app via browser at 'http://localhost:8080'")

	http.ListenAndServe(":8080", nil)
}
