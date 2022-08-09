from flask import Flask

app = Flask(__name__, static_folder='./out/', static_url_path='')


@app.route('/')
def index():
    return app.send_static_file('index.html')

@app.route('/<type>')
def type(type):
    return app.send_static_file('[type].html')

@app.route('/<type>/<entity>')
def entity(type, entity):
    return app.send_static_file('[type]/[entity].html')

if __name__ == '__main__':
    app.run(threaded=True, port=5000)