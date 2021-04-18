from core.celery import app

@app.task
def check():
 print('I am checking your stuff')