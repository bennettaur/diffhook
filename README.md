## Setting Up Slack
* Register a slack app in your org here: https://api.slack.com/apps?new_app=1
  * You can customize the name, if you like, or keep it `changelink`
  * Choose the slack workspace you want changelink to post in
* Under OAuth & Permissions add the following scopes:
  * `chat:write`, `chat:write.public`, `channels:read`, `groups:read`
  * A Future version will provide a CLI for configuring this (https://api.slack.com/authentication/oauth-v2)
* Install the Slack App into your workspace (you should be prompted for this)
  * Select any channel
* Add the bot to any private channels you want it to post to
  * Just send `@changelink` in the channel 
* Use the generated Bot OAuth Token to set the `SLACK_TOKEN` environment variable

## Todo
 - [ ] Documentation for how to setup and use
 - [ ] Add better logging [zerolog?](https://github.com/rs/zerolog)
 - [ ] Split out web server & [mongo](https://github.com/Kamva/mgm) integration
 - [ ] Add validation with [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
 - [ ] Slack OAuth setup
 - [ ] Integrate with Github API