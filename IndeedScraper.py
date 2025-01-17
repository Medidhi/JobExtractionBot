import csv
from datetime import datetime
import requests
from bs4 import BeautifulSoup
import pandas as pd
import sys
import os


"""get_url: function to get the indeed url based on position and location"""
def get_url(position, location):
    url = "https://www.indeed.com/jobs?q={}&l={}".format(position, location)
    return url

"""get_Records: get each record extracted from the html dom"""
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
        job_records.append(record)
    return job_records

"""extract_jobs: funtion to extract jobs"""
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

"""main funtion which invokes inbuilt funtions to run the functionality"""
if __name__ == "__main__":
    position = sys.argv[1].replace("#", " ")
    location = sys.argv[2].replace("#", " ")
    otputFilePath = sys.argv[3].replace("#", " ")

    data = pd.DataFrame(extract_jobs(position, location))
    path = otputFilePath + str(datetime.date(datetime.now())) + "/"

    print(path)
    try:
        os.mkdir(path)
    except FileExistsError:
        print("file already exists")
    data.to_excel((path + "{}_{}.xlsx").format(position, location))