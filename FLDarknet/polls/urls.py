from django.urls import path

from . import views

urlpatterns = [
    path('', views.thing, name='not_index'),
]