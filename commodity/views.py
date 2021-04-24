"""Module to render view templates"""
from django.http import HttpResponse
from django.template import loader

from .models import Commodity


def index(request):
    data = Commodity.objects.all()
    template = loader.get_template('commodity/index.html')
    context = {
        'data': data,
    }
    return HttpResponse(template.render(context, request))
