from rest_framework import serializers
from .models import Ship
from .admin import ShipAdmin
# Serializers define the API representation.


class ShipSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Ship
        fields = list(ShipAdmin.list_display)
