from itemadapter import ItemAdapter
import psycopg2
import uuid
from datetime import datetime
import os

class PostgresPipeline:
    def open_spider(self, spider):
        dbname = os.getenv('DB_NAME', 'default_db')
        user = os.getenv('DB_USER', 'default_user')
        password = os.getenv('DB_PASSWORD', 'default_password')
        host = os.getenv('DB_HOST', 'localhost')
        port = os.getenv('DB_PORT', '5432')
        self.conn = psycopg2.connect(dbname=dbname, user=user, password=password, host=host, port=port)
        self.cursor = self.conn.cursor()
        self.job_id = str(uuid.uuid4())
    def process_item(self, item, spider):
        self.cursor.execute(
            "INSERT INTO scrap.items (text, author, tags, job_id, processed_at) VALUES (%s, %s, %s, %s, %s)",
            (item['text'], item['author'], ', '.join(item['tags']), self.job_id, datetime.now())
        )
        self.conn.commit()
        return item
    def close_spider(self, spider):
        self.cursor.close()
        self.conn.close()
