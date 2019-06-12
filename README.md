![es](https://raw.githubusercontent.com/txn2/es/master/mast.jpg)

**es** is minimal, opinionated Elasticsearch utility lib for TXN2 services.


## Demo

### Setup

Bring up a single node Elasticsearch server and Kibana instance.
```bash
docker-compose up
```

- Elasticsearch is available at http://localhost:9200
  - see: http://localhost:9200/_cluster/health
- Kibana is available at http://localhost:5601

### Run Example

```bash
go run ./example/es.go
```