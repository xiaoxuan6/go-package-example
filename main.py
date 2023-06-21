import os
import sys

import requests
from lxml import etree


def writer(url, desc):
    filename = os.getcwd() + '/README.md'
    with open(filename, encoding='utf-8', mode='r') as f:
        oldContent = f.read()

        if r"{}".format(oldContent).find(url) == -1:
            item = r"||{}|{}|".format(url, desc)
            newContent = r"{}{}{}".format(oldContent, item, "<br>\n")

            with open(filename, encoding='utf-8', mode='w') as fin:
                fin.write(newContent)


if __name__ == '__main__':
    url = ''
    desc = ''
    if len(sys.argv) == 3:
        url = str(sys.argv[1]).replace('https://', '').replace('http://', '')
        desc = sys.argv[2]
    elif len(sys.argv) == 2:
        param = sys.argv[1]
        if str(param).startswith("http"):
            url = param
        else:
            url = "https://" + param

        response = requests.get(url=url).text.encode('utf-8')
        tree = etree.HTML(response)
        title = tree.xpath('//*[@id="responsive-meta-container"]/div/p/text()')

        url = str(param).replace('https://', '').replace('http://', '')
        desc = str(title[0]).strip()
    else:
        sys.exit(0)

    writer(url, desc)
