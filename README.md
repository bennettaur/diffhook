## Setting Up Slack
* Register a slack app in your org here: https://api.slack.com/apps?new_app=1
  * You can customize the name, if you like, or keep it `diffhook`
  * Choose the slack workspace you want diffhook to post in
* Under OAuth & Permissions add the following scopes:
  * `chat:write`, `chat:write.public`, `channels:read`, `groups:read`
  * A Future version will provide a CLI for configuring this (https://api.slack.com/authentication/oauth-v2)
* Install the Slack App into your workspace (you should be prompted for this)
  * Select any channel
* Add the bot to any private channels you want it to post to
  * Just send `@diffhook` in the channel 
* Use the generated Bot OAuth Token to set the `SLACK_TOKEN` environment variable

## Todo
 - [ ] Documentation for how to setup and use
 - [ ] Add better logging [zerolog?](https://github.com/rs/zerolog)
 - [ ] Split out web server & [mongo](https://github.com/Kamva/mgm) integration
 - [ ] Add validation with [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
 - [ ] Slack OAuth setup
 - [ ] Integrate with Github API
 - [ ] Add easy setup of git hooks
 - [ ] Add a local/ci mode so that certain actions aren't triggered when running local (ex. slack) kind of like a log only mode
 - [ ] Add an action filter to CLI so only user specific action types are run
 - [ ] CLI to list configured watchers & actions
 - [ ] CLI to test actions
 - [ ] CLI specify .diffhook.yml file to use
 - [ ] Generate update for watchers when lines changes
 - [ ] Comment support? Add a tag as comment in code to watch it?