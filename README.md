# GoRecordurbate WebUI
This project offers a simple, self-hosted web UI for managing and recording streams.

__API NOTE__: The api is basically un-tested after adding login credentials. The plan is to use API key in the future for credentials via API. This will probably not be done before the initial release.
## Core Features:
- Zero delay when checking online status. New method in use that is not restricted by rate-limit.
- Login with cookies for each client (prevents unsupervised access)
- Start, restart, and stop recordings
- Add/delete streamers, with import/export options
- View logs and recorded videos directly in the web UI
- Docker and service examples for easier setup
- Manage and add multiple users for webUI login
- Add new streamers without restarting the recorder
- View active recorders
- Watch live streams directly in the web UI
  
## Usage

|Username|Password|
|-|-|
|`admin`| `password`|

### Setup
- Copy [`.env.example`](https://github.com/luna-nightbyte/Recordurbate-WebUI/blob/main/.env.example) to `.env` and add your own session key. 
    - Can be generated with: `head -c 32 /dev/urandom | base64`
      I.e:
      ```bash
      user@user:~/Recordurbate$ head -c 32 /dev/urandom | base64
      # Output:
      Fl60B6sTqAUARyDiC6GIor8AIu6QXLF2RMvWK1Wz3eE=
      ```

#### Optional config settings
The main settings can be found in [`settings.json`](https://github.com/luna-nightbyte/Recordurbate-WebUI/blob/main/internal/settings/settings.json):
```json
{
  "app": {
    "port": 8055,
    "loop_interval_in_minutes": 2,
    "video_output_folder": "output/videos",
    "rate_limit": {
      "enable": true,
      "time": 5
    },
    "default_export_location": "./output/list.txt"

  },
  "youtube-dl": {
    "binary": "youtube-dl"
  },
  "auto_reload_config": true
}
```
#### Reset password
To change forgotten password, start the program with the `reset-pwd` argument. I.e:
```
./GoRecordurbate reset-pwd admin newpassword 
```
New login for the user `admin` would then be `newpassword`
### Build
Building the code wil create a binary for your os system. Golang is [cross-compatible](https://go.dev/wiki/GccgoCrossCompilation) for windows, linux and mac.
```bash
go mod init GoRecordurbate # Only run this line once
go mod tidy
go build
./GoRecordurbate #windows will have 'GoRecordurbate.exe'
```
### Source
```bash
go mod init GoRecordurbate # Only run this line once
go mod tidy
go run main.go
```

## Notes
This is un-tested on Windows and Mac, but golang is cross-compatible which means that this should run just as fine on Windows as on Linux.

## WebUI (Will probably be modified)


<p align="center">
  <img src="https://github.com/user-attachments/assets/35e4633b-702b-45f9-9075-a8522a6b334b" alt="Login page"/>

  
  <img src="https://github.com/user-attachments/assets/b9419caf-f2b9-4f4f-a8a0-ddd490bc9cef" alt="Control tab"/>
  <img src="https://github.com/user-attachments/assets/fa5a9008-b21c-47ef-bb90-bbcb379053bc" alt="User settings tab"/>
  <img src="https://github.com/user-attachments/assets/24744566-bc52-4e80-8504-90d63abb4903" alt="Streamers tab"/>
  <img src="https://github.com/user-attachments/assets/74f7c222-4163-4dd9-9c06-c69044e7c845" alt="Livestream tab"/>
</p>



## Other

### Todo, but not planned/prioritized 
- Log online-time of streamers and save to csv for graph plotting. Can help understand the work-hours of different streamers.
- Option to login to the streaming site and use follower list instead of config?
- Option for max video length (and size?)
- ~~headless mode without webui~~ (Abandoned)
- Move frontend to Vue (?)
- Build a default docker image
- Individual recorders in UI
  - Stop/Restart individual recorders
  - view current recording length
  - view current recording video





## Thanks

Special thanks to [oliverjrose99](https://github.com/oliverjrose99) for the inspiration and their work on [Recordurbate](https://github.com/oliverjrose99/Recordurbate)
