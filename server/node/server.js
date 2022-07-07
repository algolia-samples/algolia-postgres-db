const express = require("express");
const { resolve } = require("path");
const bodyParser = require("body-parser");
const algoliaearch = require("algoliasearch");

const envFilePath = resolve(__dirname, './.env');
require("dotenv").config({ path: envFilePath });

const algoliaClient = algoliaearch(process.env.ALGOLIA_APP_ID, process.env.ALGOLIA_API_KEY)
const index = algoliaClient.initIndex(process.env.ALGOLIA_INDEX_NAME);

const app = express();

app.use(express.static(process.env.STATIC_DIR));
app.use(bodyParser.json())

app.get("/", (req, res) => {
    const path = resolve(process.env.STATIC_DIR + "/index.html");
    res.sendFile(path);
});

app.listen(4242, () => console.log(`Node server listening on port ${4242}!`));