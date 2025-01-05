#!/bin/zsh

BASE_URL="http://localhost:5001/get?key="

while true; do
  RANDOM_KEY="key-$((RANDOM))"
  
  echo "Querying key: $RANDOM_KEY"
  RESPONSE=$(curl -s "$BASE_URL$RANDOM_KEY")
  
  SHARD=$(echo "$RESPONSE" | grep -o 'Shard : [0-9]' | awk '{print $3}')
  VALUE=$(echo "$RESPONSE" | grep -o 'Value : ".*"' | cut -d'"' -f2)
  

  if [[ "$SHARD" == "0" && -n "$VALUE" ]]; then
    echo "Success: Found Shard 1 with a non-empty value!"
    echo "Key: $RANDOM_KEY"
    echo "Response: $RESPONSE"
    break
  elif [[ "$SHARD" == "0" ]]; then
    echo "Shard 1 found, but value is empty. Continuing..."
  else
    echo "Shard is not 1 or value is empty. Continuing..."
  fi
done
