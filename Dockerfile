FROM oven/bun:1 AS bun-build

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY tailwind.config.js ./
COPY app/view/style.css ./app/view/style.css
COPY app/view/tmpl/*.html ./app/view/tmpl/
RUN bun run tailwindcss -m -i ./app/view/style.css -o ./app/view/dist/style.css

FROM golang:1.24 AS go-build

WORKDIR /app

COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

COPY --from=bun-build /app/app/view/dist/style.css ./app/view/dist/style.css

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o main ./app/cmd/app

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=go-build /app/main /

CMD ["/main"]
