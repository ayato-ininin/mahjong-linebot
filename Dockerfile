FROM golang:1.19-alpine

WORKDIR /app

#COPY go.mod ./
#COPY go.sum ./
COPY . ./

#gcloud auth application-default login　の代替
#local実行の際は、key.jsonを配置して、下記コメントを外す。
#ENV GOOGLE_APPLICATION_CREDENTIALS secret-manager-access-key.json

RUN go mod download

RUN go build -o /mahjong-linebot .

EXPOSE 8080

CMD [ "/mahjong-linebot" ]
#CMD [ "./main" ]
