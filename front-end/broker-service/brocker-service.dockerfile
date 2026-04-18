# base fo image
FROM golang:alpine as builder  
# esta es la imagen base de golang, la que se usa para compilar el codigo, "builder" es el nombre que le damos a esta etapa

RUN mkdir /app
# crea una carpeta llamada app

COPY . /app
# copia todo el codigo de la carpeta actual a la carpeta app

WORKDIR /app
# entra a la carpeta app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api
# compila el codigo, "CGO_ENABLED=0" es para que no use cgo, "-o brokerApp" es para que el ejecutable se llame brokerApp, "./cmd/api" es la ruta del codigo
# CGO es un comando que se usa para compilar codigo en go
RUN chmod +x /app/brokerApp
# le da permisos de ejecucion al archivo brokerApp

#build a tiny docker image
FROM alpine:latest
# esta es la imagen base de alpine, la que se usa para crear la imagen final, "alpine" es una distribucion de linux muy ligera

RUN mkdir /app
# crea una carpeta llamada app

COPY --from=builder /app/brokerApp /app
# copia el archivo brokerApp de la carpeta builder a la carpeta app

CMD ["/app/brokerApp"]
# este es el comando que se ejecuta cuando se inicia el contenedor