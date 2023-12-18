"module to make auto login into admin interface for guests"
from django.template.defaulttags import register
from django.template import loader
from django.http import HttpResponse
from django.contrib import auth
from django.contrib.auth import (
    authenticate,
    get_user_model,
)
from django.contrib.auth.models import Permission
from django.shortcuts import redirect
from django.conf import settings
from django.http import JsonResponse


def check_perm(user):
    "checks and adds view permissions for model_obj to user"
    permissions = Permission.objects.all()
    # s2 = [p.codename for p in s1]
    for perm in permissions:
        if 'view' not in perm.codename:
            continue

        if (perm.content_type.app_label not in settings.ADDED_APPS
                and not settings.DEBUG):
            continue

        permis = perm.name + '.' + perm.codename

        if not user.has_perm(permis):
            user.user_permissions.add(perm)
            user.save()


# Create your views here.


def login(request):
    "login to admin interface"

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


def get_index(request):
    template = loader.get_template('get_index.html')
    context = {'data': {}, 'fields': []}
    return HttpResponse(template.render(context, request))


def table(request):
    template = loader.get_template('table.html')
    context = {'data': {}, 'fields': []}
    return HttpResponse(template.render(context, request))


@register.filter
def get_item(dictionary, key):
    return getattr(dictionary, key)


def get_server(request):

    return JsonResponse(
        {"server_address": f'{request.scheme}://{request.get_host()}'})
