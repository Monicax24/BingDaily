import requests

CAT_IMG = open("./cat.jpg", "rb").read()
DOG_IMG = open("./dog.jpg", "rb").read()

HEIGHT = 1200
WIDTH = 630

def generate_random_img(height, width) -> bytes:
    res = requests.get(f"https://picsum.photos/{height}/{width}")
    return res.content

def upload_cat(url):
    res =  requests.put(url, data=CAT_IMG)
    return res.status_code == 200

def upload_dog(url):
    res = requests.put(url, data=DOG_IMG)
    return res.status_code == 200

def upload_random(url, height=HEIGHT, width=WIDTH):
    random_img = generate_random_img(height, width)
    res = requests.put(url, data=random_img)
    return res.status_code == 200
