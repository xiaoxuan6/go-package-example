import argparse
import os
from urllib.parse import urlparse

import requests


def parse_url(uri):
    parsed_url = urlparse(uri)
    path = parsed_url.path.strip("/")
    return path.split("/")


def get_filename(owner, repo):
    uri = 'https://api.github.com/repos/{}/{}/languages'.format(owner, repo)
    response = requests.get(uri).json()
    print('response', response)

    if 'Go' in response:
        return os.getcwd() + '/README.md'
    elif 'PHP' in response:
        return os.getcwd() + '/README_PHP.md'
    else:
        return os.getcwd() + '/README_OTHER.md'


def writer(uri, description, demoUrl):
    paths = parse_url(uri)
    filename = get_filename(paths[0], paths[1])

    uri = f'{uri}'.replace('https://', '')
    with open(filename, encoding='utf-8', mode='r') as f:
        oldContent = f.read()

        if r"{}".format(oldContent).find(uri) == -1:

            if f'{filename}'.__contains__('README.md'):
                if len(demoUrl) > 1:
                    demo = f"[demo]({demoUrl})"
                    item = r"|{}|{}|{}|".format(demo, uri, description)
                else:
                    item = r"||{}|{}|".format(uri, description)
            else:
                item = r"|{}|{}|".format(uri, description)

            newContent = r"{}{}{}".format(oldContent, item, "<br>\n")
            with open(filename, encoding='utf-8', mode='w') as fin:
                fin.write(newContent)


def get_description(uri):
    paths = parse_url(uri)
    uri = 'https://api.github.com/repos/{}/{}'.format(paths[0], paths[1])
    response = requests.get(uri)
    if response.status_code == 200:
        return response.json()['description']


if __name__ == '__main__':
    parse = argparse.ArgumentParser()
    parse.add_argument('--url', type=str)
    parse.add_argument('--description', type=str, default='')
    parse.add_argument('--demo_url', type=str, default='')
    args = parse.parse_args()

    url = args.url
    desc = args.description
    if len(desc) < 1:
        desc = get_description(url)

    writer(url, desc, args.demo_url)
