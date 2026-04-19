package main
import(
	"log"

	"go.mongodb.org/mongo-driver/mongo" // read documentation for connection instructions
)

const (
	webPort = "80"
	rpcPont = "5001"
	mongoURL = "mongodb://,pmgp:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	} ()
}

func connectToMongo()(*mongo.Client, error){
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect to mongo
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return c, nil
}