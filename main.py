import argparse
import os
from lxml import etree

import requests


def get_filename(language):
    if 'Go' == language:
        return os.getcwd() + '/README.md'
    elif 'PHP' == language:
        return os.getcwd() + '/README_PHP.md'
    else:
        return os.getcwd() + '/README_OTHER.md'


def writer(uri, description, language, demoUrl):
    filename = get_filename(language)
    uri = f'{uri}'.replace('https://', '')
    with open(filename, encoding='utf-8', mode='r') as f:
        oldContent = f.read()

        if r"{}".format(oldContent).find(uri) == -1:

            if language == 'Go':
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
    if not str(uri).startswith("http"):
        uri = "https://" + uri

    response = requests.get(url=uri).text.encode('utf-8')
    tree = etree.HTML(response)
    title = tree.xpath('//*[@id="responsive-meta-container"]/div/p/text()')

    return str(title[0]).strip()


if __name__ == '__main__':
    parse = argparse.ArgumentParser()
    parse.add_argument('--url', type=str)
    parse.add_argument('--description', type=str, default='')
    parse.add_argument('--language', type=str, default='Go')
    parse.add_argument('--demo_url', type=str, default='')
    args = parse.parse_args()

    url = args.url
    desc = args.description
    if len(desc) < 1:
        desc = get_description(url)

    writer(url, desc, args.language, args.demo_url)
