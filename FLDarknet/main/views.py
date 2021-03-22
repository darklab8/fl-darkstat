from django.shortcuts import render
from django.contrib import auth
from django.contrib.auth import authenticate
from django.contrib.auth.backends import ModelBackend
from django.contrib.auth.models import User
from django.shortcuts import redirect

def check_perm(user, Obj):
    from django.contrib.auth.models import Permission
    from django.contrib.contenttypes.models import ContentType
    content_type = ContentType.objects.get_for_model(Obj, for_concrete_model=False)
    s1 = Permission.objects.filter(content_type=content_type)
    #s2 = [p.codename for p in s1]
    for perm in s1:
        if 'view' in perm.codename:
            permis = content_type.name + '.' + perm.codename
            if (not user.has_perm(permis)):
                user.user_permissions.add(perm)
                user.save()

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

    from commodities.models import Commodity
    from ship.models import Ship
    check_perm(user, Commodity)
    check_perm(user, Ship)

    user = authenticate(username='guest', password='guest')
    request.user = user
    auth.login(request, user)

    return redirect('/admin')