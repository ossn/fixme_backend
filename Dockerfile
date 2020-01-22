FROM gobuffalo/buffalo:v0.14.9 as builder

RUN apt update && apt install -y git ssh ca-certificates
RUN mkdir -p /go/src/github.com/ossn/fixme_backend
WORKDIR /go/src/github.com/ossn/fixme_backend
RUN mkdir ~/.ssh && ssh-keyscan -t rsa github.com >~/.ssh/known_hosts

COPY . .
RUN GO111MODULE=on buffalo build --environment=production --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash ca-certificates curl

WORKDIR /bin/

COPY --from=builder /bin/app .

# Uncomment to run the binary in "production" mode:
ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

CMD /bin/app migrate; /bin/app
