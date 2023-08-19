import configparser
import os
from pprint import pprint

from atlassian import Bitbucket


def read_config():
    config_file = configparser.ConfigParser()
    config_file_path = os.path.join(os.path.dirname(__file__), '..', 'config.ini')
    config_file.read(config_file_path)
    return config_file


if __name__ == "__main__":
    config = read_config()
    jira_url = config["bitbucket"]["jira"]
    bitbucket_url = config["bitbucket"]["url"]

    bitbucket = Bitbucket(url=bitbucket_url,
                          username=config["bitbucket"]["username"],
                          password=config["bitbucket"]["password"])

    pull_requests = bitbucket.check_inbox_pull_requests(role='author')
    # TODO: Find out and save how many comments are made on each pull request
        


