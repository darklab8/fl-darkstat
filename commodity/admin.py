""""Module to show model commodity in admin interface"""
from django.contrib import admin
from rest_framework import serializers, viewsets
from .models import Commodity


class CommodityAdmin(admin.ModelAdmin):
    """Rewriting standard model viewer with custom one to show all necessary commodity colums"""

    # print(tuple([f.name for f in Commodity._meta.get_fields()]))
    list_display = ('name', 'nickname', 'ids_name', 'ids_info', 'units_per_container',
                    'decay_per_second', 'hit_pts', 'volume', 'loot_appearance', 'pod_appearance')

    #tuple([f.name for f in Commodity._meta.get_fields()])
    list_per_page = 1000
    # search_fields = ['address']
    # list_filter = ['visible']
    # actions = [make_saved, make_loaded, make_visible, make_invisible]
    # inlines = [ChoiceInline, ]


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

# Register your models here.

# nickname = models.CharField(max_length=50, db_index=True, blank=True, null=True)

#     ids_name = models.IntegerField(db_index=True, blank=True, null=True)
#     ids_info = models.IntegerField(db_index=True, blank=True, null=True)

#     units_per_container = models.IntegerField(db_index=True, blank=True, null=True)
#     decay_per_second = models.IntegerField(blank=True, null=True)
#     hit_pts = models.IntegerField(blank=True, null=True)

#     volume = models.FloatField(blank=True, null=True)

#     loot_appearance = models.CharField(max_length=50, db_index=True, blank=True, null=True)
#     pod_appearance = models.CharField(max_length=50, db_index=True, blank=True, null=True)

#     name = models.CharField(max_length=50, db_index=True, blank=True, null=True)
#     infocard = models.CharField(max_length=500, db_index=True, blank=True, null=True)
