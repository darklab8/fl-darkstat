"module to make auto login into admin interface for guests"
from django.contrib import auth
from django.contrib.auth import (
    authenticate,
    get_user_model,
)
from django.contrib.auth.models import Permission
from django.shortcuts import redirect
from django.core import management


def check_perm(user):
    "checks and adds view permissions for model_obj to user"
    permissions = Permission.objects.all()
    # s2 = [p.codename for p in s1]
    for perm in permissions:
        if 'view' not in perm.codename:
            continue
        permis = perm.name + '.' + perm.codename
        if not user.has_perm(permis):
            user.user_permissions.add(perm)
            user.save()

# Create your views here.


def login(request):
    "login to admin interface"

    management.call_command('migrate')

    username = 'guest'
    user_model = get_user_model()
    try:
        user = user_model.objects.get(username=username)
    except user_model.DoesNotExist:
        user = user_model(username=username)
        user.set_password('guest')
        user.is_staff = True
        user.save()

    check_perm(user)

    user = authenticate(username='guest', password='guest')
    request.user = user
    auth.login(request, user)

    return redirect('/admin')
