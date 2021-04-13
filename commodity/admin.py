""""Module to show model commodity in admin interface"""
from django.contrib import admin
from rest_framework import serializers, viewsets
from .models import Commodity


class CommodityAdmin(admin.ModelAdmin):
    """Rewriting standard model viewer with custom one to show all necessary commodity colums"""

    # print(tuple([f.name for f in Commodity._meta.get_fields()]))
    list_display = (
        # str
        "nickname",
        "loot_appearance",
        "pod_appearance",
        # float
        "volume",
        # int
        "ids_info",
        "units_per_container",
        "decay_per_second",
        "hit_pts",
        # SPECIAL
        "ids_name",
        "name",
    )

    list_per_page = 1000


admin.site.register(Commodity, CommodityAdmin)

# Serializers define the API representation.
class CommoditySerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Commodity
        fields = list(CommodityAdmin.list_display)


# ViewSets define the view behavior.
class CommodityViewSet(viewsets.ModelViewSet):
    queryset = Commodity.objects.all()
    serializer_class = CommoditySerializer
