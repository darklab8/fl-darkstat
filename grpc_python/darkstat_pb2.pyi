from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Empty(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetInfocardsInput(_message.Message):
    __slots__ = ("nicknames",)
    NICKNAMES_FIELD_NUMBER: _ClassVar[int]
    nicknames: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, nicknames: _Optional[_Iterable[str]] = ...) -> None: ...

class GetInfocardsReply(_message.Message):
    __slots__ = ("answers",)
    ANSWERS_FIELD_NUMBER: _ClassVar[int]
    answers: _containers.RepeatedCompositeFieldContainer[GetInfocardAnswer]
    def __init__(self, answers: _Optional[_Iterable[_Union[GetInfocardAnswer, _Mapping]]] = ...) -> None: ...

class GetInfocardAnswer(_message.Message):
    __slots__ = ("query", "infocard", "error")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    INFOCARD_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    query: str
    infocard: Infocard
    error: str
    def __init__(self, query: _Optional[str] = ..., infocard: _Optional[_Union[Infocard, _Mapping]] = ..., error: _Optional[str] = ...) -> None: ...

class Infocard(_message.Message):
    __slots__ = ("lines",)
    LINES_FIELD_NUMBER: _ClassVar[int]
    lines: _containers.RepeatedCompositeFieldContainer[InfocardLine]
    def __init__(self, lines: _Optional[_Iterable[_Union[InfocardLine, _Mapping]]] = ...) -> None: ...

class InfocardLine(_message.Message):
    __slots__ = ("phrases",)
    PHRASES_FIELD_NUMBER: _ClassVar[int]
    phrases: _containers.RepeatedCompositeFieldContainer[InfocardPhrase]
    def __init__(self, phrases: _Optional[_Iterable[_Union[InfocardPhrase, _Mapping]]] = ...) -> None: ...

class InfocardPhrase(_message.Message):
    __slots__ = ("phrase", "link", "bold")
    PHRASE_FIELD_NUMBER: _ClassVar[int]
    LINK_FIELD_NUMBER: _ClassVar[int]
    BOLD_FIELD_NUMBER: _ClassVar[int]
    phrase: str
    link: str
    bold: bool
    def __init__(self, phrase: _Optional[str] = ..., link: _Optional[str] = ..., bold: bool = ...) -> None: ...

class HealthReply(_message.Message):
    __slots__ = ("is_healthy",)
    IS_HEALTHY_FIELD_NUMBER: _ClassVar[int]
    is_healthy: bool
    def __init__(self, is_healthy: bool = ...) -> None: ...

class GetEquipmentInput(_message.Message):
    __slots__ = ("include_market_goods", "include_tech_compat", "filter_to_useful", "filter_nicknames")
    INCLUDE_MARKET_GOODS_FIELD_NUMBER: _ClassVar[int]
    INCLUDE_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    FILTER_TO_USEFUL_FIELD_NUMBER: _ClassVar[int]
    FILTER_NICKNAMES_FIELD_NUMBER: _ClassVar[int]
    include_market_goods: bool
    include_tech_compat: bool
    filter_to_useful: bool
    filter_nicknames: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, include_market_goods: bool = ..., include_tech_compat: bool = ..., filter_to_useful: bool = ..., filter_nicknames: _Optional[_Iterable[str]] = ...) -> None: ...

class GetGunsInput(_message.Message):
    __slots__ = ("include_market_goods", "include_tech_compat", "filter_to_useful", "include_damage_bonuses", "filter_nicknames")
    INCLUDE_MARKET_GOODS_FIELD_NUMBER: _ClassVar[int]
    INCLUDE_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    FILTER_TO_USEFUL_FIELD_NUMBER: _ClassVar[int]
    INCLUDE_DAMAGE_BONUSES_FIELD_NUMBER: _ClassVar[int]
    FILTER_NICKNAMES_FIELD_NUMBER: _ClassVar[int]
    include_market_goods: bool
    include_tech_compat: bool
    filter_to_useful: bool
    include_damage_bonuses: bool
    filter_nicknames: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, include_market_goods: bool = ..., include_tech_compat: bool = ..., filter_to_useful: bool = ..., include_damage_bonuses: bool = ..., filter_nicknames: _Optional[_Iterable[str]] = ...) -> None: ...

class GetBasesInput(_message.Message):
    __slots__ = ("include_market_goods", "filter_to_useful", "filter_nicknames", "filter_market_good_category")
    INCLUDE_MARKET_GOODS_FIELD_NUMBER: _ClassVar[int]
    FILTER_TO_USEFUL_FIELD_NUMBER: _ClassVar[int]
    FILTER_NICKNAMES_FIELD_NUMBER: _ClassVar[int]
    FILTER_MARKET_GOOD_CATEGORY_FIELD_NUMBER: _ClassVar[int]
    include_market_goods: bool
    filter_to_useful: bool
    filter_nicknames: _containers.RepeatedScalarFieldContainer[str]
    filter_market_good_category: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, include_market_goods: bool = ..., filter_to_useful: bool = ..., filter_nicknames: _Optional[_Iterable[str]] = ..., filter_market_good_category: _Optional[_Iterable[str]] = ...) -> None: ...

class GetTractorsInput(_message.Message):
    __slots__ = ("include_market_goods", "filter_to_useful", "filter_nicknames")
    INCLUDE_MARKET_GOODS_FIELD_NUMBER: _ClassVar[int]
    FILTER_TO_USEFUL_FIELD_NUMBER: _ClassVar[int]
    FILTER_NICKNAMES_FIELD_NUMBER: _ClassVar[int]
    include_market_goods: bool
    filter_to_useful: bool
    filter_nicknames: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, include_market_goods: bool = ..., filter_to_useful: bool = ..., filter_nicknames: _Optional[_Iterable[str]] = ...) -> None: ...

class GetBasesReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Base]
    def __init__(self, items: _Optional[_Iterable[_Union[Base, _Mapping]]] = ...) -> None: ...

class Base(_message.Message):
    __slots__ = ("name", "archetypes", "nickname", "faction_name", "system", "system_nickname", "region", "strid_name", "infocard_id", "file", "bgcs_base_run_by", "pos", "sector_coord", "is_transport_unreachable", "reachable", "is_pob", "market_goods_per_nick")
    class MarketGoodsPerNickEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    ARCHETYPES_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    FACTION_NAME_FIELD_NUMBER: _ClassVar[int]
    SYSTEM_FIELD_NUMBER: _ClassVar[int]
    SYSTEM_NICKNAME_FIELD_NUMBER: _ClassVar[int]
    REGION_FIELD_NUMBER: _ClassVar[int]
    STRID_NAME_FIELD_NUMBER: _ClassVar[int]
    INFOCARD_ID_FIELD_NUMBER: _ClassVar[int]
    FILE_FIELD_NUMBER: _ClassVar[int]
    BGCS_BASE_RUN_BY_FIELD_NUMBER: _ClassVar[int]
    POS_FIELD_NUMBER: _ClassVar[int]
    SECTOR_COORD_FIELD_NUMBER: _ClassVar[int]
    IS_TRANSPORT_UNREACHABLE_FIELD_NUMBER: _ClassVar[int]
    REACHABLE_FIELD_NUMBER: _ClassVar[int]
    IS_POB_FIELD_NUMBER: _ClassVar[int]
    MARKET_GOODS_PER_NICK_FIELD_NUMBER: _ClassVar[int]
    name: str
    archetypes: _containers.RepeatedScalarFieldContainer[str]
    nickname: str
    faction_name: str
    system: str
    system_nickname: str
    region: str
    strid_name: int
    infocard_id: int
    file: str
    bgcs_base_run_by: str
    pos: Pos
    sector_coord: str
    is_transport_unreachable: bool
    reachable: bool
    is_pob: bool
    market_goods_per_nick: _containers.MessageMap[str, MarketGood]
    def __init__(self, name: _Optional[str] = ..., archetypes: _Optional[_Iterable[str]] = ..., nickname: _Optional[str] = ..., faction_name: _Optional[str] = ..., system: _Optional[str] = ..., system_nickname: _Optional[str] = ..., region: _Optional[str] = ..., strid_name: _Optional[int] = ..., infocard_id: _Optional[int] = ..., file: _Optional[str] = ..., bgcs_base_run_by: _Optional[str] = ..., pos: _Optional[_Union[Pos, _Mapping]] = ..., sector_coord: _Optional[str] = ..., is_transport_unreachable: bool = ..., reachable: bool = ..., is_pob: bool = ..., market_goods_per_nick: _Optional[_Mapping[str, MarketGood]] = ...) -> None: ...

