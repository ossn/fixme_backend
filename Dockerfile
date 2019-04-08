FROM gobuffalo/buffalo:v0.13.3 as builder

RUN apt update && apt install -y git ssh jq curl ca-certificates
RUN mkdir -p /go/src/github.com/ossn/fixme_backend
WORKDIR /go/src/github.com/ossn/fixme_backend
RUN mkdir ~/.ssh && \
  ssh-keyscan -t rsa github.com > ~/.ssh/known_hosts

# Install dep
RUN curl -fsSL -o /usr/local/bin/dep $(curl -s https://api.github.com/repos/golang/dep/releases/latest | jq -r ".assets[] | select(.name | test(\"dep-linux-amd64\")) |.browser_download_url") && chmod +x /usr/local/bin/dep

# Build app
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only

COPY . .
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

# Uncomment to run the binary in "production" mode:
ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

CMD /bin/app migrate; /bin/app
