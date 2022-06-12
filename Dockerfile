FROM golang:alpine

RUN mkdir /app
WORKDIR /app

COPY . .
COPY .env .

RUN go get -d -v ./...
RUN go install -v ./...

RUN cd cmd && go build -o /port-domain-service && cd ../

CMD [ "/port-domain-service" ]