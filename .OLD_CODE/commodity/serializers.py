from rest_framework import serializers
from .models import Commodity
from .admin import CommodityAdmin


class CommoditySerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Commodity
        fields = list(CommodityAdmin.list_display)
