import os
import sys


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
    if len(sys.argv) > 1:
        writer(sys.argv[1], sys.argv[2])
