FROM golang as build-stage
WORKDIR /app

COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY main.go /app

#RUN cd /app && GOOS=linux go build -a -installsuffix cgo -o orders-gen
RUN cd /app && GOOS=linux go build -o orders-gen

RUN ls -lrt
CMD ["/app/orders-gen"]

#FROM alpine
#COPY --from=build-stage /app/orders-gen /
#RUN ls -lrt
#CMD ["/orders-gen"]