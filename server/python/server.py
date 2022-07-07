"""Sample application template in Python"""

import os

from algoliasearch.search_client import SearchClient
from dotenv import load_dotenv, find_dotenv
from flask import Flask, render_template

load_dotenv(find_dotenv())

# Setup the Algolia client
algolia = SearchClient.create(
    app_id=os.getenv('ALGOLIA_APP_ID'),
    api_key=os.getenv('ALGOLIA_API_KEY')
)
employees_index = algolia.init_index(os.getenv('ALGOLIA_INDEX_NAME'))

# Setup Flask
STATIC_DIR = str(
    os.path.abspath(os.path.join(
        __file__,
        '..',
        os.getenv('STATIC_DIR')
    ))
)
app = Flask(
    __name__,
    static_folder=STATIC_DIR,
    static_url_path="",
    template_folder=STATIC_DIR
)

@app.route('/', methods=['GET'])
def scanner():
    """Display the index."""
    return render_template('index.html')