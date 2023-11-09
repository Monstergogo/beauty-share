# beauty-share
a toy application for test, used minio to store images and mongo db as store backend.

discovery: [nacos-go](https://github.com/nacos-group/nacos-sdk-go) \
gate-way: [apisix](https://github.com/apache/apisix)

interfaces:

| name           | type | route          | remark                              | 
|----------------|------|----------------|-------------------------------------|
| OssUpload      | http | /v1/oss/upload | upload images to minio              |
| Ping           | http | /v1/ping       | used to check service health status |
| AddShare       | grpc | --             | add share info to mongo             |
| GetShareByPage | grpc | --             | get share info by page              |


## how to run
1. ```cd example && docker compose up -d```
2. open 127.0.0.1:9000 in browser, create a minio bucket and create access key
3. open 127.0.0.1:8848 in browser, in nacos console, config info like this:

    | Data_id           | Group         | 配置内容                                         | 内容格式 |
    |-------------------|---------------|----------------------------------------------|------|
    | mongo_uri         | DEFAULT_GROUP | mongodb://root:root123@localhost:27017/share | TEXT |
    | minio_endpoint    | DEFAULT_GROUP | localhost:9000                               | TEXT |
    | minio_id          | DEFAULT_GROUP | minio access key                             | TEXT |
    | minio_secret      | DEFAULT_GROUP | minio access secret key                      | TEXT |
    | share_bucket_name | DEFAULT_GROUP | minio bucket name                            | TEXT |

4. run server in cmd/server/main.go local