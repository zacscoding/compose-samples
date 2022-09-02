echo "Waiting for DynamoDB (1s)"
sleep 1s
host=$HOST

for f in /dynamodb/tables/*.json
do
    echo "aws dynamodb create-table --cli-input-json file://$f --endpoint-url http://$host --region ap-northeast-2"
    aws dynamodb create-table --cli-input-json file://$f --endpoint-url http://$host --region ap-northeast-2
done