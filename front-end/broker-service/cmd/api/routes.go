package main
import(
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)
func (app *Config) routes() http.Handler { // func that returns an http.handler

	mux := chi.NewRouter() // variable that is the comun name when using routing in go, its the abreviation of Multiplexor (it takes
	// multiple entry signals and redirect them where their belongs
	

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{ // cors is used to allow other origins to connect to our api
		AllowedOrigins: []string{"http://*","http://*"},//the allowed ones r everyone for dev purposes
	    AllowedMethods: []string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders: []string{"Accept","Authorization","Content-Type","X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	mux.Use(middleware.Heartbeat("/ping")) 
	// this is a healthcheck endpoint, it is used to check if the server is alive, if it is it will return a 200 status code (enterying the url /ping)

	mux.Post("/", app.Broker)

	mux.Post("/handle", app.HandleSubmission)

	return mux
}

/*
Aquí utilizaste el paquete Chi, que es un enrutador extremadamente popular, ligero y estándar en la industria de Go.

func (app *Config) routes() http.Handler: Al poner (app *Config) antes del nombre de la función, estás "pegando" esta función al struct que creaste en main.go. Por eso allá podías llamarla con app.routes(). Devuelve un objeto que entiende las peticiones HTTP.
mux := chi.NewRouter(): Este es el objeto principal al que le asignas rutas. "Mux" es una abreviatura universal para Multiplexor (que toma múltiples señales de entrada y las redirecciona donde pertenecen).
El bloque de Middlewares (mux.Use(...)): Los middlewares son como guardias de seguridad en la puerta de la discoteca. Toda petición HTTP pasa por ellos ANTES de ejecutar el código de tu programa.

CORS (Cross-Origin Resource Sharing): ¡La tortura número uno de los desarrolladores Web! Este bloque cors.Handler(...) configura quién puede hablarle a tu API.
¿Para qué sirve?: Cuando un navegador (en un Frontend misitio.com) hace una petición en segundo plano a un Backend distinto (mi-api.com), el navegador bloquea la conexión por seguridad. Para permitirlo, tu API tiene que enviar Cabeceras (Headers) diciéndole al navegador "Ey, sí conozco misitio.com, dejalo entrar".
AllowedOrigins: Quién tiene acceso. (Tienes comodines * para local/dev).
AllowedMethods: Qué tipo de peticiones tienen permitido hacer (GET para pedir cosas, POST para mandar... etc).
MaxAge: 300: El navegador primero hace una preguntita llamada Pre-flight de tipo OPTIONS para ver si tiene permisos. Para no saturar a tu servidor recibiendo esta pregunta doble todo el tiempo, el MaxAge le dice al navegador "Guarda esta respuesta de permiso por 300 segundos, no me sigas preguntando".
Heartbeat ("/ping"): Le dices al router que si alguien entra a http://localhost/ping, simplemente responda un "OK" vacío pero exitoso por debajo. Es valioso para que en el futuro un orquestador (como Docker Compose o Kubernetes) pregunte si tu API sigue viva (Healthcheck).
Las rutas explícitas (mux.Post("/", app.Broker)): Básicamente le dices: "Cuando un Frontend me haga una petición tipo POST a la URL raíz "/", mándamelo con el empleado encargado llamado app.Broker".
*/

/*Una función con un "Receiver" (Pegada)
go
// Ese (app *Config) es el "pegamento".
func (app *Config) routes() http.Handler {
   // La función ahora le pertenece a la caja.
   // Aquí adentro tienes acceso a cualquier cosa que esté dentro de 'app' mágicamente.
}
Al poner (app *Config) justo antes del nombre de la función routes(), le estás diciendo a Go: "A partir de hoy, la función routes le pertenece exclusivamente a la caja Config. Y cuando estemos adentro de routes, llamaremos a la caja con el pronombre app".*/

/*
En Go, para poder "pegarle" un método a un struct (es decir, crear un Receiver func (app *Config) ...), el Struct y la Función deben vivir exactamente en el mismo package (paquete)
Que vivan en un mismo package  para Go es como si fuera un solo archivo gigante y unido..
*/

/*
Lo que NO puedes hacer (Las limitaciones)
1. No puedes añadir métodos a structs de otras personas (o paquetes distintos) Imagina que usas una librería de otra persona o un paquete estándar de Go. Por ejemplo, el paquete de tiempo constante en Go usa un struct llamado Time que viene del paquete time (time.Time).

Si tú en tu main.go intentas hacer esto:

go
package main
import "time"
// ERROR: Esto no compilará
func (t time.Time) MiNuevoMetodo() {
    // ...
}
Go te va a dar un error rotundo. Te dirá: "No puedes definirle métodos a tipos que no fueron creados en este mismo paquete".

2. No puedes hacerlo con tipos "primitivos" directamente No puedes pegarle un método a variables básicas como string, int, o bool así nomás.

go
// ERROR: No puedes hacer esto
func (texto string) Gritar() { ... }
Truco: Si alguna vez necesitas hacer eso, existe un atajo que es crear un "Alias" o tipo personalizado en tu propio paquete:

go
type MiTexto string // Creas tú un tipo basado en string, ¡ahora te pertenece!
func (t MiTexto) Gritar() { ... } // Listo, esto sí funciona.
*/