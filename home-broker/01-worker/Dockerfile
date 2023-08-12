# Criando a imagem intermediária de build da aplicação
FROM golang:1.20.3-bullseye as build

WORKDIR /go/src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# Compilando a aplicação
RUN GOOS=linux CGO_ENABLED=1 go build -o trade cmd/trade/main.go

# Criação da imagem final
FROM golang:1.20.3 as app

# Definindo o Workdir
WORKDIR /app

# Copy files and directories from the build stage
COPY --from=build /go/src/trade .

# Comando de inicialização da aplicação
CMD ["./trade"]