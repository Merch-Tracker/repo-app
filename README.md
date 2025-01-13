# Repository-api

1. Fill in config.env
2. docker buildx build . -t <image-name>
3. docker run --name <container-name> --env-file config.env -p 9010:9010 -p 9050:9050 <image-name>
