from django.urls import path
from rest_framework.urlpatterns import format_suffix_patterns
from rest_framework.decorators import api_view
from rest_framework.reverse import reverse
from rest_framework.response import Response
from . import views


@api_view(['GET'])
def api_route(request, format=None):
    return Response({
        'ship get': reverse(
            'ship-get',
            request=request,
            format=format
        )
    })


# API endpoints
urlpatterns = format_suffix_patterns([
    path('', api_route,
         name='ship-root'),
    path('list',
         views.ViewList.as_view(),
         name='ship-get'),
])
