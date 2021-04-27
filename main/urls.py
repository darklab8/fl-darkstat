"""Sneaky autologin for all incoming users"""
from django.urls import path

from . import views

urlpatterns = [
    # ex: /polls/5/
    path('login/', views.login, name='detail'),
    path('', views.index, name='detail'),
]
