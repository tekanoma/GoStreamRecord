# GoRecordurbate WebUI

Readme is incomplete and will be completely refactored at a later point. Still relevant tho. 

Checkput the API doc for info on how to use the api.
## Intro
A golang version of recordurbate with some differences. One key difference is that this doesent use a deamon, but instead runs a webserver (Web UI). It can instead be compiled into a binary file and started as a service, docker container or whatever you prefer. 

Sends a request to check if the [spesific](https://github.com/luna-nightbyte/GoRecordurbate/blob/ec0b1fa79e2bb82cf948bef3415ace3aac52e523/modules/bot/bot.go#L176) user is online rather than requesting a [list of 500](https://github.com/luna-nightbyte/GoRecordurbate/blob/ec0b1fa79e2bb82cf948bef3415ace3aac52e523/modules/bot/bot.go#L175) and checking that for the correct user. 
### Notes
This is un-tested on Windows and Mac, but golang is cross-compatible which means that this should run just as fine on Windows as on Linux.

A release will be made once i have finished fixing the bare minimum below:

- [x] Start recording
- [x] Restart recording
- [x] stop recording
- [x] Add streamer
- [x] Remove streamer
- [ ] Import streamers
- [ ] Export streamers
- [ ] Docker example
- [ ] Service example
- [ ] Embed index file into code
#### WebUI
- [x] Start, stop, restart
- [x] Add / delete streamer
- [ ] Import streamers
- [ ] Export streamers
- [ ] Show log
- [ ] View videos in web ui

## WebUI (Will probably be modified)
![image](https://github.com/user-attachments/assets/a02e40a1-1a39-4cd4-9b53-a8b88568f38b)
![image](https://github.com/user-attachments/assets/9aabf47f-62eb-4065-b7ef-278cb98bd916)
![image](https://github.com/user-attachments/assets/3e0f4f3a-dd0b-42cc-8929-1439b26495fe)
![image](https://github.com/user-attachments/assets/b1ce631d-4d1a-4ffb-928a-e1f524bc327d)
![image](https://github.com/user-attachments/assets/a0325883-ffe8-4e5a-96dc-7727e1b50380)
![image](https://github.com/user-attachments/assets/117803c1-031d-4a9c-8173-7ea983b064f4)







## Thanks

Special thanks to oliverjrose99 for the inspiration and their work on [Recordurbate](https://github.com/oliverjrose99/Recordurbate)
