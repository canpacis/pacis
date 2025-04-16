FROM golang:latest

# Set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Install dependencies and build
RUN go install ./pages/pcpg \
    && mkdir -p /app/bin \
    && cp $(go env GOPATH)/bin/pcpg /app/bin/ \
    && /app/bin/pcpg i \
    && cp /usr/local/bin/pcpg_tw /app/bin/ \
    && /app/bin/pcpg c www \
    && go build -o pacis-www ./www

ENV PATH="/app/bin:${PATH}"

EXPOSE 8080

CMD ["./pacis-www"]