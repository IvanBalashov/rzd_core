# RZD_CORE
This is core for rzd_bot in Telegram. Use this project for example Clean Arch too.

### Main thinks
I'm hope what i can end this project. Now this app have only web-server for control business logic. Main feature is use 
one instance for serving one part of data. Use kubernetes in production for run more than one instance in main app.

## Configuring
Now app configuring with environment variables: 
```text
HTTP_HOST - host where app running. Must be set.
HTTP_PORT - port where app running. Must be set.
POSTGRES_URL - url for connection to postgres. Must be set.
RABBITMQ_URL - url for connection to rabbitMQ. Must be set.
APP_NAME - app name was placed in logs. Must be set.
MONGODB_URL - url for connection to MongoDB. Must be set.
```

Now app can read environment variables from `.env` file. For fast start use follow command.
```text
cp .env_example .env
```

### Main architecture
```text
+----------+
| Rzd API  |
+----------+
     |
+----------+           +----------+ 
| Rzd Core | --------- | MongoDB  |
+----------+           +----------+
 | | | | |  \
+----------+ \         +----------+
| RabbitMQ |  \________| Memcache |
+----------+           +----------+
 | | | | | 
+----------+
|  TG BOT  |
+----------+
```

### Data flow

```text
1) Bot Send in RabbitMQ message like
{MessageID:[int], Event:[string], User: UserObj, TrainArgs: ArgsForSearch}.
2) First free node get data from Queue.
    - Request to RZDApi  for getting trains on route.
    - Send answer to bot like {[]array_with_trains}.
    - Waiting answer from bot about train to be saved.
    - Got Request from bot about what train need to be saved.
    - Save train and user info in MongoDB.
3) 
```