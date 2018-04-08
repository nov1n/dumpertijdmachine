go build -o out/server server/main.go
docker build . -t "gcr.io/dumpertijdmachine/server:0.7"
gcloud docker -- push gcr.io/dumpertijdmachine/server:0.7
