package main
import(
	"fmt"
	"log"
	"net/http"
)
const webPort = "80"

type Config struct {}
func main() {
// i want a root that responses to a JSON app, do something with it and respond yes
// i get the chi pack for rooter managment, the middleware and the cors one for cors protection btw front and back.
	app := Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	//define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),  // puedo llamarla directamente como app.routes() porque la func routes() tiene como receiver a app *Config

	}

	//start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
/*Este archivo es donde arranca la aplicación de Go.

const webPort = "80": Definis el puerto donde va a escuchar. Es el estandar para tráfico HTTP no cifrado, aunque para local a veces se usan puertos como :8080.
type Config struct {}: Esto es vital. Es un contenedor de estado. A medida que tu proyecto crezca, a este struct le vas a agregar cosas (por ejemplo: la conexión a la base de datos de logs, colas de mensajería como RabbitMQ, o clientes HTTP). Al inicializar esto como app := Config{}, y luego enviar esta app a todos lados, te aseguras de que el resto del código tenga acceso fácil a esos recursos. (Se le llama Inyección de Dependencias mediante Receiver Pattern).
srv := &http.Server{...}: Levanta la infraestructura HTTP nativa incorporada de Go (una de las cosas más rápidas y optimizadas que tiene este lenguaje).
Addr: es la dirección y el puerto.
Handler: app.routes(): Aquí es donde todo se conecta. Le estás diciendo al servidor genérico de Go: "Oye, todo el tráfico entrante mandáselo a mi propio sistema de enrutamiento que yo definiré en routes()".*/

