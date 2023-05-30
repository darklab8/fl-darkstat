from django.urls import path
from rest_framework.urlpatterns import format_suffix_patterns
from rest_framework.decorators import api_view
from rest_framework.reverse import reverse
from rest_framework.response import Response
from . import views


@api_view(['GET'])
def commodity_api_route(request, format=None):
    return Response({
        'commodity get':
        reverse('commodity-get', request=request, format=format)
    })


# API endpoints
urlpatterns = format_suffix_patterns([
    path('', commodity_api_route, name='commodity-root'),
    path('list', views.ViewList.as_view(), name='commodity-get'),
    path('get_main', views.index, name='index'),
    path('get_one', views.get_one, name='get_one')
])

# urlpatterns = [
#     # ex: /polls/
#     path('', views.index, name='index'),
# ]
