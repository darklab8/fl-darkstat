from django.shortcuts import render
from django.contrib import auth
from django.contrib.auth import authenticate
from django.contrib.auth.backends import ModelBackend
from django.contrib.auth.models import User
from django.shortcuts import redirect

# Create your views here.
def login(request):
    username = 'guest'
    try:
        user = User.objects.get(username=username)
    except User.DoesNotExist:
        user = User(username=username)
        user.set_password('guest')
        user.is_staff = True
        user.save()

    user = authenticate(username='guest', password='guest')
    request.user = user
    auth.login(request, user)
    return redirect('/admin')