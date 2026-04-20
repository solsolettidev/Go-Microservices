package data

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models {
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry

}

type LogEntry struct{
	ID string `bson: "_id,omitempty" json:"id,omitempty"`  // BSON is the type used for mongo, json is for the web
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	Data string `bson:"data,omitempty" json:"data,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// funcs that allows us to interact to the database

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs") // create a collection called logs in the database called logs, mongo creates it if it doesn't exist

	_, err := collection.InsertOne(context.TODO(), Logentry{ // insert the entry into the collection, context.TODO() is used to create a context that will be used to cancel the operation if it takes too long
		Name: entry.Name,
		Data: entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}) 
	if err != nil {
		log.Println("Error inserting log entry:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel :=  context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	collection := client.Database("logs").Collection("logs")
	
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}}) // sort by created_at in descending order

	cursor, err := collection.Find(context.TODO(), bson.D{{}, opts}) // find all documents in the collection
	if err != nil{
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx) // close the cursor

	var logs []*LogEntry // create a slice of log entries

	for cursor.Next(ctx){ // iterate through the cursor and decode each document
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil{
			log.Println("Error decoding log into slice:", err)
			return nil, err
		}
		logs = append(logs, &item) // append the decoded entry to the slice
	}

	return logs, nil
}