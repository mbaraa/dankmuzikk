const express = require("express");
const scraper = require("./scraper");
const app = express();
const process = require("process");

function stop() {
  console.log("gracefully shutting down the server...");
  process.exit();
}

process.on("SIGINT", stop); // Ctrl+C
process.on("SIGTERM", stop); // docker stop

app.get("/api/search", (req, res) => {
  scraper
    .youtube(req.query.q)
    .then((x) => res.json(x))
    .catch((e) => res.send(e));
});

app.listen(8080, () => {
  console.log("Listening on port 8080 (Press CTRL+C to quit)");
});
