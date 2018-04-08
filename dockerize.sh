go build server/main.go
docker build . -t "gcr.io/dumpertijdmachine/server:0.3"
gcloud docker -- push gcr.io/dumpertijdmachine/server:0.3
