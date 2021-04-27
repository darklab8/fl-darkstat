"""Module to render view templates"""
from django.http import HttpResponse
from django.template import loader
from rest_framework import generics
from rest_framework import mixins
from rest_framework.response import Response

from .models import Commodity
from .serializers import CommoditySerializer
from .admin import CommodityAdmin
# ViewSets define the view behavior.


class ViewList(mixins.RetrieveModelMixin,
               generics.GenericAPIView):
    """:route: **list commodities**

    | lists commodities

        [GET]: http://server:port/commodity/list

    Returns:
        JSON: status_code: 200

    .. literalinclude:: write/list.json
        :language: JSON

    """
    queryset = Commodity.objects.all()
    serializer_class = CommoditySerializer

    # @universal.required_key(settings.FRONT_API_KEY)
    def get(self, request, *args, **kwargs):
        commodity = Commodity.objects.all()
        serializer = CommoditySerializer(commodity, many=True)
        return Response(serializer.data)

        # return self.retrieve(request, *args, **kwargs)


def index(request):
    data = Commodity.objects.all()
    template = loader.get_template('table.html')
    context = {
        'data': data,
        'fields': list(CommodityAdmin.list_display)
    }
    return HttpResponse(template.render(context, request))
