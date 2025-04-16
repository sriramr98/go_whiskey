get:
	curl --request GET -sL --url 'localhost:8080/hello'

post:
	curl --request POST -sL --url 'localhost:8080/hello' -H 'Content-Type: application/json' -d '{"key1":"value1", "key2":"value2"}'