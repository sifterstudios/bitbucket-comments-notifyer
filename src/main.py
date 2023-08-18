from atlassian import Bitbucket
import configparser
import os

def read_config():
    config = configparser.ConfigParser()
    config_path = os.path.join(os.path.dirname(__file__), '..', 'config.ini')
    config.read(config_path)
    return config

def fetch_comments(bitbucket, repo_slug):
    # Fetch the list of comments from a repository
    comments = bitbucket.comments.get_comments(repo_slug=repo_slug)
    return comments

if __name__ == "__main__":
    repository_slug = "repo-name"

    config = read_config()
    bitbucket_url = config["bitbucket"]["url"]
    access_token = config["bitbucket"]["access_token"]

    try:
        bitbucket = Bitbucket(url=bitbucket_url, access_token=access_token)
        comments = fetch_comments(bitbucket, repository_slug)
        
        for comment in comments:
            print(f"Author: {comment['author']['user']['displayName']}")
            print(f"Comment: {comment['text']}")
            print("=" * 30)
    except Exception as e:
        print("An error occurred:", e)
