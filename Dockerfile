# build stage
FROM golang:1.22 as builder

WORKDIR /twilio-v2v/src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /twilio-v2v/dist

# TODO: add a test stage to the build

# production image
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=builder /twilio-v2v/dist /twilio-v2v

EXPOSE ${TWILIO_V2V_PORT}

USER nonroot:nonroot

ENTRYPOINT [ "/twilio-v2v" ]
