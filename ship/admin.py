"module to render ships in admin interface"
from django.contrib import admin
from rest_framework import serializers, viewsets
from .models import Ship


class ShipAdmin(admin.ModelAdmin):
    """class to rewrite standard model view to
    what we need showing in admin interface"""
    # print(tuple([f.name for f in Commodity._meta.get_fields()]))
    list_display = (
        # str
        "nickname",
        # infocardish real name
        "info_name",
        # SPECIAL
        "ids_name",
        "name",
        "ship_class",
        "hold_size",

        # powercore
        "capacity",
        "charge_rate",
        # engine
        "cruise_speed",
        "impulse_speed",
        # float
        "mass",
        "linear_drag",
        "camera_angular_acceleration",
        "camera_turn_look_ahead_slerp_amount",
        "nudge_force",
        "explosion_resistance",
        # int
        "ship_class",
        "hold_size",
        "hit_pts",
        "nanobot_limit",
        "shield_battery_limit",
        "ids_info",
        "max_bank_angle",
        "camera_horizontal_turn_angle",
        "camera_vertical_turn_up_angle",
        "camera_vertical_turn_down_angle",
        # "ids_info1",
        # "ids_info2",
        # "ids_info3",
        "strafe_force",
        "strafe_power_usage",

        # type
        "typeof",


    )

    # different original name: type

    # tuple([f.name for f in Commodity._meta.get_fields()])
    list_per_page = 1000
    # search_fields = ['address']
    # list_filter = ['visible']
    # actions = [make_saved, make_loaded, make_visible, make_invisible]
    # inlines = [ChoiceInline, ]


admin.site.register(Ship, ShipAdmin)


# Serializers define the API representation.
class ShipSerializer(serializers.HyperlinkedModelSerializer):
    class Meta:
        model = Ship
        fields = list(ShipAdmin.list_display)


# ViewSets define the view behavior.
class ShipViewSet(viewsets.ModelViewSet):
    queryset = Ship.objects.all()
    serializer_class = ShipSerializer
