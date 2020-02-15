# Twitchbot Notes

## Features
- ~~Respond to a the !hello bot command~~
- ~~!project command -> give a description of what I'm currently working on~~
- Greet someone when they enter the chatroom
- Send a welcome message with channel info on an interval


- Web API to set configuration
- Web interface for stats + configuration
    - Set project info
    - Set greeting
    - Stream logs

- Deploy somewhere
    - Digitalocean?
    - Github action

- Refactors
    - Logging library

## Implementation
- Written using go 1.13 with go modules
- Start with basic Twitch IRC implementation, move to some twitchbot library

## Prerequisites
- Personal Twitch account
- Bot Twitch account
- Bot API key
- Go installed
