#!/bin/bash
cvlc $1 --stop-time $2 --sout "#transcode{vcodec=none,acodec=mp4a,ab=128,channels=2,samplerate=44100}:file{dst=$3}" vlc://quit