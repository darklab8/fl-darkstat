from django.urls import path

from . import views

urlpatterns = [
    # ex: /polls/5/
    path('', views.login, name='detail'),
]
