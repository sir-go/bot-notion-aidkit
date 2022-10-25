### Telegram bot for search the medicine in my home first aid kit

> https://t.me/mserg_apteka_bot

___
### How it works

The bot listens to users' requests and searches records via Notion API in the database. 
If something is found, a user gets the answer including a photo of the medicine and tags, and location info.

The data updates from the database periodically or by the command `/update`.

The database easily administrates at the Notion service.
___
### Configuration file

Option `-c` sets the configuration file path (default `./conf.toml`).

```toml
[tg_api]
    token = '...'                                   # Tg bot token, BotFather (https://t.me/BotFather) helps

[notion_api]
    token = '...'                                   # notion account API token
    timeout = '10s'                                 # request timeout
    version = '2021-08-16'                          # API version
    db_id = '...'                                   # notion database ID
    search_url = 'https://api.notion.com/v1/search' # seach API URL
    update_interval = '45m'                         # autoupdate data time interval

```
___
### Build & run
```bash
go mod download
go build -o bot ./cmd/bot
```
___
### Screencast

#### DB in the Notion table
![](notion.png)

#### Bot interaction
![](sc.gif)
