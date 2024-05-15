<div align="center">
  <a href="https://dankmuzikk.com" target="_blank"><img src="https://dankmuzikk.com/static/android-chrome-512x512.png" width="150" /></a>

  <h1>DankMuzikk</h1>
  <p>
    <strong>Create, Share and Play Music Playlists.</strong>
  </p>
  <p>
    <a href="https://goreportcard.com/report/github.com/mbaraa/dankmuzikk"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/mbaraa/dankmuzikk"/></a>
    <a href="https://github.com/mbaraa/dankmuzikk/actions/workflows/rex-deploy.yml"><img alt="rex-deployment" src="https://github.com/mbaraa/dankmuzikk/actions/workflows/rex-deploy.yml/badge.svg"/></a>
  </p>
</div>

## About

**DankMuzikk** is a music player designed for colloborative playlists, where a playlist's collaborators can add and vote for music in a playlist, and the other music player stuff.

_Note: this is a fling side-project that will die in a while so don't get your hopes up._

## Contributing

IDK, it would be really nice of you to contribute, check the poorly written [CONTRIBUTING.md](/CONTRIBUTING.md) for more info.

### Roadmap

- [x] Search YouTube for music
- [x] Web UI
- [x] Audio player
- [x] Accounts and Profiles
- [x] Playlists
- [ ] Share playlists
- [ ] Vote songs in playlists
- [ ] Write a better YouTube scraper (or try to fix the quota thing)
- [ ] Refactor the code (never gonna happen)

## Run locally

1. Clone the repo.

```bash
git clone https://github.com/mbaraa/dankmuzikk
```

2. Create the docker environment file

```bash
cp .env.example .env.docker
```

3. Run it with compose.

```bash
docker compose up
```

3. Visit http://localhost:20250
4. Don't ask why I chose this weird port.

## Acknowledgements

- The background was taken from dankpods.net
- Frank’s original image was taken from dankpods.biz
- Colorscheme is inspired from [Dankpods](https://www.youtube.com/@DankPods)
- templ was used to make this project’s views: MIT licensed by [Adrian Hesketh](https://github.com/a-h)
- htmx was used to make this project’s client more dynamic: No licence but just wanted to say that I used htmx BTW
- youtube-scrape was used to search videos without using the actual YouTube API (small quota): MIT licenses by [Herman Fassett](https://github.com/HermanFassett)
