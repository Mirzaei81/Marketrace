import os 
import requests
# from icrawler.builtin import  BingImageCrawler,GoogleImageCrawler,BaiduImageCrawler
from bing_image_downloader import downloader

from portalTypes import variant_from_dict

token = os.environ.get("token")
# bingCrawler = BaiduImageCrawler(storage={'root_dir': './images'})
def getPages(page):
	url = f"https://modernhyperindustry.com/site/api/v1/manage/store/products/variants?page={page}&size=100"
	payload = {}
	headers = {
	'Authorization': f'Bearer {token}'
	}

	response = requests.request("GET", url, headers=headers, data=payload)
	data = response.json()
	print(data)

	res = variant_from_dict(data)
	for variant  in res.variants:
		downloader.download(variant.title,limit=5,output_dir=f"./images/{variant.id}")
	if res.count * page<res.total:
		getPages(page+1)
getPages(1)
