# Usar la imagen oficial de Go como base
FROM golang:1.23.3-alpine

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar go mod y go sum para asegurar que las dependencias se instalen correctamente
COPY go.mod go.sum ./
RUN go mod tidy

# Copiar el resto del código fuente de la aplicación
COPY . .

# Establecer las variables de entorno necesarias para la conexión a MongoDB
ENV MONGO_URI=mongodb://52.205.123.191:27017/DeleteCustomerDB
ENV MONGO_DB_NAME=DeleteCustomerDB
ENV PORT=8084
ENV CREATE_CUSTOMER_SERVICE=http://localhost:8081
ENV READ_CUSTOMER_SERVICE=http://localhost:8082
ENV UPDATE_CUSTOMER_SERVICE=http://localhost:8083

# Exponer el puerto de la aplicación
EXPOSE 8084

# Copiar el archivo .env al contenedor
COPY .env .env

# Compilar la aplicación
RUN go build -o main .

# Comando para ejecutar la aplicación
CMD ["./main"]