class MiningInfo(_message.Message):
    __slots__ = ("dynamic_loot_min", "dynamic_loot_max", "dynamic_loot_difficulty", "mined_good")
    DYNAMIC_LOOT_MIN_FIELD_NUMBER: _ClassVar[int]
    DYNAMIC_LOOT_MAX_FIELD_NUMBER: _ClassVar[int]
    DYNAMIC_LOOT_DIFFICULTY_FIELD_NUMBER: _ClassVar[int]
    MINED_GOOD_FIELD_NUMBER: _ClassVar[int]
    dynamic_loot_min: int
    dynamic_loot_max: int
    dynamic_loot_difficulty: int
    mined_good: MarketGood
    def __init__(self, dynamic_loot_min: _Optional[int] = ..., dynamic_loot_max: _Optional[int] = ..., dynamic_loot_difficulty: _Optional[int] = ..., mined_good: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...

class MarketGood(_message.Message):
    __slots__ = ("nickname", "ship_nickname", "name", "price_base", "hp_type", "category", "level_required", "rep_required", "price_base_buys_for", "price_base_sells_for", "volume", "ship_class", "base_sells", "is_server_side_override", "not_buyable", "is_transport_unreachable", "base_info")
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    SHIP_NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_BASE_FIELD_NUMBER: _ClassVar[int]
    HP_TYPE_FIELD_NUMBER: _ClassVar[int]
    CATEGORY_FIELD_NUMBER: _ClassVar[int]
    LEVEL_REQUIRED_FIELD_NUMBER: _ClassVar[int]
    REP_REQUIRED_FIELD_NUMBER: _ClassVar[int]
    PRICE_BASE_BUYS_FOR_FIELD_NUMBER: _ClassVar[int]
    PRICE_BASE_SELLS_FOR_FIELD_NUMBER: _ClassVar[int]
    VOLUME_FIELD_NUMBER: _ClassVar[int]
    SHIP_CLASS_FIELD_NUMBER: _ClassVar[int]
    BASE_SELLS_FIELD_NUMBER: _ClassVar[int]
    IS_SERVER_SIDE_OVERRIDE_FIELD_NUMBER: _ClassVar[int]
    NOT_BUYABLE_FIELD_NUMBER: _ClassVar[int]
    IS_TRANSPORT_UNREACHABLE_FIELD_NUMBER: _ClassVar[int]
    BASE_INFO_FIELD_NUMBER: _ClassVar[int]
    nickname: str
    ship_nickname: str
    name: str
    price_base: int
    hp_type: str
    category: str
    level_required: int
    rep_required: float
    price_base_buys_for: int
    price_base_sells_for: int
    volume: float
    ship_class: int
    base_sells: bool
    is_server_side_override: bool
    not_buyable: bool
    is_transport_unreachable: bool
    base_info: BaseInfo
    def __init__(self, nickname: _Optional[str] = ..., ship_nickname: _Optional[str] = ..., name: _Optional[str] = ..., price_base: _Optional[int] = ..., hp_type: _Optional[str] = ..., category: _Optional[str] = ..., level_required: _Optional[int] = ..., rep_required: _Optional[float] = ..., price_base_buys_for: _Optional[int] = ..., price_base_sells_for: _Optional[int] = ..., volume: _Optional[float] = ..., ship_class: _Optional[int] = ..., base_sells: bool = ..., is_server_side_override: bool = ..., not_buyable: bool = ..., is_transport_unreachable: bool = ..., base_info: _Optional[_Union[BaseInfo, _Mapping]] = ...) -> None: ...

class BaseInfo(_message.Message):
    __slots__ = ("base_nickname", "base_name", "system_name", "region", "faction_name", "base_pos", "sector_coord")
    BASE_NICKNAME_FIELD_NUMBER: _ClassVar[int]
    BASE_NAME_FIELD_NUMBER: _ClassVar[int]
    SYSTEM_NAME_FIELD_NUMBER: _ClassVar[int]
    REGION_FIELD_NUMBER: _ClassVar[int]
    FACTION_NAME_FIELD_NUMBER: _ClassVar[int]
    BASE_POS_FIELD_NUMBER: _ClassVar[int]
    SECTOR_COORD_FIELD_NUMBER: _ClassVar[int]
    base_nickname: str
    base_name: str
    system_name: str
    region: str
    faction_name: str
    base_pos: Pos
    sector_coord: str
    def __init__(self, base_nickname: _Optional[str] = ..., base_name: _Optional[str] = ..., system_name: _Optional[str] = ..., region: _Optional[str] = ..., faction_name: _Optional[str] = ..., base_pos: _Optional[_Union[Pos, _Mapping]] = ..., sector_coord: _Optional[str] = ...) -> None: ...

class Pos(_message.Message):
    __slots__ = ("x", "y", "z")
    X_FIELD_NUMBER: _ClassVar[int]
    Y_FIELD_NUMBER: _ClassVar[int]
    Z_FIELD_NUMBER: _ClassVar[int]
    x: float
    y: float
    z: float
    def __init__(self, x: _Optional[float] = ..., y: _Optional[float] = ..., z: _Optional[float] = ...) -> None: ...

class GetCommoditiesInput(_message.Message):
    __slots__ = ("include_market_goods", "filter_to_useful", "filter_nicknames")
    INCLUDE_MARKET_GOODS_FIELD_NUMBER: _ClassVar[int]
    FILTER_TO_USEFUL_FIELD_NUMBER: _ClassVar[int]
    FILTER_NICKNAMES_FIELD_NUMBER: _ClassVar[int]
    include_market_goods: bool
    filter_to_useful: bool
    filter_nicknames: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, include_market_goods: bool = ..., filter_to_useful: bool = ..., filter_nicknames: _Optional[_Iterable[str]] = ...) -> None: ...

class GetCommoditiesReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Commodity]
    def __init__(self, items: _Optional[_Iterable[_Union[Commodity, _Mapping]]] = ...) -> None: ...

class Commodity(_message.Message):
    __slots__ = ("nickname", "price_base", "name", "combinable", "volume", "ship_class", "name_id", "infocard_id", "bases", "price_best_base_buys_for", "price_best_base_sells_for", "proffit_margin", "mass")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_BASE_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    COMBINABLE_FIELD_NUMBER: _ClassVar[int]
    VOLUME_FIELD_NUMBER: _ClassVar[int]
    SHIP_CLASS_FIELD_NUMBER: _ClassVar[int]
    NAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFOCARD_ID_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    PRICE_BEST_BASE_BUYS_FOR_FIELD_NUMBER: _ClassVar[int]
    PRICE_BEST_BASE_SELLS_FOR_FIELD_NUMBER: _ClassVar[int]
    PROFFIT_MARGIN_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    nickname: str
    price_base: int
    name: str
    combinable: bool
    volume: float
    ship_class: int
    name_id: int
    infocard_id: int
    bases: _containers.MessageMap[str, MarketGood]
    price_best_base_buys_for: int
    price_best_base_sells_for: int
    proffit_margin: int
    mass: float
    def __init__(self, nickname: _Optional[str] = ..., price_base: _Optional[int] = ..., name: _Optional[str] = ..., combinable: bool = ..., volume: _Optional[float] = ..., ship_class: _Optional[int] = ..., name_id: _Optional[int] = ..., infocard_id: _Optional[int] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., price_best_base_buys_for: _Optional[int] = ..., price_best_base_sells_for: _Optional[int] = ..., proffit_margin: _Optional[int] = ..., mass: _Optional[float] = ...) -> None: ...

class GetAmmoReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Ammo]
    def __init__(self, items: _Optional[_Iterable[_Union[Ammo, _Mapping]]] = ...) -> None: ...

class Ammo(_message.Message):
    __slots__ = ("name", "price", "hit_pts", "volume", "munition_lifetime", "nickname", "name_id", "info_id", "seeker_type", "seeker_range", "seeker_fov_deg", "bases", "discovery_tech_compat", "ammo_limit", "mass")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    HIT_PTS_FIELD_NUMBER: _ClassVar[int]
    VOLUME_FIELD_NUMBER: _ClassVar[int]
    MUNITION_LIFETIME_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFO_ID_FIELD_NUMBER: _ClassVar[int]
    SEEKER_TYPE_FIELD_NUMBER: _ClassVar[int]
    SEEKER_RANGE_FIELD_NUMBER: _ClassVar[int]
    SEEKER_FOV_DEG_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    AMMO_LIMIT_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    name: str
    price: int
    hit_pts: int
    volume: float
    munition_lifetime: float
    nickname: str
    name_id: int
    info_id: int
    seeker_type: str
    seeker_range: int
    seeker_fov_deg: int
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    ammo_limit: AmmoLimit
    mass: float
    def __init__(self, name: _Optional[str] = ..., price: _Optional[int] = ..., hit_pts: _Optional[int] = ..., volume: _Optional[float] = ..., munition_lifetime: _Optional[float] = ..., nickname: _Optional[str] = ..., name_id: _Optional[int] = ..., info_id: _Optional[int] = ..., seeker_type: _Optional[str] = ..., seeker_range: _Optional[int] = ..., seeker_fov_deg: _Optional[int] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ..., ammo_limit: _Optional[_Union[AmmoLimit, _Mapping]] = ..., mass: _Optional[float] = ...) -> None: ...

class DiscoveryTechCompat(_message.Message):
    __slots__ = ("techcompat_by_id", "tech_cell")
    class TechcompatByIdEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: float
        def __init__(self, key: _Optional[str] = ..., value: _Optional[float] = ...) -> None: ...
    TECHCOMPAT_BY_ID_FIELD_NUMBER: _ClassVar[int]
    TECH_CELL_FIELD_NUMBER: _ClassVar[int]
    techcompat_by_id: _containers.ScalarMap[str, float]
    tech_cell: str
    def __init__(self, techcompat_by_id: _Optional[_Mapping[str, float]] = ..., tech_cell: _Optional[str] = ...) -> None: ...

class TechCompatAnswer(_message.Message):
    __slots__ = ("tech_compat", "error", "nickname")
    TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    tech_compat: DiscoveryTechCompat
    error: str
    nickname: str
    def __init__(self, tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ..., error: _Optional[str] = ..., nickname: _Optional[str] = ...) -> None: ...

class GetTechCompatInput(_message.Message):
    __slots__ = ("nicknames",)
    NICKNAMES_FIELD_NUMBER: _ClassVar[int]
    nicknames: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, nicknames: _Optional[_Iterable[str]] = ...) -> None: ...

class GetTechCompatReply(_message.Message):
    __slots__ = ("answers",)
    ANSWERS_FIELD_NUMBER: _ClassVar[int]
    answers: _containers.RepeatedCompositeFieldContainer[TechCompatAnswer]
    def __init__(self, answers: _Optional[_Iterable[_Union[TechCompatAnswer, _Mapping]]] = ...) -> None: ...

class GetCounterMeasuresReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[CounterMeasure]
    def __init__(self, items: _Optional[_Iterable[_Union[CounterMeasure, _Mapping]]] = ...) -> None: ...

class CounterMeasure(_message.Message):
    __slots__ = ("name", "price", "hit_pts", "ai_range", "lifetime", "range", "diversion_pctg", "lootable", "nickname", "name_id", "info_id", "bases", "discovery_tech_compat", "ammo_limit", "mass")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    HIT_PTS_FIELD_NUMBER: _ClassVar[int]
    AI_RANGE_FIELD_NUMBER: _ClassVar[int]
    LIFETIME_FIELD_NUMBER: _ClassVar[int]
    RANGE_FIELD_NUMBER: _ClassVar[int]
    DIVERSION_PCTG_FIELD_NUMBER: _ClassVar[int]
    LOOTABLE_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFO_ID_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    AMMO_LIMIT_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    name: str
    price: int
    hit_pts: int
    ai_range: int
    lifetime: int
    range: int
    diversion_pctg: int
    lootable: bool
    nickname: str
    name_id: int
    info_id: int
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    ammo_limit: AmmoLimit
    mass: float
    def __init__(self, name: _Optional[str] = ..., price: _Optional[int] = ..., hit_pts: _Optional[int] = ..., ai_range: _Optional[int] = ..., lifetime: _Optional[int] = ..., range: _Optional[int] = ..., diversion_pctg: _Optional[int] = ..., lootable: bool = ..., nickname: _Optional[str] = ..., name_id: _Optional[int] = ..., info_id: _Optional[int] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ..., ammo_limit: _Optional[_Union[AmmoLimit, _Mapping]] = ..., mass: _Optional[float] = ...) -> None: ...

class GetEnginesReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Engine]
    def __init__(self, items: _Optional[_Iterable[_Union[Engine, _Mapping]]] = ...) -> None: ...

class Engine(_message.Message):
    __slots__ = ("name", "price", "cruise_speed", "cruise_charge_time", "linear_drag", "max_force", "reverse_fraction", "impulse_speed", "hp_type", "flame_effect", "trail_effect", "nickname", "name_id", "info_id", "bases", "discovery_tech_compat", "mass")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    CRUISE_SPEED_FIELD_NUMBER: _ClassVar[int]
    CRUISE_CHARGE_TIME_FIELD_NUMBER: _ClassVar[int]
    LINEAR_DRAG_FIELD_NUMBER: _ClassVar[int]
    MAX_FORCE_FIELD_NUMBER: _ClassVar[int]
    REVERSE_FRACTION_FIELD_NUMBER: _ClassVar[int]
    IMPULSE_SPEED_FIELD_NUMBER: _ClassVar[int]
    HP_TYPE_FIELD_NUMBER: _ClassVar[int]
    FLAME_EFFECT_FIELD_NUMBER: _ClassVar[int]
    TRAIL_EFFECT_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFO_ID_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    name: str
    price: int
    cruise_speed: int
    cruise_charge_time: int
    linear_drag: int
    max_force: int
    reverse_fraction: float
    impulse_speed: float
    hp_type: str
    flame_effect: str
    trail_effect: str
    nickname: str
    name_id: int
    info_id: int
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    mass: float
    def __init__(self, name: _Optional[str] = ..., price: _Optional[int] = ..., cruise_speed: _Optional[int] = ..., cruise_charge_time: _Optional[int] = ..., linear_drag: _Optional[int] = ..., max_force: _Optional[int] = ..., reverse_fraction: _Optional[float] = ..., impulse_speed: _Optional[float] = ..., hp_type: _Optional[str] = ..., flame_effect: _Optional[str] = ..., trail_effect: _Optional[str] = ..., nickname: _Optional[str] = ..., name_id: _Optional[int] = ..., info_id: _Optional[int] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ..., mass: _Optional[float] = ...) -> None: ...

class GetFactionsInput(_message.Message):
    __slots__ = ("include_reputations", "include_bribes", "filter_to_useful")
    INCLUDE_REPUTATIONS_FIELD_NUMBER: _ClassVar[int]
    INCLUDE_BRIBES_FIELD_NUMBER: _ClassVar[int]
    FILTER_TO_USEFUL_FIELD_NUMBER: _ClassVar[int]
    include_reputations: bool
    include_bribes: bool
    filter_to_useful: bool
    def __init__(self, include_reputations: bool = ..., include_bribes: bool = ..., filter_to_useful: bool = ...) -> None: ...

class GetFactionsReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Faction]
    def __init__(self, items: _Optional[_Iterable[_Union[Faction, _Mapping]]] = ...) -> None: ...

class Faction(_message.Message):
    __slots__ = ("name", "short_name", "nickname", "object_destruction", "mission_success", "mission_failure", "mission_abort", "infoname_id", "infocard_id", "reputations", "bribes")
    NAME_FIELD_NUMBER: _ClassVar[int]
    SHORT_NAME_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    OBJECT_DESTRUCTION_FIELD_NUMBER: _ClassVar[int]
    MISSION_SUCCESS_FIELD_NUMBER: _ClassVar[int]
    MISSION_FAILURE_FIELD_NUMBER: _ClassVar[int]
    MISSION_ABORT_FIELD_NUMBER: _ClassVar[int]
    INFONAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFOCARD_ID_FIELD_NUMBER: _ClassVar[int]
    REPUTATIONS_FIELD_NUMBER: _ClassVar[int]
    BRIBES_FIELD_NUMBER: _ClassVar[int]
    name: str
    short_name: str
    nickname: str
    object_destruction: float
    mission_success: float
    mission_failure: float
    mission_abort: float
    infoname_id: int
    infocard_id: int
    reputations: _containers.RepeatedCompositeFieldContainer[Reputation]
    bribes: _containers.RepeatedCompositeFieldContainer[Bribe]
    def __init__(self, name: _Optional[str] = ..., short_name: _Optional[str] = ..., nickname: _Optional[str] = ..., object_destruction: _Optional[float] = ..., mission_success: _Optional[float] = ..., mission_failure: _Optional[float] = ..., mission_abort: _Optional[float] = ..., infoname_id: _Optional[int] = ..., infocard_id: _Optional[int] = ..., reputations: _Optional[_Iterable[_Union[Reputation, _Mapping]]] = ..., bribes: _Optional[_Iterable[_Union[Bribe, _Mapping]]] = ...) -> None: ...

class Reputation(_message.Message):
    __slots__ = ("name", "rep", "empathy", "nickname")
    NAME_FIELD_NUMBER: _ClassVar[int]
    REP_FIELD_NUMBER: _ClassVar[int]
    EMPATHY_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    name: str
    rep: float
    empathy: float
    nickname: str
    def __init__(self, name: _Optional[str] = ..., rep: _Optional[float] = ..., empathy: _Optional[float] = ..., nickname: _Optional[str] = ...) -> None: ...

class Bribe(_message.Message):
    __slots__ = ("base_nickname", "chance", "base_info")
    BASE_NICKNAME_FIELD_NUMBER: _ClassVar[int]
    CHANCE_FIELD_NUMBER: _ClassVar[int]
    BASE_INFO_FIELD_NUMBER: _ClassVar[int]
    base_nickname: str
    chance: float
    base_info: BaseInfo
    def __init__(self, base_nickname: _Optional[str] = ..., chance: _Optional[float] = ..., base_info: _Optional[_Union[BaseInfo, _Mapping]] = ...) -> None: ...

class GetGunsReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Gun]
    def __init__(self, items: _Optional[_Iterable[_Union[Gun, _Mapping]]] = ...) -> None: ...

class Gun(_message.Message):
    __slots__ = ("bases", "discovery_tech_compat", "nickname", "name", "type", "price", "hp_type", "ids_name", "ids_info", "volume", "hit_pts", "power_usage", "refire", "range", "toughness", "is_auto_turret", "lootable", "required_ammo", "hull_damage", "energy_damage", "shield_damage", "avg_shield_damage", "damage_type", "life_time", "speed", "gun_turn_rate", "dispersion_angle", "hull_damage_per_sec", "avg_shield_damage_per_sec", "energy_damage_per_sec", "power_usage_per_sec", "avg_efficiency", "hull_efficiency", "shield_efficiency", "energy_damage_efficiency", "damage_bonuses", "missile", "gun_detailed", "num_barrels", "burst_fire", "ammo_limit", "mass", "disco_gun")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    CLASS_FIELD_NUMBER: _ClassVar[int]
    HP_TYPE_FIELD_NUMBER: _ClassVar[int]
    IDS_NAME_FIELD_NUMBER: _ClassVar[int]
    IDS_INFO_FIELD_NUMBER: _ClassVar[int]
    VOLUME_FIELD_NUMBER: _ClassVar[int]
    HIT_PTS_FIELD_NUMBER: _ClassVar[int]
    POWER_USAGE_FIELD_NUMBER: _ClassVar[int]
    REFIRE_FIELD_NUMBER: _ClassVar[int]
    RANGE_FIELD_NUMBER: _ClassVar[int]
    TOUGHNESS_FIELD_NUMBER: _ClassVar[int]
    IS_AUTO_TURRET_FIELD_NUMBER: _ClassVar[int]
    LOOTABLE_FIELD_NUMBER: _ClassVar[int]
    REQUIRED_AMMO_FIELD_NUMBER: _ClassVar[int]
    HULL_DAMAGE_FIELD_NUMBER: _ClassVar[int]
    ENERGY_DAMAGE_FIELD_NUMBER: _ClassVar[int]
    SHIELD_DAMAGE_FIELD_NUMBER: _ClassVar[int]
    AVG_SHIELD_DAMAGE_FIELD_NUMBER: _ClassVar[int]
    DAMAGE_TYPE_FIELD_NUMBER: _ClassVar[int]
    LIFE_TIME_FIELD_NUMBER: _ClassVar[int]
    SPEED_FIELD_NUMBER: _ClassVar[int]
    GUN_TURN_RATE_FIELD_NUMBER: _ClassVar[int]
    DISPERSION_ANGLE_FIELD_NUMBER: _ClassVar[int]
    HULL_DAMAGE_PER_SEC_FIELD_NUMBER: _ClassVar[int]
    AVG_SHIELD_DAMAGE_PER_SEC_FIELD_NUMBER: _ClassVar[int]
    ENERGY_DAMAGE_PER_SEC_FIELD_NUMBER: _ClassVar[int]
    POWER_USAGE_PER_SEC_FIELD_NUMBER: _ClassVar[int]
    AVG_EFFICIENCY_FIELD_NUMBER: _ClassVar[int]
    HULL_EFFICIENCY_FIELD_NUMBER: _ClassVar[int]
    SHIELD_EFFICIENCY_FIELD_NUMBER: _ClassVar[int]
    ENERGY_DAMAGE_EFFICIENCY_FIELD_NUMBER: _ClassVar[int]
    DAMAGE_BONUSES_FIELD_NUMBER: _ClassVar[int]
    MISSILE_FIELD_NUMBER: _ClassVar[int]
    GUN_DETAILED_FIELD_NUMBER: _ClassVar[int]
    NUM_BARRELS_FIELD_NUMBER: _ClassVar[int]
    BURST_FIRE_FIELD_NUMBER: _ClassVar[int]
    AMMO_LIMIT_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    DISCO_GUN_FIELD_NUMBER: _ClassVar[int]
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    nickname: str
    name: str
    type: str
    price: int
    hp_type: str
    ids_name: int
    ids_info: int
    volume: float
    hit_pts: str
    power_usage: float
    refire: float
    range: float
    toughness: float
    is_auto_turret: bool
    lootable: bool
    required_ammo: bool
    hull_damage: int
    energy_damage: int
    shield_damage: int
    avg_shield_damage: int
    damage_type: str
    life_time: float
    speed: float
    gun_turn_rate: float
    dispersion_angle: float
    hull_damage_per_sec: float
    avg_shield_damage_per_sec: float
    energy_damage_per_sec: float
    power_usage_per_sec: float
    avg_efficiency: float
    hull_efficiency: float
    shield_efficiency: float
    energy_damage_efficiency: float
    damage_bonuses: _containers.RepeatedCompositeFieldContainer[DamageBonus]
    missile: Missile
    gun_detailed: GunDetailed
    num_barrels: int
    burst_fire: BurstFire
    ammo_limit: AmmoLimit
    mass: float
    disco_gun: DiscoGun
    def __init__(self, bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ..., nickname: _Optional[str] = ..., name: _Optional[str] = ..., type: _Optional[str] = ..., price: _Optional[int] = ..., hp_type: _Optional[str] = ..., ids_name: _Optional[int] = ..., ids_info: _Optional[int] = ..., volume: _Optional[float] = ..., hit_pts: _Optional[str] = ..., power_usage: _Optional[float] = ..., refire: _Optional[float] = ..., range: _Optional[float] = ..., toughness: _Optional[float] = ..., is_auto_turret: bool = ..., lootable: bool = ..., required_ammo: bool = ..., hull_damage: _Optional[int] = ..., energy_damage: _Optional[int] = ..., shield_damage: _Optional[int] = ..., avg_shield_damage: _Optional[int] = ..., damage_type: _Optional[str] = ..., life_time: _Optional[float] = ..., speed: _Optional[float] = ..., gun_turn_rate: _Optional[float] = ..., dispersion_angle: _Optional[float] = ..., hull_damage_per_sec: _Optional[float] = ..., avg_shield_damage_per_sec: _Optional[float] = ..., energy_damage_per_sec: _Optional[float] = ..., power_usage_per_sec: _Optional[float] = ..., avg_efficiency: _Optional[float] = ..., hull_efficiency: _Optional[float] = ..., shield_efficiency: _Optional[float] = ..., energy_damage_efficiency: _Optional[float] = ..., damage_bonuses: _Optional[_Iterable[_Union[DamageBonus, _Mapping]]] = ..., missile: _Optional[_Union[Missile, _Mapping]] = ..., gun_detailed: _Optional[_Union[GunDetailed, _Mapping]] = ..., num_barrels: _Optional[int] = ..., burst_fire: _Optional[_Union[BurstFire, _Mapping]] = ..., ammo_limit: _Optional[_Union[AmmoLimit, _Mapping]] = ..., mass: _Optional[float] = ..., disco_gun: _Optional[_Union[DiscoGun, _Mapping]] = ..., **kwargs) -> None: ...

