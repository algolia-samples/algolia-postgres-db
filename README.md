# Sample application template

This is a template you can use to create other Algolia sample applications. It contains a variety of features that every Algolia sample app should ideally include. You can use the [Github repository template](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) functionality to create your sample app from this template.

## Features

The sample app uses the following features:

- Three back-end implementations in different languages
- ...

## Demo (Try it yourself!)

Adding a live demo (e.g., on [CodeSandbox](https://codesandbox.io/)) will let the people quickly test your sample application!

## How to run the sample app locally

The sample app implements three servers in the following programming languages:

- [Python](server/python)
- [Node.js/JavaScript](server/node)
- [Go](server/go)

The [client](client) is a single HTML page.

### 1. Clone this repository

```
git clone https://github.com/algolia-samples/chatbot-with-algolia-answers
```

Copy the file `.env.example` to the directory of the server you want to use and rename it to `.env`. For example, to use the Python implementation:

```bash
cp .env.example server/python/.env
```

### 2. Set up Algolia

To use this sample app, you need an Algolia account. If you don't have one already, [create an account for free](https://www.algolia.com/users/sign-up). Note your [Application ID](https://deploy-preview-5789--algolia-docs.netlify.app/doc/guides/sending-and-managing-data/send-and-update-your-data/how-to/importing-with-the-api/#application-id).

In the `.env` file, set the environment variables `ALGOLIA_APP_ID`:

```bash
ALGOLIA_APP_ID=<replace-with-your-algolia-app-id>
```

### 3. Create your Algolia index and upload data

After you set up your Algolia account and Algolia application, [create and populate an index](https://www.algolia.com/doc/guides/sending-and-managing-data/prepare-your-data/).

To upload your data, you can use the [Algolia dashboard](https://www.algolia.com/doc/guides/sending-and-managing-data/send-and-update-your-data/how-to/importing-from-the-dashboard/) or use on of Algolia's [API clients](https://www.algolia.com/developers/#integrations).

After creating the index and uploading the data, set the environment variables `ALGOLIA_INDEX_NAME` and `ALGOLIA_API_KEY` in the `.env` file:

```bash
ALGOLIA_INDEX_NAME=<replace-with-your-algolia-index-name>
ALGOLIA_API_KEY=<replace-with-your-algolia-api-key>
```

### 6. Follow the instructions in the server directory

Each server directory has a file with instructions:

- [Node.js](server/node/README)
- [Python](server/python/README)
- [Go](server/go/README)

For example, to run the Python implementation of the server, follow these steps:

```bash
cd server/python # there's a README in this folder with instructions
python3 venv env
source env/bin/activate
pip3 install -r requirements.txt
export FLASK_APP=server.py
python3 -m flask run --port=4242
```

## Resources

- [GitHub's repository template](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-repository-from-a-template) functionality

## Contributing

This template is open source and welcomes contributions. All contributions are subject to our [Code of Conduct](https://github.com/algolia-samples/.github/blob/master/CODE_OF_CONDUCT.md).

## Authors

- [@cdenoix](https://twitter.com/cdenoix)
