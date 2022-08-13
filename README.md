![Build](https://img.shields.io/github/workflow/status/Jess-v/legbot-go/Build%20and%20Publish%20Image) 
![Latest Version](https://img.shields.io/docker/v/jessvv/legbot-go?sort=semver)
![Image Pulls](https://img.shields.io/docker/pulls/jessvv/legbot-go)
![Docker Stars](https://img.shields.io/docker/stars/jessvv/legbot-go)
![Image Size](https://img.shields.io/docker/image-size/jessvv/legbot-go)

# legbot-go

Do you take a medication that must be injected intramuscularly on a weekly basis into your leg? 

Do you find yourself struggling to remember which leg you're supposed to inject into when the day comes?

Perhaps legbot-go can help.

This project is an iteration on my original [legbot project](https://github.com/Jess-v/legbot), in an effort to learn Go. Overall it's an incredibly simple bot under the hood (as you can see by taking a peek at main.go) and is a near direct port of my original Python based bot.

## Environment Variables
| Variable            | Default | Description            | Required? |
| --------------------|---------|------------------------|-----------|
| `DISCORD_API_TOKEN` | `None`  | Discord Bot Auth Token | `True`    |

## Usage

First, if you do not have a Discord bot token, you will need that first. One guide on how to create one can be found [here](https://www.writebots.com/discord-bot-token/).

Once you have your API token, create a file named `.env` and insert the following, but with your API token after the equals sign:

```bash
DISCORD_API_TOKEN=
```

Next, decide where you want your folder that stores leg-related data to live. Make note of the path to this folder.

Finally, run the following:

```bash
docker run --rm --env-file .env -v <YOUR/LOCAL/PATH/TO/FOLDER>:/app/users/ --name="legbot" -d jessvv/legbot-go:latest
```
