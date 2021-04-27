#echo -n "this object will have only 1 instance" | openssl dgst -sha256 -binary | base64
#aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=

curl -v jojo:9001/objects/test4_1 -XPUT -d"this object will have only 1 instance" -H "Digest: SHA-256=aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY="

curl localhost:9001/locate/test4_1

echo
curl localhost:9002/objects/test
echo