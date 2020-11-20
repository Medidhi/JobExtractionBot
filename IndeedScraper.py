import csv
from datetime import datetime
import requests
from bs4 import BeautifulSoup
import pandas as pd
import sys

def get_url(position, location):
    url = "https://www.indeed.com/jobs?q={}&l={}".format(position, location)
    return url

def get_Records(cards):
    job_records = []
    for job in cards:
        atag = job.h2.a
        record = {}
        record["titel"] = atag.get('title')
        record["job_url"] = "https://www.indeed.com/" + atag.get('href')
        record["company"] = job.find("span", "company").text.strip()
        record["location"] = job.find("div", "recJobLoc").get("data-rc-loc")
        record["posted_date"] = job.find("span", "date").text
        record["salary"] = job.find("span", "date").text
        job_records.append(record)
    return job_records

def extract_jobs(position, location):
    url = get_url(position, location)
    response = requests.get(url)
    soup = BeautifulSoup(response.text, 'html.parser')
    cards = soup.find_all("div", "jobsearch-SerpJobCard")

    jobrecord = []
    jobrecord.extend(get_Records(cards))
    while True:
        try:
            url = "https://www.indeed.com" + soup.find("a", {"arial-label": "Next"}).get("href")
        except AttributeError:
            break
        response = requests.get(url)
        soup = BeautifulSoup(response.text, 'html.parser')
        cards = soup.find_all("div", "jobsearch-SerpJobCard")
        jobrecord.extend(get_Records(cards))
    return jobrecord

if name == "__main__":
    data = pd.DataFrame(extract_jobs(sys.argv[1], sys.argv[2]))
    data.to_excel((sys.argv[3]+"{}_{}_{}.xlsx").format(datetime.today(), sys.argv[1], sys.argv[2]))