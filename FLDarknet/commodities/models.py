from django.db import models

# Create your models here.
class Commodity(models.Model):
    class Meta:
        verbose_name_plural = "commodities"
    
    
    nickname = models.CharField(max_length=50, db_index=True, blank=True, null=True)

    ids_name = models.IntegerField(db_index=True, blank=True, null=True)
    ids_info = models.IntegerField(db_index=True, blank=True, null=True)

    units_per_container = models.IntegerField(db_index=True, blank=True, null=True)
    decay_per_second = models.IntegerField(blank=True, null=True)
    hit_pts = models.IntegerField(blank=True, null=True)

    volume = models.FloatField(blank=True, null=True)

    loot_appearance = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    pod_appearance = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    
    name = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    infocard = models.CharField(max_length=500, db_index=True, blank=True, null=True)