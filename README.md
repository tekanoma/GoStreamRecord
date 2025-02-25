# GoRecordurbate WebUI

## Introduction 

I started this project as my own Recordurbate V2 with webUI because i had some "issues" like long [restart time](https://github.com/oliverjrose99/Recordurbate/issues/77), [not all users are found](https://github.com/oliverjrose99/Recordurbate/issues/76), [status about recorders](https://github.com/oliverjrose99/Recordurbate/issues/75) (Still in developmen). And others that you'll see further down this page.

All of my example commands will be in linux. But its very few, and most if not all are the same on windows/mac.

__API NOTE__: The API is still open after adding adding login. The plan is to use browser cookies and optinally API secret key in the future for credentials via API. But every API endpoint will be reachable through a simple __curl__ command for now. This is not really any issue for anyone who only plan to use this inside their home network. But __Do NOT__ expose the port for this app on your router!! 

## Core Features
### Recorder:
- Streamer status checks bypasses rate limits; recording rate limits still under testing.
- Start, stop, and restart all recordings.
- View active recorders in real-time (More data will be added)
- Stop individual recordings as needed.
- Add or delete streamers dynamically (no need to restart the recorder).
- Import/export streamer lists for easier management.

### WebUI:
- Secure login using cookies for each client (prevents unauthorized access).
- Manage multiple user accounts.
- Directly view logs and recorded videos through the WebUI.
- Watch live streams directly from the WebUI.
### Setup & Deployment:
- Docker configuration and service examples for straightforward deployment.
- Install [Golang](https://go.dev/doc/install) to run the source or build binary.
## Usage

|Username|Password|
|-|-|
|`admin`|see [this](https://github.com/luna-nightbyte/Recordurbate-WebUI/tree/main?tab=readme-ov-file#reset-password)|

### Setup
__important__: You will still need to have the `internal/settings` folder and it's content in the same folder structure when running this app. That means that you'll have to copy that along with any binary you build.

- Download this repo and open a terminal in this folder. Ask ChatGPT how to find the folder path and how to move into it via cli if you dont know.

- Copy [`.env.example`](https://github.com/luna-nightbyte/Recordurbate-WebUI/blob/main/.env.example) to `.env` and add your own session key. (se the one below as long as the app is not exposed outside your local network__IF__ you dont know how to create one. _hint: just ask chatGPT_). 
    - Can be generated with this command on linux.: `head -c 32 /dev/urandom | base64`
      I.e:
      ```bash
      user@user:~/Recordurbate$ head -c 32 /dev/urandom | base64
      # Output
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
- Select and delete videos
- Option to login to the streaming site and use follower list instead of config?
- Option for max video length (and size?)
- ~~headless mode without webui~~ (Abandoned because i will not create all the logic for handling the various arguments myself. Others can create a PR if they want to.)
- Move frontend to Vue
  - Btter for organizing components being re-used
- Build a default docker image
- Individual recorders in UI
  - Start/Restart individual recorders (in progress)
  - set max lenght/size (could be optional to use one of either)
  - view current recording length
  - view current recording video





## Thanks

Special thanks to [oliverjrose99](https://github.com/oliverjrose99) for the inspiration and their work on [Recordurbate](https://github.com/oliverjrose99/Recordurbate)
