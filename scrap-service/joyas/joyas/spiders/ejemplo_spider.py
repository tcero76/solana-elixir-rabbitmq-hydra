import scrapy
from bs4 import BeautifulSoup
from joyas.items import JoyasItem

class EjemploSpider(scrapy.Spider):
    name = "ejemplo"
    start_urls = ["https://quotes.toscrape.com/"]

    def parse(self, response):
        soup = BeautifulSoup(response.text, "html.parser")
        for quote in soup.select("div.quote"):
            item = JoyasItem()
            item['text'] = quote.select_one('span.text').get_text(strip=True)
            item['author'] = quote.select_one('small.author').get_text(strip=True)
            item['tags'] =[tag.get_text(strip=True) for tag in quote.select('div.tags a.tag')]
            yield item