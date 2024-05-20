from yt_dlp import YoutubeDL
from yt_dlp.utils import DownloadError
from fastapi import FastAPI, status, Response
from fastapi.requests import Request
import os.path
from threading import Lock, Thread
import time
import signal
import mariadb
import sys

DOWNLOAD_PATH = os.environ.get("YOUTUBE_MUSIC_DOWNLOAD_PATH")
DB_NAME     = os.environ.get("DB_NAME")
DB_HOST     = os.environ.get("DB_HOST")
DB_USERNAME = os.environ.get("DB_USERNAME")
DB_PASSWORD = os.environ.get("DB_PASSWORD")

## DB stuff
conn = None

def open_db_conn():
    try:
        db_host = DB_HOST
        db_port = 3306
        if ":" in db_host:
            db_host = DB_HOST[:DB_HOST.index(":")]
            db_port = int(DB_HOST[DB_HOST.index(":")+1:])
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
    cur.execute("UPDATE songs SET fully_downloaded=1 WHERE yt_id=?", (id,))
    conn.commit()


def song_exists(id: str) -> bool:
    cur = conn.cursor()
    cur.execute("SELECT id FROM songs WHERE yt_id=? AND fully_downloaded=1", (id,))
    result = cur.fetchone()
    return result[0] if result else False

## Download Video stuff

class MutexArray:
    def __init__(self, initial_array: []):
        self._lock = Lock()
        self._array = initial_array.copy()

    def exists(self, item) -> bool:
        with self._lock:
            return item in self._array

    def get(self, index):
        with self._lock:
            return self._array[index]

    def set(self, index, value):
        with self._lock:
            self._array[index] = value

    def append(self, value):
        with self._lock:
            self._array.append(value)

    def remove(self, value):
        with self._lock:
            self._array.remove(value)

    def get_array_and_clear(self):
        with self._lock:
            clone = self._array.copy()
            self._array.clear()
            return clone

    def length(self):
        with self._lock:
            return len(self._array)

    def release(self):
        self._lock.release()

background_download_list = MutexArray([])
to_be_downloaded = MutexArray([])


ytdl = YoutubeDL({
    "format": "bestaudio/best",
    "postprocessors": [{
        "key": "FFmpegExtractAudio",
        "preferredcodec": "mp3",
        "preferredquality": "192",
    }],
    "outtmpl": f"{DOWNLOAD_PATH}/%(id)s.%(ext)s"
})


def download_song(id: str) -> int:
    """
        download_song downloads the given song's ids using yt_dlp,
        and returns the operation's status code.
    """
    try:
        if id is None or len(id) == 0:
            return

        ## wait list
        while to_be_downloaded.exists(id):
            time.sleep(1)
            pass

        ## download the stuff
        if song_exists(id):
            to_be_downloaded.remove(id)
            return 0

        to_be_downloaded.append(id)
        ytdl.download(f"https://www.youtube.com/watch?v={id}")
        to_be_downloaded.remove(id)
        update_song_status(id)

        return 0
    except DownloadError:
        return 1
    except Exception:
        return 2


def download_songs_from_queue():
    """
        download_songs_from_queue fetches the current songs in the download queue,
        and starts the download process.
    """
    if background_download_list.length() == 0:
        return
    for id in background_download_list.get_array_and_clear():
        download_song(id)


def add_song_to_queue(id: str):
    """
        add_song_to_queue adds a song's id to the download queue.
    """
    background_download_list.append(id)


## BG downloader thread

def download_songs_in_background(interval=1):
    """
        download_songs_in_background runs every given interval time in seconds (default is 1),
        and downloads the songs in the queue in the background.
    """
    while True:
        download_songs_from_queue()
        time.sleep(interval)


download_thread = Thread(target=download_songs_in_background, args=(1,))

## FastAPI Stuff

app = FastAPI(
    title="DankMuzikk's YouTube Downloader",
    description="Apparently the CLI's overhead and limitation has got the best of me.",
)


@app.on_event("startup")
def on_startup():
    open_db_conn()
    global download_thread
    download_thread.start()


@app.on_event("shutdown")
def on_shutdown():
    print("Stopping background download thread...")
    background_download_list.release()
    to_be_downloaded.release()
    download_thread.join()
    print("Closing MariaDB's connection...")
    conn.close()


@app.get("/download/queue/{id}", status_code=status.HTTP_200_OK)
def handle_add_download_song_to_queue(id: str, response: Response):
    add_song_to_queue(id)


@app.get("/download/{id}", status_code=status.HTTP_200_OK)
def handle_download_song(id: str, response: Response):
    err = download_song(id)
    if err != 0:
        response.status_code = status.HTTP_400_BAD_REQUEST


@app.get("/download/multi/{ids}",  status_code=status.HTTP_200_OK)
def handle_download_songs(ids: str, response: Response):
    for id in ids.split(","):
        err = download_song(id)
        if err != 0:
            response.status_code = status.HTTP_400_BAD_REQUEST


@app.get("/download/queue/multi/{ids}",  status_code=status.HTTP_200_OK)
def handle_add_download_songs_to_queue(ids: str, response: Response):
    for id in ids.split(","):
        add_song_to_queue(id)

