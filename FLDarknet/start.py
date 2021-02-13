"""
This script runs the application using a development server.
It contains the definition of routes and views for the application.
"""
import flint as fl
from app import create_app, db, cli
#from app.models import User, Post, Message, Notification, Task

app = create_app()
cli.register(app)

# @app.shell_context_processor
# def make_shell_context():
    # return {'db': db, 'User': User, 'Post': Post, 'Message': Message,
            # 'Notification': Notification, 'Task': Task}


# Make the WSGI interface available at the top level so wfastcgi can get it.
wsgi_app = app.wsgi_app


# @app.route('/')
# def hello():
    # """Renders a sample page."""
    # return "Hello World!"

if __name__ == '__main__':
    import os
    HOST = os.environ.get('SERVER_HOST', '0.0.0.0')
    try:
        PORT = int(os.environ.get('SERVER_PORT', '5000'))
    except ValueError:
        PORT = 5555
    app.run(HOST, PORT)
