"module to render ships in admin interface"
from django.contrib import admin
from .models import Ship


class ShipAdmin(admin.ModelAdmin):
    """class to rewrite standard model view to
    what we need showing in admin interface"""
    # print(tuple([f.name for f in Commodity._meta.get_fields()]))
    list_display = (
        "nickname",  # str
        "info_name",  # str - infocardish real name
        "name",  # str - ini name
        "ship_class",  # int
        "typeof",  # str - type
        "hold_size",  # int - cargo hold
        "nanobot_limit",  # int
        "shield_battery_limit",  # int
        "capacity",  # int - powercore capacity
        "charge_rate",  # int - powercore charge
        "cruise_speed",  # int - engine
        "impulse_speed",  # int - engine
        "hit_pts",  # int - health points
        "id",
    )

    list_per_page = 1000


admin.site.register(Ship, ShipAdmin)
