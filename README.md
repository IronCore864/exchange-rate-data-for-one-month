# Tiexin's Solutoin to Exchange Rate

## Language

Here I used golang. Mainly because it's a compile language, and this problem requires docker, so the image size can be much smaller.

Golang is easier to write than C/Rust and still can maintain a very small image size (in this example, 14MB).

If using script language like python, or languages that requires JVM to run like scala/Java, the docker image size would be much bigger, so deployment time is way too long compared to run time.

## DB

Here I usd redis.

In fact any other kv store, doc db, or even rdb could do it, but since this is a very simple task, rdb probably would be an overkill.

One important reason I chose redis over others is the fact that redis set supports expiration, which fits perfectly the 30 day store requirement of this task.

## Local Deploy/Run/Test without Docker/K8s

### Unit Test

```
git clone git@github.com:IronCore864/exchange-rate-data-for-one-month.git
cd exchange-rate-data-for-one-month
go test ./...
```

I only implemented a little UT for the sake of demostration, no test coverage is guaranteed.

### Local Run

On Mac OS, For a faster build/run/test without having to build docker images:

```
brew install redis
# keeps running
redis-server

# open another tab
git clone git@github.com:IronCore864/exchange-rate-data-for-one-month.git
cd exchange-rate-data-for-one-month
go get ./...
go build
./exchange-rate-data-for-one-month
```

## Build Docker Image

```
git clone git@github.com:IronCore864/exchange-rate-data-for-one-month.git
cd exchange-rate-data-for-one-month
docker build -t ironcore864/exchange-rate-data-for-one-month .
docker push ironcore864/exchange-rate-data-for-one-month
```

## Test Docker Image Locally without K8s

```
# keeps running
redis-server
# in another tab
docker run --rm -e RedisHost=host.docker.internal -e RedisPort=6379 -e RedisPassword="" ironcore864/exchange-rate-data-for-one-month
```

This will let the image access the redis running on localhost.

## Deploy in K8s (be it local or production)

### Dependency

- a k8s cluster running
- helm already installed locally
- helm init done in k8s

### Redis

```
helm install stable/redis --name redis
```

### Run One Time (e.g., local k8s test)

```
git clone git@github.com:IronCore864/exchange-rate-data-for-one-month.git
cd exchange-rate-data-for-one-month/k8s
kubectl apply -f job.yaml
```

### Update on a Daily Basis

```
git clone git@github.com:IronCore864/exchange-rate-data-for-one-month.git
cd exchange-rate-data-for-one-month/k8s
kubectl apply -f cron.yaml
```

### Connect to Redis in K8s

```
export REDIS_PASSWORD=$(kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 --decode)
kubectl run --namespace default redis-client --rm --tty -i --restart='Never' --env REDIS_PASSWORD=$REDIS_PASSWORD --image docker.io/bitnami/redis:4.0.14 -- bash
# in the pod
redis-cli -h redis-master -a $REDIS_PASSWORD
```

## Notes

- at the moment (13:32, 05.04.2019), https://api.exchangeratesapi.io/latest does not work for me (maybe because I accessed so many times, maybe it's down), so I changed the url to https://api.exchangeratesapi.io/2019-04-05.
- the task asks to access https://api.exchangeratesapi.io/latest only; in fact there is another approach which is by changing the URL so access 30 days' data in one run and store it. Didn't use this method becuase it seems this is not what the task wants.
- cron is configured in the midnight so that it makes a minimum impact on the cluster performance. I didn't test the cron.yaml becuase I don't want to wait till tomorrow.
- I merged both config.json (used when testing locally without docker and no need to set env) and env vars (for docker/k8s usage)

