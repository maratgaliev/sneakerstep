# sneakerstep

GraphQL API for the latest sneakers releases based on Go

## Run application
```
go run main.go
```

### Now server is running on port 8080

## Load sneaker list: 
```
curl -g 'http://localhost:8080/graphql?query={sneakerList{id,text,done}}'
```

## Get single sneaker: 
```
curl -g 'http://localhost:8080/graphql?query={sneaker(id:"b"){id,text,done}}'
```

