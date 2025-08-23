from flask import Flask
import os
import signal
from yt_dlp import YoutubeDL
import shutil

DOWNLOAD_PATH = os.environ.get("BLOBS_DIR")
if DOWNLOAD_PATH is None or DOWNLOAD_PATH == "":
    print("Missing BLOBS_DIR suka")
    exit(1)

MUZIKKX_DIR = f"{DOWNLOAD_PATH}/muzikkx/"
PIX_DIR = f"{DOWNLOAD_PATH}/pix/"

YT_ERROR = {
    0: "none",
    1: "age restriction",
    2: "video unavailable",
    3: "other youtube error",
}

def download_yt_song(id: str) -> int:
    try:
        ytdl = YoutubeDL({
            "format": "251,140,249,250,233,234,139,bestaudio/best",
            "postprocessors": [{
                "key": "FFmpegExtractAudio",
                "preferredcodec": "mp3",
                "preferredquality": "192",
            }],
            "writethumbnail": True,
            "outtmpl": f"{MUZIKKX_DIR}/%(id)s.%(ext)s"
        })
        ytdl.download("https://www.youtube.com/watch?v=" + id)
    except:
        return 3

    try:
        shutil.move(f"{MUZIKKX_DIR}/{id}.jpg", f"{PIX_DIR}/{id}.webp")
    except:
        pass
    try:
        shutil.move(f"{MUZIKKX_DIR}/{id}.webp", f"{PIX_DIR}/{id}.webp")
    except:
        pass
    try:
        shutil.move(f"{MUZIKKX_DIR}/{id}.png", f"{PIX_DIR}/{id}.webp")
    except:
        pass

    return 0


app = Flask(__name__)

@app.route("/download/<id>")
def handle_download_song(id):
    res = download_yt_song(id)
    if res != 0:
        return {"error": YT_ERROR[res]}
    return {"msg": "woohoo"}


def close_server(arg1, arg2):
    print("signal shit", arg1, arg2)
    exit(0)


signal.signal(signal.SIGINT, close_server)
signal.signal(signal.SIGTERM, close_server)

if __name__ == '__main__':
    app.run(port=4321)
