How to use

```
DB=$(docker run -t -i  gosocks-proxy bootstrap_db)
docker run -d --privileged -p 8080:8080/tcp -p  8000:8000/tcp --volumes-from $DB  --restart=always gosocks-proxy 
```
