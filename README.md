<div align="center">
  <a href="https://dankmuzikk.com" target="_blank"><img src="https://dankmuzikk.com/static/android-chrome-512x512.png" width="150" /></a>

  <h1>DankMuzikk</h1>
  <p>
    <strong>Create, Share and Play Music Playlists.</strong>
  </p>
  <p>
    <a href="https://goreportcard.com/report/github.com/mbaraa/dankmuzikk"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/mbaraa/dankmuzikk"/></a>
    <a href="https://github.com/mbaraa/dankmuzikk/actions/workflows/rex-deploy.yml"><img alt="rex-deployment" src="https://github.com/mbaraa/dankmuzikk/actions/workflows/rex-deploy.yml/badge.svg"/></a>
    <a href="https://github.com/mbaraa/dankmuzikk/actions/workflows/rex-deploy-beta.yml"><img alt="rex-deployment" src="https://github.com/mbaraa/dankmuzikk/actions/workflows/rex-deploy-beta.yml/badge.svg"/></a>
  </p>
</div>

## About

**DankMuzikk** is a music player designed for colloborative playlists, where a playlist's collaborators can add and vote for music in a playlist, and the other music player stuff.

_Note: this is a fling side-project that could die anytime so don't get your hopes up._

## Contributing

IDK, it would be really nice of you to contribute, check the poorly written [CONTRIBUTING.md](/CONTRIBUTING.md) for more info.

## Roadmap

- [x] Search YouTube for music
- [x] Web UI
- [x] Audio player
- [x] Accounts and Profiles
- [x] Playlists
- [x] Share playlists
- [x] Vote songs in playlists
- [ ] Songs queue
  - [x] Add songs to queue
  - [ ] View/Edit queue
- [ ] Expandable player
  - [x] Mobile
  - [ ] Desktop 
- [ ] Share songs
- [ ] Player's menu
- [ ] Import **public** playlists from YouTube
- [ ] Update profile
- [ ] Cross device control
- [ ] Offline support
- [ ] Songs' metadata fixer
- [ ] Lyrics
- [ ] Refactor the code (never gonna happen)

## Technologies used

- **[Go](https://golang.org)**: Main programming language.
- **[JavaScript](https://developer.mozilla.org/en-US/docs/Web/javascript)**: Dynamic client logic.
- **[Python](https://python.org)**: Used for the YouTube download service.
- **[templ](https://templ.guide)**: The better [html/template](https://pkg.go.dev/html/template).
- **[htmx](https://htmx.org)**: The front-end library of peace.
- **[GORM](https://gorm.io)**: The fantastic ORM library for Golang.
- **[MariaDB](https://mariadb.org)**: Relational database.
- **[yt-dlp](https://github.com/yt-dlp/yt-dlp)**: YouTube download helper.
- **[pytube](https://github.com/pytube/pytube)**: YouTube download helper.
- **[minify](https://github.com/tdewolff/minify)**: Minify static text files.

## Run locally

1. Clone the repo.

```bash
git clone https://github.com/mbaraa/dankmuzikk
```

2. Create the docker environment file

```bash
cp .env.example .env.docker
```

3. Run it with docker compose.

```bash
docker compose up -f docker-compose-dev.yml
```

3. Visit http://localhost:20250
4. Don't ask why I chose this weird port.

## Acknowledgements

- **This project is not affiliated with YouTube or Google, or anyone to that matter in any sort of ways.**
- The background was taken from [dankpods.net](https://dankpods.net)
- Frank’s original image was taken from [dingusland.biz](https://dingusland.biz)
- Colorscheme is inspired from [Dankpods](https://www.youtube.com/@DankPods)
- youtube-scrape was used to search videos without using the actual YouTube API (small quota): MIT licenses by [Herman Fassett](https://github.com/HermanFassett)

---

Made with 🧉 by [Baraa Al-Masri](https://mbaraa.com)
