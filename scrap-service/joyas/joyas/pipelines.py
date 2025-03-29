from itemadapter import ItemAdapter
import csv
class JoyasPipeline:
    def open_spider(self, spider):
        self.file = open('joyas.csv', 'w', newline='', encoding='utf-8')
        self.writer = csv.writer(self.file)
        self.writer.writerow(['Text', 'Author', 'Tags'])
    def process_item(self, item, spider):
        self.writer.writerow([item['text'], item['author'], ", ".join(item['tags'])])
        return item
    def close_spider(self, spider):
        self.file.close()