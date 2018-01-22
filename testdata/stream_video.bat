ECHO ffmpeg -i %~dp0SampleVideo_1280x720_5mb.mp4 -f mpegts -codec:v mpeg1video -s 1280x720 -rtbufsize 2048M -r 30 -b:v 3000k  -q:v 6 http://localhost:8081/secret1234

ffmpeg -i "%~dp0/SampleVideo_1280x720_5mb.mp4" -f mpegts -codec:v mpeg1video -s 1280x720 -rtbufsize 2048M -r 30 -b:v 3000k  -q:v 6 http://localhost:8081/secret1234