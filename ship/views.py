"""Module to render view templates"""
from django.http import HttpResponse
from django.template import loader
from rest_framework import generics
from rest_framework import mixins
from rest_framework.response import Response

from .models import Ship
from .serializers import ShipSerializer
from .admin import ShipAdmin
# ViewSets define the view behavior.


class ViewList(mixins.RetrieveModelMixin,
               generics.GenericAPIView):
    """:route: **list commodities**

    | lists ships

        [GET]: http://server:port/ship/list

    Returns:
        JSON: status_code: 200

    .. literalinclude:: write/list.json
        :language: JSON

    """
    queryset = Ship.objects.all()
    serializer_class = ShipSerializer

    # @universal.required_key(settings.FRONT_API_KEY)
    def get(self, request, *args, **kwargs):
        commodity = Ship.objects.all()
        serializer = ShipSerializer(commodity, many=True)
        return Response(serializer.data)


def index(request):
    data = Ship.objects.all()
    template = loader.get_template('table.html')
    context = {
        'data': data,
        'fields': list(ShipAdmin.list_display)
    }
    return HttpResponse(template.render(context, request))
