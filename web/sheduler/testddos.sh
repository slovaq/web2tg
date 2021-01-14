for (( i=1; i <= 1000; i++ ))
do
echo "number is $i"
    url="http://127.0.0.1:3001/add?message=test$i&date=2021-01-14&time=13:47:10"
curl $url
done