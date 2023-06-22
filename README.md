# uploader

```
sudo docker build -t olelarsen/uploader --build-arg NODE_ENV=production --build-arg APP_NAME=uploader --build-arg PORT=1234 --build-arg SESSION_SECRET=12345abscdefg --build-arg X_TOKEN=abcdefg12345 --build-arg USE_HASH=true --build-arg DB_SQL_HOST=localhost --build-arg DB_SQL_PORT=5432 --build-arg DB_SQL_USERNAME=postgres --build-arg DB_SQL_PASSWORD=postgres --build-arg DB_SQL_DATABASE=files --build-arg USE_DB=true .
