package tests

//s4bmongo:collection users
type User struct {
	Email   string  `bson:"email"`
	Profile Profile `bson:"profile"`
	Address struct {
		City string
	} `bson:"address"`
	Fin         *Finance
	Permissions map[string]interface{}
	Ids         map[string]int
}

type Profile struct {
	FirstName string `bson:"first_name"`
	LastName  string `bson:"last_name"`
}

type Finance struct {
	Income int64
}
