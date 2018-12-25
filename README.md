# RZD_CORE
This is core for rzd_bot in Telegram. Use this project for example Clean Arch.

### Main thinks
I'm hope what i can end this project. Now this app have only web-server for control business logic. Main feature is use 
one instance for serving one part of data. Use kubernetes in production for run more than one instance in main app.

##Configuring
Now app configuring with environment variables: 
```text
HTTP_HOST - host where app running. Must be set.
HTTP_PORT - port where app running. Must be set.
POSTGRES_URL - url for connection to postgres. Must be set.
```
