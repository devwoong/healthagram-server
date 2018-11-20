curl -i -X POST -F uplTheFile=@edit.png -F "name=uplTheFile" localhost:8000/upload

curl -i -X POST -H "Content-Type: application/json" -d '{"id":"UDID","os":"1","phone":"PhoneNumber"}' localhost:8000/json