class DamageBonus(_message.Message):
    __slots__ = ("type", "modifier")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    MODIFIER_FIELD_NUMBER: _ClassVar[int]
    type: str
    modifier: float
    def __init__(self, type: _Optional[str] = ..., modifier: _Optional[float] = ...) -> None: ...

class Missile(_message.Message):
    __slots__ = ("max_angular_velocity",)
    MAX_ANGULAR_VELOCITY_FIELD_NUMBER: _ClassVar[int]
    max_angular_velocity: float
    def __init__(self, max_angular_velocity: _Optional[float] = ...) -> None: ...

class GunDetailed(_message.Message):
    __slots__ = ("flash_particle_name", "const_effect", "munition_hit_effect")
    FLASH_PARTICLE_NAME_FIELD_NUMBER: _ClassVar[int]
    CONST_EFFECT_FIELD_NUMBER: _ClassVar[int]
    MUNITION_HIT_EFFECT_FIELD_NUMBER: _ClassVar[int]
    flash_particle_name: str
    const_effect: str
    munition_hit_effect: str
    def __init__(self, flash_particle_name: _Optional[str] = ..., const_effect: _Optional[str] = ..., munition_hit_effect: _Optional[str] = ...) -> None: ...

class BurstFire(_message.Message):
    __slots__ = ("sustained_refire", "ammo", "reload_time", "sustained_hull_damage_per_sec", "sustained_avg_shield_damage_per_sec", "sustained_energy_damage_per_sec", "sustained_power_usage_per_sec")
    SUSTAINED_REFIRE_FIELD_NUMBER: _ClassVar[int]
    AMMO_FIELD_NUMBER: _ClassVar[int]
    RELOAD_TIME_FIELD_NUMBER: _ClassVar[int]
    SUSTAINED_HULL_DAMAGE_PER_SEC_FIELD_NUMBER: _ClassVar[int]
    SUSTAINED_AVG_SHIELD_DAMAGE_PER_SEC_FIELD_NUMBER: _ClassVar[int]
    SUSTAINED_ENERGY_DAMAGE_PER_SEC_FIELD_NUMBER: _ClassVar[int]
    SUSTAINED_POWER_USAGE_PER_SEC_FIELD_NUMBER: _ClassVar[int]
    sustained_refire: float
    ammo: int
    reload_time: float
    sustained_hull_damage_per_sec: float
    sustained_avg_shield_damage_per_sec: float
    sustained_energy_damage_per_sec: float
    sustained_power_usage_per_sec: float
    def __init__(self, sustained_refire: _Optional[float] = ..., ammo: _Optional[int] = ..., reload_time: _Optional[float] = ..., sustained_hull_damage_per_sec: _Optional[float] = ..., sustained_avg_shield_damage_per_sec: _Optional[float] = ..., sustained_energy_damage_per_sec: _Optional[float] = ..., sustained_power_usage_per_sec: _Optional[float] = ...) -> None: ...

class DiscoGun(_message.Message):
    __slots__ = ("armor_pen",)
    ARMOR_PEN_FIELD_NUMBER: _ClassVar[int]
    armor_pen: float
    def __init__(self, armor_pen: _Optional[float] = ...) -> None: ...

class GetMinesReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Mine]
    def __init__(self, items: _Optional[_Iterable[_Union[Mine, _Mapping]]] = ...) -> None: ...

class Mine(_message.Message):
    __slots__ = ("name", "price", "ammo_price", "nickname", "projectile_archetype", "ids_name", "ids_info", "hull_damage", "energy_damange", "shield_damage", "power_usage", "value", "refire", "detonation_distance", "radius", "seek_distance", "top_speed", "acceleration", "linear_drag", "life_time", "owner_safe", "toughness", "hit_pts", "lootable", "ammo_limit", "mass", "bases", "discovery_tech_compat")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    AMMO_PRICE_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    PROJECTILE_ARCHETYPE_FIELD_NUMBER: _ClassVar[int]
    IDS_NAME_FIELD_NUMBER: _ClassVar[int]
    IDS_INFO_FIELD_NUMBER: _ClassVar[int]
    HULL_DAMAGE_FIELD_NUMBER: _ClassVar[int]
    ENERGY_DAMANGE_FIELD_NUMBER: _ClassVar[int]
    SHIELD_DAMAGE_FIELD_NUMBER: _ClassVar[int]
    POWER_USAGE_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    REFIRE_FIELD_NUMBER: _ClassVar[int]
    DETONATION_DISTANCE_FIELD_NUMBER: _ClassVar[int]
    RADIUS_FIELD_NUMBER: _ClassVar[int]
    SEEK_DISTANCE_FIELD_NUMBER: _ClassVar[int]
    TOP_SPEED_FIELD_NUMBER: _ClassVar[int]
    ACCELERATION_FIELD_NUMBER: _ClassVar[int]
    LINEAR_DRAG_FIELD_NUMBER: _ClassVar[int]
    LIFE_TIME_FIELD_NUMBER: _ClassVar[int]
    OWNER_SAFE_FIELD_NUMBER: _ClassVar[int]
    TOUGHNESS_FIELD_NUMBER: _ClassVar[int]
    HIT_PTS_FIELD_NUMBER: _ClassVar[int]
    LOOTABLE_FIELD_NUMBER: _ClassVar[int]
    AMMO_LIMIT_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    name: str
    price: int
    ammo_price: int
    nickname: str
    projectile_archetype: str
    ids_name: int
    ids_info: int
    hull_damage: int
    energy_damange: int
    shield_damage: int
    power_usage: float
    value: float
    refire: float
    detonation_distance: float
    radius: float
    seek_distance: int
    top_speed: int
    acceleration: int
    linear_drag: float
    life_time: float
    owner_safe: int
    toughness: float
    hit_pts: int
    lootable: bool
    ammo_limit: AmmoLimit
    mass: float
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    def __init__(self, name: _Optional[str] = ..., price: _Optional[int] = ..., ammo_price: _Optional[int] = ..., nickname: _Optional[str] = ..., projectile_archetype: _Optional[str] = ..., ids_name: _Optional[int] = ..., ids_info: _Optional[int] = ..., hull_damage: _Optional[int] = ..., energy_damange: _Optional[int] = ..., shield_damage: _Optional[int] = ..., power_usage: _Optional[float] = ..., value: _Optional[float] = ..., refire: _Optional[float] = ..., detonation_distance: _Optional[float] = ..., radius: _Optional[float] = ..., seek_distance: _Optional[int] = ..., top_speed: _Optional[int] = ..., acceleration: _Optional[int] = ..., linear_drag: _Optional[float] = ..., life_time: _Optional[float] = ..., owner_safe: _Optional[int] = ..., toughness: _Optional[float] = ..., hit_pts: _Optional[int] = ..., lootable: bool = ..., ammo_limit: _Optional[_Union[AmmoLimit, _Mapping]] = ..., mass: _Optional[float] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ...) -> None: ...

