from flask import Flask
import mariadb
import os
import os.path
import signal
import threading
import time
from yt_dlp import YoutubeDL


##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################
## Environmental variables
##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################

def get_env(key) -> str:
    val = os.environ.get(key)
    if val is None or val == "":
        print(f"Missing {key} suka")
        exit(1)
    return val


DB_NAME = get_env("DB_NAME")
DB_HOST = get_env("DB_HOST")
DB_USERNAME = get_env("DB_USERNAME")
DB_PASSWORD = get_env("DB_PASSWORD")
DOWNLOAD_PATH = get_env("YOUTUBE_MUSIC_DOWNLOAD_PATH")

##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################
## DB

## Execute those on the target's mariadb
##
## SET @@GLOBAL.wait_timeout=31536000;
## SET @@GLOBAL.interactive_timeout=31536000;

##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################

conn = None


def open_db_conn():
    try:
        db_host = DB_HOST
        db_port = 3306
        if ":" in db_host:
            db_host = DB_HOST[:DB_HOST.index(":")]
            db_port = int(DB_HOST[DB_HOST.index(":") + 1:])
        global conn
        conn = mariadb.connect(
            user=DB_USERNAME,
            password=DB_PASSWORD,
            host=db_host,
            port=db_port,
            database=DB_NAME
        )
    except mariadb.Error as e:
        print(f"Error connecting to MariaDB Platform: {e}")
        return 1


def update_song_status(id: str):
    cur = conn.cursor()
    try:
        cur.execute("UPDATE songs SET fully_downloaded=1 WHERE yt_id=?", (id,))
        conn.commit()
    finally:
        cur.close()


def song_exists(song_yt_id: str) -> bool:
    cur = conn.cursor()
    try:
        cur.execute("SELECT id FROM songs WHERE yt_id=? AND fully_downloaded=1", (song_yt_id,))
        result = cur.fetchone()
        return result[0] if result else False
    except:
        cur.close()
        return False
    finally:
        cur.close()
    return False


open_db_conn()

##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################
## Downloader
##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################

YT_ERROR = {
    0: "none",
    1: "age restiction",
    2: "video unavailable",
    3: "other youtube error",
}

def download_yt_song(id: str) -> int:
    try:
        ytdl = YoutubeDL({
            "format": "bestaudio/mp3",
            "postprocessors": [{
                "key": "FFmpegExtractAudio",
                "preferredcodec": "mp3",
                "preferredquality": "192",
            }],
            "outtmpl": f"{DOWNLOAD_PATH}/%(id)s.%(ext)s"
        })
        ytdl.download("https://www.youtube.com/watch?v=" + id)
    except:
        return 3

    return 0


##############################################################################################################################################################################################################################

to_download_lock = threading.Lock()
to_download_stop_event = threading.Event()
to_download_queue = set([])

currently_downloading_lock = threading.Lock()
currently_downloading_stop_event = threading.Event()
currently_downloading_queue = set([])


def background_task():
    while not to_download_stop_event.is_set():
        with to_download_lock:
            if to_download_queue:
                song_yt_id = to_download_queue.pop()
                print(f"Downloading {song_yt_id} from the queue.")
                res = download_song(song_yt_id)
                if res != 0:
                    print(f"Error downloading {song_yt_id}, error: {YT_ERROR[res]}")
        time.sleep(0.5)


def add_song_to_queue(song_yt_id: str) -> int:
    """
        add_song_to_queue adds a song's id to the download queue.
    """
    with to_download_lock:
        if song_exists(song_yt_id):
            print(f"The song with id {song_yt_id} was already downloaded 😬")
            return 0

        to_download_queue.add(song_yt_id)
        print(f"Added song {song_yt_id} to the download queue.")
    return 0


def download_song(song_yt_id: str) -> int:
    """
        download_song downloads the given song's ids using yt_dlp,
        and returns the operation's status code.
    """
    if song_exists(song_yt_id):
        print(f"The song with id {song_yt_id} was already downloaded 😬")
        return 0

    if not currently_downloading_stop_event.is_set():
        with currently_downloading_lock:
            print(f"Downloading song with id {song_yt_id} ...")
            while song_yt_id in currently_downloading_queue:
                print("waiting suka")
                time.sleep(0.5)
                pass

            currently_downloading_queue.add(song_yt_id)
            res = download_yt_song(song_yt_id)
            currently_downloading_queue.remove(song_yt_id)
            if res != 0:
                print(f"error: {YT_ERROR[res]} when downloading {song_yt_id}")
                return res
            update_song_status(song_yt_id)
            print("Successfully downloaded " + song_yt_id)
            return 0
    return 3


thread = threading.Thread(target=background_task)
thread.start()

##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################
##############################################################################################################################################################################################################################

app = Flask(__name__)


@app.route("/download/queue/<id>")
def handle_add_download_song_to_queue(id):
    res = add_song_to_queue(id)
    if res != 0:
        return {"error": YT_ERROR[res]}
    return {"msg": "woohoo"}


@app.route("/download/<id>")
def handle_download_song(id):
    res = download_song(id)
    if res != 0:
        return {"error": YT_ERROR[res]}
    return {"msg": "woohoo"}


def close_server(arg1, arg2):
    print("signal shit", arg1, arg2)
    print("Stopping background download thread...")
    global to_download_stop_event
    to_download_stop_event.set()
    global currently_downloading_stop_event
    currently_downloading_stop_event.set()
    global thread
    thread.join()
    print("Closing MariaDB's connection...")
    global conn
    conn.close()
    exit(0)


signal.signal(signal.SIGINT, close_server)
signal.signal(signal.SIGTERM, close_server)

if __name__ == '__main__':
    app.run(port=4321)
