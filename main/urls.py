"""Sneaky autologin for all incoming users"""
from django.urls import path

from . import views

urlpatterns = [
    # ex: /polls/5/
    path('login/', views.login, name='detail'),
    path('get_index', views.get_index, name='detail'),
    path('', views.table, name='detail'),
    path('server', views.get_server, name='detail'),
]
