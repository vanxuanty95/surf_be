FROM golang:1.16-alpine

# We want to populate the module cache based on the go.{mod,sum} files.
WORKDIR /home/
COPY surf_be.bin .
COPY configurations/ configurations/


# This container exposes port 8080 to the outside world
EXPOSE 8082

# Run the binary program produced by `go install`
CMD ["./surf_be.bin -env dev"]