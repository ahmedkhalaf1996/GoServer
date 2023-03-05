## -------------------------|  live Streming Docs  |-------------------------------- #

## create session
## Post Request
localhost:9000/session

Body {
     "title":"ahmed",
     "host":"ahmed"
   # example "title":"ahmed",
   # example "host":"ahmed"
}
>> will return 
{
  "socket": stirng
 ## Example     "socket": "c17eaac2e7714f5ec0c485b27629a944bfa4523b"
}

## Get Session GetSession - Checks if session exists.
## Get Request
localhost:9000/connect?url=SessionUrl

>> return 200 stuatus Response if exists
>> Or  "error": "Socket connection not found." as json data

## !! ConnectSession - Given a host returns the session object.
## Post Reques
localhost:9000/connect/SessionUrl
body {
    "host" : stirng
  # example  "host":"ahmed"
}


# ------------- #
## Ws coonection

localhost:9000/ws/:SocketString