class AmmoLimit(_message.Message):
    __slots__ = ("amount_in_catridge", "max_catridges")
    AMOUNT_IN_CATRIDGE_FIELD_NUMBER: _ClassVar[int]
    MAX_CATRIDGES_FIELD_NUMBER: _ClassVar[int]
    amount_in_catridge: int
    max_catridges: int
    def __init__(self, amount_in_catridge: _Optional[int] = ..., max_catridges: _Optional[int] = ...) -> None: ...

class GetScannersReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Scanner]
    def __init__(self, items: _Optional[_Iterable[_Union[Scanner, _Mapping]]] = ...) -> None: ...

class Scanner(_message.Message):
    __slots__ = ("name", "price", "range", "cargo_scan_range", "lootable", "nickname", "name_id", "info_id", "mass", "bases", "discovery_tech_compat")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    RANGE_FIELD_NUMBER: _ClassVar[int]
    CARGO_SCAN_RANGE_FIELD_NUMBER: _ClassVar[int]
    LOOTABLE_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFO_ID_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    name: str
    price: int
    range: int
    cargo_scan_range: int
    lootable: bool
    nickname: str
    name_id: int
    info_id: int
    mass: float
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    def __init__(self, name: _Optional[str] = ..., price: _Optional[int] = ..., range: _Optional[int] = ..., cargo_scan_range: _Optional[int] = ..., lootable: bool = ..., nickname: _Optional[str] = ..., name_id: _Optional[int] = ..., info_id: _Optional[int] = ..., mass: _Optional[float] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ...) -> None: ...

class GetShieldsReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Shield]
    def __init__(self, items: _Optional[_Iterable[_Union[Shield, _Mapping]]] = ...) -> None: ...

class Shield(_message.Message):
    __slots__ = ("name", "type", "technology", "price", "capacity", "regeneration_rate", "constant_power_draw", "value", "rebuild_power_draw", "off_rebuild_time", "toughness", "hit_pts", "lootable", "nickname", "hp_type", "ids_name", "ids_info", "mass", "bases", "discovery_tech_compat")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    CLASS_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    TECHNOLOGY_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    CAPACITY_FIELD_NUMBER: _ClassVar[int]
    REGENERATION_RATE_FIELD_NUMBER: _ClassVar[int]
    CONSTANT_POWER_DRAW_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    REBUILD_POWER_DRAW_FIELD_NUMBER: _ClassVar[int]
    OFF_REBUILD_TIME_FIELD_NUMBER: _ClassVar[int]
    TOUGHNESS_FIELD_NUMBER: _ClassVar[int]
    HIT_PTS_FIELD_NUMBER: _ClassVar[int]
    LOOTABLE_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    HP_TYPE_FIELD_NUMBER: _ClassVar[int]
    IDS_NAME_FIELD_NUMBER: _ClassVar[int]
    IDS_INFO_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    name: str
    type: str
    technology: str
    price: int
    capacity: int
    regeneration_rate: int
    constant_power_draw: int
    value: float
    rebuild_power_draw: int
    off_rebuild_time: int
    toughness: float
    hit_pts: int
    lootable: bool
    nickname: str
    hp_type: str
    ids_name: int
    ids_info: int
    mass: float
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    def __init__(self, name: _Optional[str] = ..., type: _Optional[str] = ..., technology: _Optional[str] = ..., price: _Optional[int] = ..., capacity: _Optional[int] = ..., regeneration_rate: _Optional[int] = ..., constant_power_draw: _Optional[int] = ..., value: _Optional[float] = ..., rebuild_power_draw: _Optional[int] = ..., off_rebuild_time: _Optional[int] = ..., toughness: _Optional[float] = ..., hit_pts: _Optional[int] = ..., lootable: bool = ..., nickname: _Optional[str] = ..., hp_type: _Optional[str] = ..., ids_name: _Optional[int] = ..., ids_info: _Optional[int] = ..., mass: _Optional[float] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ..., **kwargs) -> None: ...

class GetShipsReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Ship]
    def __init__(self, items: _Optional[_Iterable[_Union[Ship, _Mapping]]] = ...) -> None: ...

class Ship(_message.Message):
    __slots__ = ("nickname", "name", "type", "price", "armor", "hold_size", "nanobots", "batteries", "mass", "power_capacity", "power_recharge_rate", "cruise_speed", "linear_drag", "engine_max_force", "impulse_speed", "thruster_speed", "reverse_fraction", "thrust_capacity", "thrust_recharge", "max_angular_speed_deg_s", "angular_distance_from0_to_half_sec", "time_to90_max_angular_speed", "nudge_force", "strafe_force", "name_id", "info_id", "slots", "biggest_hardpoint", "ship_packages", "bases", "discovery_tech_compat", "disco_ship")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    CLASS_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    ARMOR_FIELD_NUMBER: _ClassVar[int]
    HOLD_SIZE_FIELD_NUMBER: _ClassVar[int]
    NANOBOTS_FIELD_NUMBER: _ClassVar[int]
    BATTERIES_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    POWER_CAPACITY_FIELD_NUMBER: _ClassVar[int]
    POWER_RECHARGE_RATE_FIELD_NUMBER: _ClassVar[int]
    CRUISE_SPEED_FIELD_NUMBER: _ClassVar[int]
    LINEAR_DRAG_FIELD_NUMBER: _ClassVar[int]
    ENGINE_MAX_FORCE_FIELD_NUMBER: _ClassVar[int]
    IMPULSE_SPEED_FIELD_NUMBER: _ClassVar[int]
    THRUSTER_SPEED_FIELD_NUMBER: _ClassVar[int]
    REVERSE_FRACTION_FIELD_NUMBER: _ClassVar[int]
    THRUST_CAPACITY_FIELD_NUMBER: _ClassVar[int]
    THRUST_RECHARGE_FIELD_NUMBER: _ClassVar[int]
    MAX_ANGULAR_SPEED_DEG_S_FIELD_NUMBER: _ClassVar[int]
    ANGULAR_DISTANCE_FROM0_TO_HALF_SEC_FIELD_NUMBER: _ClassVar[int]
    TIME_TO90_MAX_ANGULAR_SPEED_FIELD_NUMBER: _ClassVar[int]
    NUDGE_FORCE_FIELD_NUMBER: _ClassVar[int]
    STRAFE_FORCE_FIELD_NUMBER: _ClassVar[int]
    NAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFO_ID_FIELD_NUMBER: _ClassVar[int]
    SLOTS_FIELD_NUMBER: _ClassVar[int]
    BIGGEST_HARDPOINT_FIELD_NUMBER: _ClassVar[int]
    SHIP_PACKAGES_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    DISCO_SHIP_FIELD_NUMBER: _ClassVar[int]
    nickname: str
    name: str
    type: str
    price: int
    armor: int
    hold_size: int
    nanobots: int
    batteries: int
    mass: float
    power_capacity: int
    power_recharge_rate: int
    cruise_speed: int
    linear_drag: float
    engine_max_force: int
    impulse_speed: float
    thruster_speed: _containers.RepeatedScalarFieldContainer[int]
    reverse_fraction: float
    thrust_capacity: int
    thrust_recharge: int
    max_angular_speed_deg_s: float
    angular_distance_from0_to_half_sec: float
    time_to90_max_angular_speed: float
    nudge_force: float
    strafe_force: float
    name_id: int
    info_id: int
    slots: _containers.RepeatedCompositeFieldContainer[EquipmentSlot]
    biggest_hardpoint: _containers.RepeatedScalarFieldContainer[str]
    ship_packages: _containers.RepeatedCompositeFieldContainer[ShipPackage]
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    disco_ship: DiscoShip
    def __init__(self, nickname: _Optional[str] = ..., name: _Optional[str] = ..., type: _Optional[str] = ..., price: _Optional[int] = ..., armor: _Optional[int] = ..., hold_size: _Optional[int] = ..., nanobots: _Optional[int] = ..., batteries: _Optional[int] = ..., mass: _Optional[float] = ..., power_capacity: _Optional[int] = ..., power_recharge_rate: _Optional[int] = ..., cruise_speed: _Optional[int] = ..., linear_drag: _Optional[float] = ..., engine_max_force: _Optional[int] = ..., impulse_speed: _Optional[float] = ..., thruster_speed: _Optional[_Iterable[int]] = ..., reverse_fraction: _Optional[float] = ..., thrust_capacity: _Optional[int] = ..., thrust_recharge: _Optional[int] = ..., max_angular_speed_deg_s: _Optional[float] = ..., angular_distance_from0_to_half_sec: _Optional[float] = ..., time_to90_max_angular_speed: _Optional[float] = ..., nudge_force: _Optional[float] = ..., strafe_force: _Optional[float] = ..., name_id: _Optional[int] = ..., info_id: _Optional[int] = ..., slots: _Optional[_Iterable[_Union[EquipmentSlot, _Mapping]]] = ..., biggest_hardpoint: _Optional[_Iterable[str]] = ..., ship_packages: _Optional[_Iterable[_Union[ShipPackage, _Mapping]]] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ..., disco_ship: _Optional[_Union[DiscoShip, _Mapping]] = ..., **kwargs) -> None: ...

