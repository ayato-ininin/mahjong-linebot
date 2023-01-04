FROM golang:1.19-alpine

WORKDIR /app

#COPY go.mod ./
#COPY go.sum ./
COPY . ./
RUN go mod download

RUN go build -o /mahjong-linebot .

EXPOSE 8080

CMD [ "/mahjong-linebot" ]
#CMD [ "./main" ]
