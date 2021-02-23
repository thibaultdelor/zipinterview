from datetime import datetime
import json
import os
import re

from elasticsearch import Elasticsearch

es = Elasticsearch([os.getenv('ELASTICSEARCH_URL', "http://localhost:9200/")])
iso_format = '%Y-%m-%dT%H:%M:%SZ'
time_window_in_sec = int(os.getenv('TIME_WINDOW', "180"))


def handler(event, context):

    # event date in milliseconds minus the time window
    date_millis = int(datetime.strptime(
        event["time"], iso_format).timestamp() * 1000) - time_window_in_sec * 1000
    count = es.search(size=0, body={
        "query": {
            "bool": {
                "filter": [
                    {"range": {"@timestamp": {"gte": date_millis}}}
                ]
            }
        },
        "aggs": {
            "byuser": {
                "cardinality": {
                    "field": "user_id"
                }
            }
        }
    }
    )
    return {
        'total': count['hits']['total']['value'],
        'countUsers': count['aggregations']['byuser']['value']
    }


if __name__ == "__main__":
    # For local dev purpose, try a sample event as described here:
    # https://docs.aws.amazon.com/AmazonCloudWatch/latest/events/RunLambdaSchedule.html
    print(handler({
        "version": "0",
        "id": "53dc4d37-cffa-4f76-80c9-8b7d4a4d2eaa",
        "detail-type": "Scheduled Event",
        "source": "aws.events",
        "account": "123456789012",
        "time": datetime.now().strftime(iso_format),
        "region": "us-east-1",
        "resources": [
            "arn:aws:events:us-east-1:123456789012:rule/my-scheduled-rule"
        ],
        "detail": {}
    }, {}))
