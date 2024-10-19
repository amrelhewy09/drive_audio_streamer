### Google Drive Audio Streamer CLI
I often have many large .wav audio files stored on Google Drive but struggle with slow internet speeds, making it frustrating to wait for entire files to load before I can listen to them. To solve this, I developed a simple command-line interface (CLI) tool that lists all the audio files in my Google Drive and buffers them for playback. This way, I can start listening without waiting for the entire file to download. It's still very beta and needs alot of features but it works :)

### Install
Steps are very simple just before starting make sure you have [ffmpeg](https://github.com/FFmpeg/FFmpeg) installed on your host machine
Also make sure to have a google drive project with a valid `credentials.json` file, visit [here](https://developers.google.com/drive/activity/v2/guides/project) for more info.

1. Clone the repository
2. Paste the `credentials.json` file in the root directory of the project
3. Run `./build/drive_audio_streamer_linux run` or `./build/drive_audio_streamer_mac run` depending on your host os
4. Google oauth flow starts
5. Once successfull authorization, take the authorization code and paste into the terminal
6. You should see a list of your audio files on google drive :)

It will attempt to start an Oauth flow with google to authorize access to google drive. Make sure port `8080` is unused because a local server runs to display the authorization code. You copy paste the code into the termianl and press enter :)

Feel free to suggest any improvements. Hope this helps in any way
