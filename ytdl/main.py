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
        global conn
        conn = mariadb.connect(
            user=DB_USERNAME,
            password=DB_PASSWORD,
            host=DB_HOST[:DB_HOST.index(":")],
            port=3307,
            database=DB_NAME
        )
    except mariadb.Error as e:
        print(f"Error connecting to MariaDB Platform: {e}")
        return 1


def update_song_status(id: str):
    cur = conn.cursor()
    cur.execute("UPDATE songs SET fully_downloaded=1 WHERE yt_id=?", (id,))
    conn.commit()


## Download Video stuff

class MutexArray:
    def __init__(self, initial_array: []):
        self._lock = Lock()
        self._array = initial_array.copy()

    def get(self, index):
        with self._lock:
            return self._array[index]

    def set(self, index, value):
        with self._lock:
            self._array[index] = value

    def append(self, value):
        with self._lock:
            self._array.append(value)

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

mutex_array = MutexArray([])


ytdl = YoutubeDL({
    "format": "bestaudio/best",
    "postprocessors": [{
        "key": "FFmpegExtractAudio",
        "preferredcodec": "mp3",
        "preferredquality": "192",
    }],
    "outtmpl": f"{DOWNLOAD_PATH}/%(id)s.%(ext)s"
})


def download_songs(ids: [str]) -> int:
    """
        download_songs downloads the given songs' ids using yt_dlp,
        and returns the operation's status code.
    """
    try:
        new_ids = []
        for id in ids:
            if not os.path.isfile(id+".mp3"):
                new_ids.append(id)

        if len(new_ids) == 0:
            return

        for id in new_ids:
            ytdl.download(f"https://www.youtube.com/watch?v={id}")
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
    if mutex_array.length() == 0:
        return
    download_songs(mutex_array.get_array_and_clear())


def add_song_to_queue(id: str):
    """
        add_song_to_queue adds a song's id to the download queue.
    """
    mutex_array.append(id)


## BG downloader thread

def download_songs_in_background(interval=1):
    """
        download_songs_in_background runs every given interval time in seconds (default is 1),
        and downloads the songs in the queue in the background.
    """
    while True:
        download_songs_from_queue()
        time.sleep(interval)


download_thread = Thread(target=download_songs_in_background, args=(30,))

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
    global mutex_array
    mutex_array.release()
    global download_thread
    download_thread.join()
    print("Closing MariaDB's connection...")
    global conn
    conn.close()


@app.get("/download/queue/{id}", status_code=status.HTTP_200_OK)
def handle_add_download_song_to_queue(id: str, response: Response):
    add_song_to_queue(id)


@app.get("/download/{id}", status_code=status.HTTP_200_OK)
def handle_download_song(id: str, response: Response):
    err = download_songs([id])
    if err != 0:
        response.status_code = status.HTTP_400_BAD_REQUEST


@app.get("/download/multi/{ids}",  status_code=status.HTTP_200_OK)
def handle_download_songs(ids: str, response: Response):
    err = download_songs(ids.split(","))
    if err != 0:
        response.status_code = status.HTTP_400_BAD_REQUEST


@app.get("/download/queue/multi/{ids}",  status_code=status.HTTP_200_OK)
def handle_add_download_songs_to_queue(ids: str, response: Response):
    for id in ids.split(","):
        add_song_to_queue(id)

