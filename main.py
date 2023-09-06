import argparse
import os

import requests
from lxml import etree


def writer(url, desc, demoUrl):
    filename = os.getcwd() + '/README.md'
    with open(filename, encoding='utf-8', mode='r') as f:
        oldContent = f.read()

        if r"{}".format(oldContent).find(url) == -1:

            if demoUrl is not None:
                demo = f"[demo]({demoUrl})"
                item = r"|{}|{}|{}|".format(demo, url, desc)
            else:
                item = r"||{}|{}|".format(url, desc)

            newContent = r"{}{}{}".format(oldContent, item, "<br>\n")
            with open(filename, encoding='utf-8', mode='w') as fin:
                fin.write(newContent)


def check_url(uri):
    url = uri
    if not str(uri).startswith("http"):
        url = "https://" + uri

    response = requests.get(url=url).text.encode('utf-8')
    tree = etree.HTML(response)
    title = tree.xpath('//*[@id="responsive-meta-container"]/div/p/text()')

    return str(title[0]).strip()


if __name__ == '__main__':

    parse = argparse.ArgumentParser()
    parse.add_argument('--url', type=str)
    parse.add_argument('--description', type=str, default=None)
    parse.add_argument('--demo_url', type=str, default=None)
    args = parse.parse_args()

    url = args.url
    desc = args.description
    if desc is None:
        desc = check_url(url)

    url = str(url).replace('https://', '').replace('http://', '')

    writer(url, desc, args.demo_url)