class EquipmentSlot(_message.Message):
    __slots__ = ("slot_name", "allowed_equip")
    SLOT_NAME_FIELD_NUMBER: _ClassVar[int]
    ALLOWED_EQUIP_FIELD_NUMBER: _ClassVar[int]
    slot_name: str
    allowed_equip: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, slot_name: _Optional[str] = ..., allowed_equip: _Optional[_Iterable[str]] = ...) -> None: ...

class ShipPackage(_message.Message):
    __slots__ = ("nickname",)
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    nickname: str
    def __init__(self, nickname: _Optional[str] = ...) -> None: ...

class DiscoShip(_message.Message):
    __slots__ = ("armor_mult",)
    ARMOR_MULT_FIELD_NUMBER: _ClassVar[int]
    armor_mult: float
    def __init__(self, armor_mult: _Optional[float] = ...) -> None: ...

class GetThrustersReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Thruster]
    def __init__(self, items: _Optional[_Iterable[_Union[Thruster, _Mapping]]] = ...) -> None: ...

class Thruster(_message.Message):
    __slots__ = ("name", "price", "max_force", "power_usage", "efficiency", "value", "hit_pts", "lootable", "nickname", "name_id", "info_id", "mass", "bases", "discovery_tech_compat")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    MAX_FORCE_FIELD_NUMBER: _ClassVar[int]
    POWER_USAGE_FIELD_NUMBER: _ClassVar[int]
    EFFICIENCY_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    HIT_PTS_FIELD_NUMBER: _ClassVar[int]
    LOOTABLE_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFO_ID_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    DISCOVERY_TECH_COMPAT_FIELD_NUMBER: _ClassVar[int]
    name: str
    price: int
    max_force: int
    power_usage: int
    efficiency: float
    value: float
    hit_pts: int
    lootable: bool
    nickname: str
    name_id: int
    info_id: int
    mass: float
    bases: _containers.MessageMap[str, MarketGood]
    discovery_tech_compat: DiscoveryTechCompat
    def __init__(self, name: _Optional[str] = ..., price: _Optional[int] = ..., max_force: _Optional[int] = ..., power_usage: _Optional[int] = ..., efficiency: _Optional[float] = ..., value: _Optional[float] = ..., hit_pts: _Optional[int] = ..., lootable: bool = ..., nickname: _Optional[str] = ..., name_id: _Optional[int] = ..., info_id: _Optional[int] = ..., mass: _Optional[float] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., discovery_tech_compat: _Optional[_Union[DiscoveryTechCompat, _Mapping]] = ...) -> None: ...

class GetTractorsReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Tractor]
    def __init__(self, items: _Optional[_Iterable[_Union[Tractor, _Mapping]]] = ...) -> None: ...

class Tractor(_message.Message):
    __slots__ = ("name", "price", "max_length", "reach_speed", "lootable", "nickname", "short_nickname", "name_id", "info_id", "bases", "mass")
    class BasesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: MarketGood
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[MarketGood, _Mapping]] = ...) -> None: ...
    NAME_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    MAX_LENGTH_FIELD_NUMBER: _ClassVar[int]
    REACH_SPEED_FIELD_NUMBER: _ClassVar[int]
    LOOTABLE_FIELD_NUMBER: _ClassVar[int]
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    SHORT_NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_ID_FIELD_NUMBER: _ClassVar[int]
    INFO_ID_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    MASS_FIELD_NUMBER: _ClassVar[int]
    name: str
    price: int
    max_length: int
    reach_speed: int
    lootable: bool
    nickname: str
    short_nickname: str
    name_id: int
    info_id: int
    bases: _containers.MessageMap[str, MarketGood]
    mass: float
    def __init__(self, name: _Optional[str] = ..., price: _Optional[int] = ..., max_length: _Optional[int] = ..., reach_speed: _Optional[int] = ..., lootable: bool = ..., nickname: _Optional[str] = ..., short_nickname: _Optional[str] = ..., name_id: _Optional[int] = ..., info_id: _Optional[int] = ..., bases: _Optional[_Mapping[str, MarketGood]] = ..., mass: _Optional[float] = ...) -> None: ...

class GetHashesReply(_message.Message):
    __slots__ = ("hashes_by_nick",)
    class HashesByNickEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: Hash
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[Hash, _Mapping]] = ...) -> None: ...
    HASHES_BY_NICK_FIELD_NUMBER: _ClassVar[int]
    hashes_by_nick: _containers.MessageMap[str, Hash]
    def __init__(self, hashes_by_nick: _Optional[_Mapping[str, Hash]] = ...) -> None: ...

class Hash(_message.Message):
    __slots__ = ("int32", "uint32", "hex")
    INT32_FIELD_NUMBER: _ClassVar[int]
    UINT32_FIELD_NUMBER: _ClassVar[int]
    HEX_FIELD_NUMBER: _ClassVar[int]
    int32: int
    uint32: int
    hex: str
    def __init__(self, int32: _Optional[int] = ..., uint32: _Optional[int] = ..., hex: _Optional[str] = ...) -> None: ...

class GetPoBsReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[PoB]
    def __init__(self, items: _Optional[_Iterable[_Union[PoB, _Mapping]]] = ...) -> None: ...

