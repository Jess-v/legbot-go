![Build](https://img.shields.io/github/actions/workflow/status/Jess-v/legbot-go/build_publish.yaml) 
![Latest Version](https://img.shields.io/docker/v/jessvv/legbot-go?sort=semver)
![Image Pulls](https://img.shields.io/docker/pulls/jessvv/legbot-go)
![Docker Stars](https://img.shields.io/docker/stars/jessvv/legbot-go)
![Image Size](https://img.shields.io/docker/image-size/jessvv/legbot-go)

# legbot-go

Do you take a medication that must be injected intramuscularly on a weekly basis into your leg? 

Do you find yourself struggling to remember which leg you're supposed to inject into when the day comes?

Perhaps legbot-go can help.

This project is an iteration on my original [legbot project](https://github.com/Jess-v/legbot), in an effort to learn Go.

## Environment Variables
| Variable            | Default     | Description                                     |
| --------------------|-------------|-------------------------------------------------|
| `BOT_TOKEN`         | `None`      | Discord Bot Auth Token                          | 
| `LOG_LEVEL`         | `info`      | Sets the bot's log level                        | 
| `DATABASE_HOST`     | `localhost` | Hostname of the connnected Postgres database    |
| `DATABASE_PORT`     | `5432`      | Connection port of Postgres database            |
| `DATABASSE_NAME`    | `legbot`    | Name of the postgres database that will be used |
| `DATABASE_USER`     | `postgres`  | Username of the postgres user that will be used |
| `DATABASE_PASSWORD` | `postgres`  | Password of the postgres user that will be used |

## Usage

The easiest method of using this bot is to get a Discord API token, and inserting this into the proper variable within the provided `docker-compose.yml`. From there, a simple 

```sh
docker-compose up -d
```

will get the bot up and running. The commands for usage of the bot are `/set`, `/where`, and `/praise`.