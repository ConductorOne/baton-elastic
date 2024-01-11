FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-elastic"]
COPY baton-elastic /