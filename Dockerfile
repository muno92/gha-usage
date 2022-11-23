FROM golang:1.19-bullseye as builder
WORKDIR /app

COPY . .
RUN make

FROM gcr.io/distroless/base-debian11
COPY --from=builder /app/github_actions_usage_calculator ./github_actions_usage_calculator
CMD ["./github_actions_usage_calculator"]
