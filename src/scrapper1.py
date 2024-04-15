import requests
import json
from bs4 import BeautifulSoup

URL = "https://immobilien.bwsg.at"

PARAMS = {
    'f[all][marketing_type]': 'rent',
	'f[all][realty_type][0]': '2',
    'f[all][realty_type][1]': '3',
    'f[all][city]': 'Wien',
	'from': '1117350'
}

req = requests.request(method='GET', url=URL, params=PARAMS)

soup = BeautifulSoup(req.text, 'html.parser')
pages = soup.find(class_ = 'pagination')
cur = pages.find(class_ = 'active')
pages = pages.find_all('li')

link = []

while True:

    panel_wrapper = soup.find_all(class_='panel-wrapper')

    for panel in panel_wrapper:
        panel_footer = panel.find(class_='panel-footer')
        link.append(panel.find('a').get('href'))
    
    if cur == pages[-1]:
        break
            
    cur_index = pages.index(cur)
    next_index = cur_index + 1
    next_link = pages[next_index].find('a').get('href')
    
    req = requests.request(method='GET', url=URL + next_link)
    soup = BeautifulSoup(req.text, 'html.parser')
    pages = soup.find(class_ = 'pagination')
    cur = pages.find(class_ = 'active')
    pages = pages.find_all('li')
    
print(len(link))

    

    



print(link)
