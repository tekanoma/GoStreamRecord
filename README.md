# GoRecordurbate WebUI
This project offers a simple, self-hosted web UI for managing and recording streams.

__API NOTE__: The api is basically un-tested after adding login credentials. The plan is to use API key in the future for credentials via API.
## Core Features:
- Login with cookies for each client (prevents unsupervised access)
- Start, restart, and stop recordings
- Add/delete streamers, with import/export options
- View logs and recorded videos directly in the web UI
- Docker and service examples for easier setup
## In Progress:
- ~~Better handling of default usernames and passwords~~
- Embedding the index file directly into the code
- API secret key
  
## Usage
Default login:

- User: `admin`
- Password: `password`

### Build
Building the code wil create a binary for your os system. Golang is [cross-compatible](https://go.dev/wiki/GccgoCrossCompilation) for windows, linux and mac.
```bash
go mod init GoRecordurbate # Only run this line once
go mod tidy
go build
```
### Source
```bash
go mod init GoRecordurbate # Only run this line once
go mod tidy
go run main.go
```

## Notes
This is un-tested on Windows and Mac, but golang is cross-compatible which means that this should run just as fine on Windows as on Linux.

A release will be made once i have finished fixing the bare minimum below:
- [x] Start recording
- [x] Restart recording
- [x] stop recording
- [x] Add / delete streamer
- [x] Import streamers
- [x] Export streamers
- [x] Show logs in web UI
- [x] Show videos in web UI
- [x] Change username/password directly in the web UI
- [x] Better implementation of default username & password
- [x] Docker example
- [x] Service example
- [ ] Embed frontend files into binary
- [ ] API secret key implementation 

## WebUI (Will probably be modified)


<p align="center">
  <img src="https://github.com/user-attachments/assets/35e4633b-702b-45f9-9075-a8522a6b334b" alt="Login page"/>
  <img src="https://github.com/user-attachments/assets/6c04598e-d3fe-4630-9bfc-e3c1216d67c5" alt="Streamers tab"/>
  <img src="https://github.com/user-attachments/assets/ab28c113-4c6a-4a07-ba88-05ed6b3a868e" alt="Control tab"/>
</p>

## Other

### Ideas, but not planned/prioritized 
- Log online-time of streamers and save to csv for graph plotting. Can help understand the work-hours of different streamers.
- Option to login to the streaming site and use follower list instead of config?
- Option for max video length (and size?)
- headless mode without webui
- Move frontend to Vue (?)




## Thanks

Special thanks to [oliverjrose99](https://github.com/oliverjrose99) for the inspiration and their work on [Recordurbate](https://github.com/oliverjrose99/Recordurbate)
