"""fldarknet URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.1/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import include, path

from rest_framework.reverse import reverse
from rest_framework.decorators import api_view
from rest_framework.response import Response


# from rest_framework import routers

# from ship.admin import ShipViewSet
# from commodity.admin import CommodityViewSet

# Routers provide an easy way of automatically determining the URL conf.
# router = routers.DefaultRouter()
# router.register(r'ship', ShipViewSet)
# router.register(r'commodity', CommodityViewSet)


@api_view(['GET'])
def api_root(request, format=None):
    return Response({
        'commodity root': reverse(
            'commodity-root',
            request=request,
            format=format
        ),
        'ship root': reverse(
            'ship-root',
            request=request,
            format=format
        )
    })


urlpatterns = [
    path('', include('main.urls')),

    path('admin/login/', include('main.urls')),
    path('admin/logout/', include('main.urls')),
    path('admin/password_change/', include('main.urls')),

    path('admin/', admin.site.urls),
    # path('inbuilt_2/', admin.site.urls),
    path('api/', api_root),
    # path('api/', include(router.urls)),
    path('api-auth/',
         include('rest_framework.urls',
                 namespace='rest_framework')
         ),

    path('commodity/', include('commodity.urls')),
    path('ship/', include('ship.urls')),
]