class PoBCore(_message.Message):
    __slots__ = ("nickname", "name", "pos", "level", "money", "health", "defense_mode", "system_nick", "system_name", "faction_nick", "faction_name", "forum_thread_url", "cargo_space_left", "base_pos", "sector_coord", "region")
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    POS_FIELD_NUMBER: _ClassVar[int]
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    MONEY_FIELD_NUMBER: _ClassVar[int]
    HEALTH_FIELD_NUMBER: _ClassVar[int]
    DEFENSE_MODE_FIELD_NUMBER: _ClassVar[int]
    SYSTEM_NICK_FIELD_NUMBER: _ClassVar[int]
    SYSTEM_NAME_FIELD_NUMBER: _ClassVar[int]
    FACTION_NICK_FIELD_NUMBER: _ClassVar[int]
    FACTION_NAME_FIELD_NUMBER: _ClassVar[int]
    FORUM_THREAD_URL_FIELD_NUMBER: _ClassVar[int]
    CARGO_SPACE_LEFT_FIELD_NUMBER: _ClassVar[int]
    BASE_POS_FIELD_NUMBER: _ClassVar[int]
    SECTOR_COORD_FIELD_NUMBER: _ClassVar[int]
    REGION_FIELD_NUMBER: _ClassVar[int]
    nickname: str
    name: str
    pos: str
    level: int
    money: int
    health: float
    defense_mode: int
    system_nick: str
    system_name: str
    faction_nick: str
    faction_name: str
    forum_thread_url: str
    cargo_space_left: int
    base_pos: Pos
    sector_coord: str
    region: str
    def __init__(self, nickname: _Optional[str] = ..., name: _Optional[str] = ..., pos: _Optional[str] = ..., level: _Optional[int] = ..., money: _Optional[int] = ..., health: _Optional[float] = ..., defense_mode: _Optional[int] = ..., system_nick: _Optional[str] = ..., system_name: _Optional[str] = ..., faction_nick: _Optional[str] = ..., faction_name: _Optional[str] = ..., forum_thread_url: _Optional[str] = ..., cargo_space_left: _Optional[int] = ..., base_pos: _Optional[_Union[Pos, _Mapping]] = ..., sector_coord: _Optional[str] = ..., region: _Optional[str] = ...) -> None: ...

class PoB(_message.Message):
    __slots__ = ("core", "shop_items")
    CORE_FIELD_NUMBER: _ClassVar[int]
    SHOP_ITEMS_FIELD_NUMBER: _ClassVar[int]
    core: PoBCore
    shop_items: _containers.RepeatedCompositeFieldContainer[ShopItem]
    def __init__(self, core: _Optional[_Union[PoBCore, _Mapping]] = ..., shop_items: _Optional[_Iterable[_Union[ShopItem, _Mapping]]] = ...) -> None: ...

class ShopItem(_message.Message):
    __slots__ = ("nickname", "name", "category", "id", "quantity", "price", "sell_price", "min_stock", "max_stock")
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    CATEGORY_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    QUANTITY_FIELD_NUMBER: _ClassVar[int]
    PRICE_FIELD_NUMBER: _ClassVar[int]
    SELL_PRICE_FIELD_NUMBER: _ClassVar[int]
    MIN_STOCK_FIELD_NUMBER: _ClassVar[int]
    MAX_STOCK_FIELD_NUMBER: _ClassVar[int]
    nickname: str
    name: str
    category: str
    id: int
    quantity: int
    price: int
    sell_price: int
    min_stock: int
    max_stock: int
    def __init__(self, nickname: _Optional[str] = ..., name: _Optional[str] = ..., category: _Optional[str] = ..., id: _Optional[int] = ..., quantity: _Optional[int] = ..., price: _Optional[int] = ..., sell_price: _Optional[int] = ..., min_stock: _Optional[int] = ..., max_stock: _Optional[int] = ...) -> None: ...

class GetPoBGoodsReply(_message.Message):
    __slots__ = ("items",)
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[PoBGood]
    def __init__(self, items: _Optional[_Iterable[_Union[PoBGood, _Mapping]]] = ...) -> None: ...

class PoBGood(_message.Message):
    __slots__ = ("nickname", "name", "total_buyable_from_bases", "total_sellable_to_bases", "best_price_to_buy", "best_price_to_sell", "category", "any_base_sells", "any_base_buys", "bases")
    NICKNAME_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_BUYABLE_FROM_BASES_FIELD_NUMBER: _ClassVar[int]
    TOTAL_SELLABLE_TO_BASES_FIELD_NUMBER: _ClassVar[int]
    BEST_PRICE_TO_BUY_FIELD_NUMBER: _ClassVar[int]
    BEST_PRICE_TO_SELL_FIELD_NUMBER: _ClassVar[int]
    CATEGORY_FIELD_NUMBER: _ClassVar[int]
    ANY_BASE_SELLS_FIELD_NUMBER: _ClassVar[int]
    ANY_BASE_BUYS_FIELD_NUMBER: _ClassVar[int]
    BASES_FIELD_NUMBER: _ClassVar[int]
    nickname: str
    name: str
    total_buyable_from_bases: int
    total_sellable_to_bases: int
    best_price_to_buy: int
    best_price_to_sell: int
    category: str
    any_base_sells: bool
    any_base_buys: bool
    bases: _containers.RepeatedCompositeFieldContainer[PoBGoodBase]
    def __init__(self, nickname: _Optional[str] = ..., name: _Optional[str] = ..., total_buyable_from_bases: _Optional[int] = ..., total_sellable_to_bases: _Optional[int] = ..., best_price_to_buy: _Optional[int] = ..., best_price_to_sell: _Optional[int] = ..., category: _Optional[str] = ..., any_base_sells: bool = ..., any_base_buys: bool = ..., bases: _Optional[_Iterable[_Union[PoBGoodBase, _Mapping]]] = ...) -> None: ...

class PoBGoodBase(_message.Message):
    __slots__ = ("shop_item", "base")
    SHOP_ITEM_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    shop_item: ShopItem
    base: PoBCore
    def __init__(self, shop_item: _Optional[_Union[ShopItem, _Mapping]] = ..., base: _Optional[_Union[PoBCore, _Mapping]] = ...) -> None: ...

class GetGraphPathsInput(_message.Message):
    __slots__ = ("queries",)
    QUERIES_FIELD_NUMBER: _ClassVar[int]
    queries: _containers.RepeatedCompositeFieldContainer[GraphPathQuery]
    def __init__(self, queries: _Optional[_Iterable[_Union[GraphPathQuery, _Mapping]]] = ...) -> None: ...

class GraphPathQuery(_message.Message):
    __slots__ = ("to",)
    FROM_FIELD_NUMBER: _ClassVar[int]
    TO_FIELD_NUMBER: _ClassVar[int]
    to: str
    def __init__(self, to: _Optional[str] = ..., **kwargs) -> None: ...

class GetGraphPathsReply(_message.Message):
    __slots__ = ("answers",)
    ANSWERS_FIELD_NUMBER: _ClassVar[int]
    answers: _containers.RepeatedCompositeFieldContainer[GetGraphPathsAnswer]
    def __init__(self, answers: _Optional[_Iterable[_Union[GetGraphPathsAnswer, _Mapping]]] = ...) -> None: ...

class GetGraphPathsAnswer(_message.Message):
    __slots__ = ("route", "time", "error")
    ROUTE_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    route: GraphPathQuery
    time: GraphPathTime
    error: str
    def __init__(self, route: _Optional[_Union[GraphPathQuery, _Mapping]] = ..., time: _Optional[_Union[GraphPathTime, _Mapping]] = ..., error: _Optional[str] = ...) -> None: ...

class GraphPathTime(_message.Message):
    __slots__ = ("transport", "frigate", "freighter")
    TRANSPORT_FIELD_NUMBER: _ClassVar[int]
    FRIGATE_FIELD_NUMBER: _ClassVar[int]
    FREIGHTER_FIELD_NUMBER: _ClassVar[int]
    transport: int
    frigate: int
    freighter: int
    def __init__(self, transport: _Optional[int] = ..., frigate: _Optional[int] = ..., freighter: _Optional[int] = ...) -> None: ...
