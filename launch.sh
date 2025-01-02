set -e

trap 'killall KeyForge' SIGINT

cd $(dirname "$0")

killall KeyForge || true
sleep 1

go install -v

KeyForge -db-location=shards/shard1.db -http-address=127.0.0.1:5000 -config-file=sharding.toml -shard=shard1 &
KeyForge -db-location=shards/shard2.db -http-address=127.0.0.1:5001 -config-file=sharding.toml -shard=shard2 &
KeyForge -db-location=shards/shard3.db -http-address=127.0.0.1:5002  -config-file=sharding.toml -shard=shard3 &

wait