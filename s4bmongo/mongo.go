package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	fmt.Println("s4bmongo " + primitive.NewObjectID().Hex())
}
