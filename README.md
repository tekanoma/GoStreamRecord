# GoRecordurbate

## Intro
A golang version of recordurbate with some differences. One key difference is that this doesent use a deamon. It can instead be compiled into a binary file and started as a service, docker container or whatever you prefer. 

### Notes
This is un-tested on Windows and Mac, but golang is cross-compatible which means that this should run just as fine on Windows as on Linux.

A release will be made once i have finished fixing the bare minimum below:
- [x] Start recording
- [ ] Restart recording
- [ ] stop recording (Might be dropped since the recording already stops when signal interrupt is sent.)
- [x] Add streamer
- [x] Remove streamer
- [ ] Import streamers
- [ ] Export streamers
- [ ] Docker example
- [ ] Service example
## Thanks

Special thanks to oliverjrose99 for the inspiration and their work API information in the source code of [Recordurbate](https://github.com/oliverjrose99/Recordurbate)
