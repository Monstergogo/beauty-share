# beauty-share
a toy application for test, used minio to store images and mongo db as store backend.

discovery: [consul](https://developer.hashicorp.com/consul) \
api gateway: [apisix](https://github.com/apache/apisix)

interfaces:

| name           | type | route          | remark                              | 
|----------------|------|----------------|-------------------------------------|
| OssUpload      | http | /v1/oss/upload | upload images to minio              |
| Ping           | http | /v1/ping       | used to check service health status |
| AddShare       | grpc | --             | add share info to mongo             |
| GetShareByPage | grpc | --             | get share info by page              |


## how to run
1. ```cd example && docker compose up -d```
2. open http://127.0.0.1:9000 in browser, create a minio bucket and access key, default **username**: root, **password**: root1234
3. open http://127.0.0.1:8500/ui in browser(consul ui), select Key/Value, creat a folder named beauty-share and config some key info:
   - key: **mongo_uri**, value: **mongodb://root:root123@localhost:27017/share**
   - key: **minio**, value:
        ```
        {
        "id": "SD6p9BBCeaia4fxxx",
        "secret": "2T84qwUEHmFSqkQOanVMAxlivaxxxx",
        "endpoint": "localhost:9000",
        "bucket": "photos"
        }
        ```
      id and secret is minio access info created in step 2, images save in bucket photos

4. run server in cmd/server/main.go local, can see server register successfully in consul ui(http://127.0.0.1:8500/ui/dc1/services)
5. create route in apisix dashboard: http://127.0.0.1:9005/ and test server